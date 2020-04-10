package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/token/internal/types"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error)
	GetOrNewAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error)
	GetAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error)
	SetAccount(ctx sdk.Context, acc types.Account) error
	UpdateAccount(ctx sdk.Context, acc types.Account) error
}

func (k Keeper) NewAccountWithAddress(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error) {
	acc = types.NewBaseAccountWithAddress(contractID, addr)
	if err = k.SetAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (k Keeper) SetAccount(ctx sdk.Context, acc types.Account) error {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(acc.GetContractID(), acc.GetAddress())
	if store.Has(accKey) {
		return sdkerrors.Wrap(types.ErrAccountExist, acc.GetAddress().String())
	}
	store.Set(accKey, k.cdc.MustMarshalBinaryBare(acc))

	// Set Account if not exists yet
	account := k.accountKeeper.GetAccount(ctx, acc.GetAddress())
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, acc.GetAddress())
		k.accountKeeper.SetAccount(ctx, account)
	}
	return nil
}

func (k Keeper) UpdateAccount(ctx sdk.Context, acc types.Account) error {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(acc.GetContractID(), acc.GetAddress())
	if !store.Has(accKey) {
		return sdkerrors.Wrap(types.ErrAccountNotExist, acc.GetAddress().String())
	}
	store.Set(accKey, k.cdc.MustMarshalBinaryBare(acc))
	return nil
}

func (k Keeper) GetOrNewAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error) {
	acc, err = k.GetAccount(ctx, contractID, addr)
	if err != nil {
		acc, err = k.NewAccountWithAddress(ctx, contractID, addr)
		if err != nil {
			return nil, err
		}
	}
	return acc, nil
}

func (k Keeper) GetAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc types.Account, err error) {
	store := ctx.KVStore(k.storeKey)
	accKey := types.AccountKey(contractID, addr)
	if !store.Has(accKey) {
		return nil, sdkerrors.Wrap(types.ErrAccountNotExist, addr.String())
	}
	bz := store.Get(accKey)
	return k.mustDecodeAccount(bz), nil
}

func (k Keeper) mustDecodeAccount(bz []byte) (acc types.Account) {
	err := k.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
