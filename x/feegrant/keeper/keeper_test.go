package keeper_test

import (
	"context"
	"testing"

	"github.com/line/lbm-sdk/x/feegrant/types"
	tmproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/feegrant/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app     *simapp.SimApp
	sdkCtx  sdk.Context
	addrs   []sdk.AccAddress
	msgSrvr types.MsgServer
	ctx     context.Context
	atom    sdk.Coins
	keeper  keeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.sdkCtx = ctx
	suite.addrs = simapp.AddTestAddrsIncremental(app, ctx, 4, sdk.NewInt(30000000))
	suite.ctx = sdk.WrapSDKContext(ctx)
	suite.keeper = suite.app.FeeGrantKeeper
	suite.msgSrvr = keeper.NewMsgServerImpl(suite.keeper)
	suite.atom = sdk.NewCoins(sdk.NewCoin("atom", sdk.NewInt(555)))
}

func (suite *KeeperTestSuite) TestKeeperCrud() {
	// some helpers
	eth := sdk.NewCoins(sdk.NewInt64Coin("eth", 123))
	exp := suite.sdkCtx.BlockTime().AddDate(1, 0, 0)
	basic := &types.BasicAllowance{
		SpendLimit: suite.atom,
		Expiration: &exp,
	}

	basic2 := &types.BasicAllowance{
		SpendLimit: eth,
		Expiration: &exp,
	}

	// let's set up some initial state here
	err := suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[0], suite.addrs[1], basic)
	suite.Require().NoError(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[0], suite.addrs[2], basic2)
	suite.Require().NoError(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[1], suite.addrs[2], basic)
	suite.Require().NoError(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[1], suite.addrs[3], basic)
	suite.Require().NoError(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[3], suite.addrs[0], basic2)
	suite.Require().NoError(err)

	// remove some, overwrite other
	_, err = suite.msgSrvr.RevokeAllowance(suite.ctx, &types.MsgRevokeAllowance{Granter: suite.addrs[0].String(), Grantee: suite.addrs[1].String()})
	suite.Require().NoError(err)
	_, err = suite.msgSrvr.RevokeAllowance(suite.ctx, &types.MsgRevokeAllowance{Granter: suite.addrs[0].String(), Grantee: suite.addrs[2].String()})
	suite.Require().NoError(err)

	// revoke non-exist fee allowance
	_, err = suite.msgSrvr.RevokeAllowance(suite.ctx, &types.MsgRevokeAllowance{Granter: suite.addrs[0].String(), Grantee: suite.addrs[2].String()})
	suite.Require().Error(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[0], suite.addrs[2], basic)
	suite.Require().NoError(err)

	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[1], suite.addrs[2], basic2)
	suite.Require().NoError(err)

	// end state:
	// addr -> addr3 (basic)
	// addr2 -> addr3 (basic2), addr4(basic)
	// addr4 -> addr (basic2)

	// then lots of queries
	cases := map[string]struct {
		grantee   sdk.AccAddress
		granter   sdk.AccAddress
		allowance types.FeeAllowanceI
	}{
		"addr revoked": {
			granter: suite.addrs[0],
			grantee: suite.addrs[1],
		},
		"addr revoked and added": {
			granter:   suite.addrs[0],
			grantee:   suite.addrs[2],
			allowance: basic,
		},
		"addr never there": {
			granter: suite.addrs[0],
			grantee: suite.addrs[3],
		},
		"addr modified": {
			granter:   suite.addrs[1],
			grantee:   suite.addrs[2],
			allowance: basic2,
		},
	}

	for name, tc := range cases {
		tc := tc
		suite.Run(name, func() {
			allow, _ := suite.keeper.GetAllowance(suite.sdkCtx, tc.granter, tc.grantee)

			if tc.allowance == nil {
				suite.Nil(allow)
				return
			}
			suite.NotNil(allow)
			suite.Equal(tc.allowance, allow)
		})
	}
	accAddr := sdk.AccAddress("link1k907plrxssuh2dxsd8k2jtp3de2t4xyhq6rkxd")
	suite.Require().NoError(err)

	// let's grant and revoke authorization to non existing account
	err = suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[3], accAddr, basic2)
	suite.Require().NoError(err)

	_, err = suite.keeper.GetAllowance(suite.sdkCtx, suite.addrs[3], accAddr)
	suite.Require().NoError(err)

	_, err = suite.msgSrvr.RevokeAllowance(suite.ctx, &types.MsgRevokeAllowance{Granter: suite.addrs[3].String(), Grantee: accAddr.String()})
	suite.Require().NoError(err)

}

func (suite *KeeperTestSuite) TestUseGrantedFee() {
	eth := sdk.NewCoins(sdk.NewInt64Coin("eth", 123))
	blockTime := suite.sdkCtx.BlockTime()
	oneYear := blockTime.AddDate(1, 0, 0)

	future := &types.BasicAllowance{
		SpendLimit: suite.atom,
		Expiration: &oneYear,
	}

	// for testing limits of the contract
	hugeAtom := sdk.NewCoins(sdk.NewInt64Coin("atom", 9999))
	smallAtom := sdk.NewCoins(sdk.NewInt64Coin("atom", 1))
	futureAfterSmall := &types.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("atom", 554)),
		Expiration: &oneYear,
	}

	// then lots of queries
	cases := map[string]struct {
		grantee sdk.AccAddress
		granter sdk.AccAddress
		fee     sdk.Coins
		allowed bool
		final   types.FeeAllowanceI
	}{
		"use entire pot": {
			granter: suite.addrs[0],
			grantee: suite.addrs[1],
			fee:     suite.atom,
			allowed: true,
			final:   nil,
		},
		"too high": {
			granter: suite.addrs[0],
			grantee: suite.addrs[1],
			fee:     hugeAtom,
			allowed: false,
			final:   future,
		},
		"use a little": {
			granter: suite.addrs[0],
			grantee: suite.addrs[1],
			fee:     smallAtom,
			allowed: true,
			final:   futureAfterSmall,
		},
	}

	for name, tc := range cases {
		tc := tc
		suite.Run(name, func() {
			err := suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[0], suite.addrs[1], future)
			suite.Require().NoError(err)

			err = suite.keeper.UseGrantedFees(suite.sdkCtx, tc.granter, tc.grantee, tc.fee, []sdk.Msg{})
			if tc.allowed {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}

			loaded, _ := suite.keeper.GetAllowance(suite.sdkCtx, tc.granter, tc.grantee)
			suite.Equal(tc.final, loaded)
		})
	}

	expired := &types.BasicAllowance{
		SpendLimit: eth,
		Expiration: &blockTime,
	}
	// creating expired feegrant
	ctx := suite.sdkCtx.WithBlockTime(oneYear)
	err := suite.keeper.GrantAllowance(ctx, suite.addrs[0], suite.addrs[2], expired)
	suite.Require().NoError(err)

	// expect error: feegrant expired
	err = suite.keeper.UseGrantedFees(ctx, suite.addrs[0], suite.addrs[2], eth, []sdk.Msg{})
	suite.Error(err)
	suite.Contains(err.Error(), "fee allowance expired")

	// verify: feegrant is revoked
	_, err = suite.keeper.GetAllowance(ctx, suite.addrs[0], suite.addrs[2])
	suite.Error(err)
	suite.Contains(err.Error(), "fee-grant not found")

}

func (suite *KeeperTestSuite) TestIterateGrants() {
	eth := sdk.NewCoins(sdk.NewInt64Coin("eth", 123))
	exp := suite.sdkCtx.BlockTime().AddDate(1, 0, 0)

	allowance := &types.BasicAllowance{
		SpendLimit: suite.atom,
		Expiration: &exp,
	}

	allowance1 := &types.BasicAllowance{
		SpendLimit: eth,
		Expiration: &exp,
	}

	suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[0], suite.addrs[1], allowance)
	suite.keeper.GrantAllowance(suite.sdkCtx, suite.addrs[2], suite.addrs[1], allowance1)

	suite.keeper.IterateAllFeeAllowances(suite.sdkCtx, func(grant types.Grant) bool {
		suite.Require().Equal(suite.addrs[1].String(), grant.Grantee)
		suite.Require().Contains([]string{suite.addrs[0].String(), suite.addrs[2].String()}, grant.Granter)
		return true
	})

}
