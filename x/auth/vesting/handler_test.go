package vesting_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth/vesting"
	"github.com/line/lbm-sdk/x/auth/vesting/types"
)

type HandlerTestSuite struct {
	suite.Suite

	handler sdk.Handler
	app     *simapp.SimApp
}

func (suite *HandlerTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)

	suite.handler = vesting.NewHandler(app.AccountKeeper, app.BankKeeper)
	suite.app = app
}

func (suite *HandlerTestSuite) TestMsgCreateVestingAccount() {
	ctx := suite.app.BaseApp.NewContext(false, ocproto.Header{Height: suite.app.LastBlockHeight() + 1})

	balances := sdk.NewCoins(sdk.NewInt64Coin("test", 1000))
	addr1 := sdk.BytesToAccAddress([]byte("addr1_______________"))
	addr2 := sdk.BytesToAccAddress([]byte("addr2_______________"))
	addr3 := sdk.BytesToAccAddress([]byte("addr3_______________"))

	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	suite.app.AccountKeeper.SetAccount(ctx, acc1)
	suite.Require().NoError(suite.app.BankKeeper.SetBalances(ctx, addr1, balances))

	testCases := []struct {
		name      string
		msg       *types.MsgCreateVestingAccount
		expectErr bool
	}{
		{
			name:      "create delayed vesting account",
			msg:       types.NewMsgCreateVestingAccount(addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin("test", 100)), ctx.BlockTime().Unix()+10000, true),
			expectErr: false,
		},
		{
			name:      "create continuous vesting account",
			msg:       types.NewMsgCreateVestingAccount(addr1, addr3, sdk.NewCoins(sdk.NewInt64Coin("test", 100)), ctx.BlockTime().Unix()+10000, false),
			expectErr: false,
		},
		{
			name:      "continuous vesting account already exists",
			msg:       types.NewMsgCreateVestingAccount(addr1, addr3, sdk.NewCoins(sdk.NewInt64Coin("test", 100)), ctx.BlockTime().Unix()+10000, false),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			res, err := suite.handler(ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				err := sdk.ValidateAccAddress(tc.msg.ToAddress)
				suite.Require().NoError(err)
				accI := suite.app.AccountKeeper.GetAccount(ctx, sdk.AccAddress(tc.msg.ToAddress))
				suite.Require().NotNil(accI)

				if tc.msg.Delayed {
					acc, ok := accI.(*types.DelayedVestingAccount)
					suite.Require().True(ok)
					suite.Require().Equal(tc.msg.Amount, acc.GetVestingCoins(ctx.BlockTime()))
				} else {
					acc, ok := accI.(*types.ContinuousVestingAccount)
					suite.Require().True(ok)
					suite.Require().Equal(tc.msg.Amount, acc.GetVestingCoins(ctx.BlockTime()))
				}
			}
		})
	}
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
