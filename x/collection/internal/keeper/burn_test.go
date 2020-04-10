package keeper

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestKeeper_BurnFT(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.NoError(t, keeper.BurnFT(ctx, defaultContractID, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))))
	require.EqualError(t, keeper.BurnFT(ctx, wrongContractID, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr1.String(), types.NewBurnPermission(wrongContractID).String()).Error())
	require.EqualError(t, keeper.BurnFT(ctx, defaultContractID, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewBurnPermission(defaultContractID).String()).Error())
	require.EqualError(t, keeper.BurnFT(ctx, defaultContractID, addr3, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr3.String(), types.NewBurnPermission(defaultContractID).String()).Error())
	require.EqualError(t, keeper.BurnFT(ctx, defaultContractID, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))), sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", addr1.String(), types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)).String()).Error())

	require.NoError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))))
	require.EqualError(t, keeper.BurnFT(ctx, defaultContractID, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2, types.NewBurnPermission(defaultContractID)).Error())
	require.EqualError(t, keeper.BurnFT(ctx, defaultContractID, addr1, types.NewCoins(types.NewCoin("0000000200000000", sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", addr1.String(), types.NewCoin("0000000200000000", sdk.NewInt(1)).String()).Error())
}

func TestKeeper_BurnFTFrom(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)
	prepareProxy(ctx, t)

	require.EqualError(t, keeper.BurnFTFrom(ctx, wrongContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), wrongContractID).Error())
	require.NoError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", addr2.String(), types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)).String()).Error())

	require.NoError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))))
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr2, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2, types.NewBurnPermission(defaultContractID)).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin("0000000200000000", sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", addr2.String(), types.NewCoin("0000000200000000", sdk.NewInt(1)).String()).Error())
}

func TestKeeper_BurnNFT(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)
	i, err := keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTCount)
	require.NoError(t, err)
	require.Equal(t, int64(5), i.Int64())
	i, err = keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTMint)
	require.NoError(t, err)
	require.Equal(t, int64(5), i.Int64())
	i, err = keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTBurn)
	require.NoError(t, err)
	require.Equal(t, int64(0), i.Int64())

	require.NoError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, defaultTokenID4))
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, defaultTokenID4), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID4).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr2, defaultTokenID4), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID4).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr3, defaultTokenID4), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID4).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, wrongContractID, addr1, defaultTokenID4), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID4).Error())

	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID2, defaultTokenID3))
	require.NoError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, defaultTokenID1))

	i, err = keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTCount)
	require.NoError(t, err)
	require.Equal(t, int64(1), i.Int64())
	i, err = keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTMint)
	require.NoError(t, err)
	require.Equal(t, int64(5), i.Int64())
	i, err = keeper.GetNFTCountInt(ctx, defaultContractID, defaultTokenType, types.QueryNFTBurn)
	require.NoError(t, err)
	require.Equal(t, int64(4), i.Int64())
}

func TestKeeper_BurnNFTFrom(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)
	prepareProxy(ctx, t)

	require.NoError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID4))
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID4), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID4).Error())
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr3, addr2, defaultTokenID4), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr3.String(), addr2.String(), defaultContractID).Error())
	require.EqualError(t, keeper.BurnNFTFrom(ctx, wrongContractID, addr1, addr2, defaultTokenID4), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), wrongContractID).Error())

	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr2, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr2, defaultTokenID2, defaultTokenID3))
	require.NoError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
}

func TestMintBurn(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	const (
		wrongTokenID = "12345678"
	)
	require.EqualError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID2, wrongTokenID, defaultName, defaultMeta, addr1)), sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", defaultContractID2, wrongTokenID[:types.TokenTypeLength]).Error())
	require.EqualError(t, keeper.MintNFT(ctx, addr3, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr3.String(), types.NewMintPermission(defaultContractID).String()).Error())

	require.NoError(t, keeper.BurnFT(ctx, defaultContractID, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, wrongTokenID), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, wrongTokenID).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr3, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr3.String(), types.NewBurnPermission(defaultContractID).String()).Error())
}

func TestBurnNFTScenario(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID4))

	require.NoError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, defaultTokenID1))

	_, err := keeper.GetNFT(ctx, defaultContractID, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultContractID, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID4, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
}

func TestBurnNFTFromSuccess(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// success case
	// addr1 has: burn permission, approved
	// addr2 has: token

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID4))

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))))

	// addr2 approves addr1 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr1, addr2))

	// test burnNFTFrom
	require.NoError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))

	_, err := keeper.GetNFT(ctx, defaultContractID, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultContractID, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID4, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenIDFT, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
}

func TestBurnNFTFromFailure1(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case1
	// addr1 has: burn permission, approved, token
	// addr2 has: nothing

	// addr2 approves addr1 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr1, addr2))

	// test burnNFTFrom, burnFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID1, addr2.String()).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", addr2.String(), types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)).String()).Error())
}

func TestBurnNFTFromFailure2(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case2
	// addr1 has: burn permission (not approved)
	// addr2 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), defaultContractID).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), defaultContractID).Error())
}

func TestBurnNFTFromFailure3(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case3
	// addr2 has: approved (no permission)
	// addr3 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr3, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr3, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))
	// addr3 approves addr2 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr2, addr3))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr2, addr3, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewBurnPermission(defaultContractID).String()).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr2, addr3, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewBurnPermission(defaultContractID).String()).Error())
}
