package keeper

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/line/lbm-sdk/x/feegrant/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

// Allowance returns fee granted to the grantee by the granter.
func (q Keeper) Allowance(c context.Context, req *types.QueryAllowanceRequest) (*types.QueryAllowanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	granterAddr := sdk.AccAddress(req.Granter)
	granteeAddr := sdk.AccAddress(req.Grantee)

	ctx := sdk.UnwrapSDKContext(c)

	feeAllowance, err := q.GetAllowance(ctx, granterAddr, granteeAddr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	msg, ok := feeAllowance.(proto.Message)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't proto marshal %T", msg)
	}

	feeAllowanceAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &types.QueryAllowanceResponse{
		Allowance: &types.Grant{
			Granter:   granterAddr.String(),
			Grantee:   granteeAddr.String(),
			Allowance: feeAllowanceAny,
		},
	}, nil
}

// Allowances queries all the allowances granted to the given grantee.
func (q Keeper) Allowances(c context.Context, req *types.QueryAllowancesRequest) (*types.QueryAllowancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err := sdk.ValidateAccAddress(req.Grantee)
	if err != nil {
		return nil, err
	}
	granteeAddr := sdk.AccAddress(req.Grantee)

	ctx := sdk.UnwrapSDKContext(c)

	var grants []*types.Grant

	store := ctx.KVStore(q.storeKey)
	grantsStore := prefix.NewStore(store, types.FeeAllowancePrefixByGrantee(granteeAddr))

	pageRes, err := query.Paginate(grantsStore, req.Pagination, func(key []byte, value []byte) error {
		var grant types.Grant

		if err := q.cdc.UnmarshalBinaryBare(value, &grant); err != nil {
			return err
		}

		grants = append(grants, &grant)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllowancesResponse{Allowances: grants, Pagination: pageRes}, nil
}
