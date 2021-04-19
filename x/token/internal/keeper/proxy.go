package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

var (
	ApprovedValue = []byte{0x01}
)

type ProxyKeeper interface {
	IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) bool
	SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error
	DeleteApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error
}

func (k Keeper) IsApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	approvedKey := types.TokenApprovedKey(k.getContractID(ctx), proxy, approver)
	return store.Has(approvedKey)
}

func (k Keeper) SetApproved(ctx sdk.Context, proxy sdk.AccAddress, approver sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.TokenKey(k.getContractID(ctx))) {
		return sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s", k.getContractID(ctx))
	}

	approvedKey := types.TokenApprovedKey(k.getContractID(ctx), proxy, approver)
	if store.Has(approvedKey) {
		return sdkerrors.Wrapf(types.ErrTokenAlreadyApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), approver.String(), k.getContractID(ctx))
	}
	store.Set(approvedKey, ApprovedValue)

	// Set Account if not exists yet
	account := k.accountKeeper.GetAccount(ctx, proxy)
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, proxy)
		k.accountKeeper.SetAccount(ctx, account)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeApproveToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyApprover, approver.String()),
		),
	})

	return nil
}

func (k Keeper) GetApprovers(ctx sdk.Context, proxy sdk.AccAddress) (accAds []sdk.AccAddress, err error) {
	_, err = k.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	k.iterateApprovers(ctx, proxy, false, func(address sdk.AccAddress) bool {
		accAds = append(accAds, address)
		return false
	})
	return accAds, nil
}

func (k Keeper) iterateApprovers(ctx sdk.Context, prefix sdk.AccAddress, reverse bool, process func(accAd sdk.AccAddress) bool) {
	store := ctx.KVStore(k.storeKey)
	prefixKey := types.TokenApproversKey(k.getContractID(ctx), prefix)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, prefixKey)
	} else {
		iter = sdk.KVStorePrefixIterator(store, prefixKey)
	}
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		bz := iter.Key()
		approver := bz[len(prefixKey):]
		if process(approver) {
			return
		}
		iter.Next()
	}
}