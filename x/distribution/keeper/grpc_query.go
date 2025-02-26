package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/distribution/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

var _ types.QueryServer = Keeper{}

// Params queries params of distribution module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)

	return &types.QueryParamsResponse{Params: params}, nil
}

// ValidatorOutstandingRewards queries rewards of a validator address
func (k Keeper) ValidatorOutstandingRewards(c context.Context, req *types.QueryValidatorOutstandingRewardsRequest) (*types.QueryValidatorOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	err := sdk.ValidateValAddress(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	rewards := k.GetValidatorOutstandingRewards(ctx, sdk.ValAddress(req.ValidatorAddress))

	return &types.QueryValidatorOutstandingRewardsResponse{Rewards: rewards}, nil
}

// ValidatorCommission queries accumulated commission for a validator
func (k Keeper) ValidatorCommission(c context.Context, req *types.QueryValidatorCommissionRequest) (*types.QueryValidatorCommissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	err := sdk.ValidateValAddress(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	commission := k.GetValidatorAccumulatedCommission(ctx, sdk.ValAddress(req.ValidatorAddress))

	return &types.QueryValidatorCommissionResponse{Commission: commission}, nil
}

// ValidatorSlashes queries slash events of a validator
func (k Keeper) ValidatorSlashes(c context.Context, req *types.QueryValidatorSlashesRequest) (*types.QueryValidatorSlashesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	if req.EndingHeight < req.StartingHeight {
		return nil, status.Errorf(codes.InvalidArgument, "starting height greater than ending height (%d > %d)", req.StartingHeight, req.EndingHeight)
	}

	ctx := sdk.UnwrapSDKContext(c)
	events := make([]types.ValidatorSlashEvent, 0)
	store := ctx.KVStore(k.storeKey)
	err := sdk.ValidateValAddress(req.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address")
	}
	slashesStore := prefix.NewStore(store, types.GetValidatorSlashEventPrefix(sdk.ValAddress(req.ValidatorAddress)))

	pageRes, err := query.FilteredPaginate(slashesStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var result types.ValidatorSlashEvent
		err := k.cdc.UnmarshalBinaryBare(value, &result)

		if err != nil {
			return false, err
		}

		if result.ValidatorPeriod < req.StartingHeight || result.ValidatorPeriod > req.EndingHeight {
			return false, nil
		}

		if accumulate {
			events = append(events, result)
		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorSlashesResponse{Slashes: events, Pagination: pageRes}, nil
}

// DelegationRewards the total rewards accrued by a delegation
func (k Keeper) DelegationRewards(c context.Context, req *types.QueryDelegationRewardsRequest) (*types.QueryDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	err := sdk.ValidateValAddress(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	val := k.stakingKeeper.Validator(ctx, sdk.ValAddress(req.ValidatorAddress))
	if val == nil {
		return nil, sdkerrors.Wrap(types.ErrNoValidatorExists, req.ValidatorAddress)
	}

	err = sdk.ValidateAccAddress(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	del := k.stakingKeeper.Delegation(ctx, sdk.AccAddress(req.DelegatorAddress), sdk.ValAddress(req.ValidatorAddress))
	if del == nil {
		return nil, types.ErrNoDelegationExists
	}

	endingPeriod := k.IncrementValidatorPeriod(ctx, val)
	rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

	return &types.QueryDelegationRewardsResponse{Rewards: rewards}, nil
}

// DelegationTotalRewards the total rewards accrued by a each validator
func (k Keeper) DelegationTotalRewards(c context.Context, req *types.QueryDelegationTotalRewardsRequest) (*types.QueryDelegationTotalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	err := sdk.ValidateAccAddress(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	k.stakingKeeper.IterateDelegations(
		ctx, sdk.AccAddress(req.DelegatorAddress),
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			val := k.stakingKeeper.Validator(ctx, valAddr)
			endingPeriod := k.IncrementValidatorPeriod(ctx, val)
			delReward := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)

			delRewards = append(delRewards, types.NewDelegationDelegatorReward(valAddr, delReward))
			total = total.Add(delReward...)
			return false
		},
	)

	return &types.QueryDelegationTotalRewardsResponse{Rewards: delRewards, Total: total}, nil
}

// DelegatorValidators queries the validators list of a delegator
func (k Keeper) DelegatorValidators(c context.Context, req *types.QueryDelegatorValidatorsRequest) (*types.QueryDelegatorValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	err := sdk.ValidateAccAddress(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	var validators []string

	k.stakingKeeper.IterateDelegations(
		ctx, sdk.AccAddress(req.DelegatorAddress),
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			validators = append(validators, del.GetValidatorAddr().String())
			return false
		},
	)

	return &types.QueryDelegatorValidatorsResponse{Validators: validators}, nil
}

// DelegatorWithdrawAddress queries Query/delegatorWithdrawAddress
func (k Keeper) DelegatorWithdrawAddress(c context.Context, req *types.QueryDelegatorWithdrawAddressRequest) (*types.QueryDelegatorWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}
	err := sdk.ValidateAccAddress(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(req.DelegatorAddress))

	return &types.QueryDelegatorWithdrawAddressResponse{WithdrawAddress: withdrawAddr.String()}, nil
}

// CommunityPool queries the community pool coins
func (k Keeper) CommunityPool(c context.Context, req *types.QueryCommunityPoolRequest) (*types.QueryCommunityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	pool := k.GetFeePoolCommunityCoins(ctx)

	return &types.QueryCommunityPoolResponse{Pool: pool}, nil
}
