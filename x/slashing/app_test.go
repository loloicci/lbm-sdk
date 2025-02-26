package slashing_test

import (
	"errors"
	"testing"

	abci "github.com/line/ostracon/abci/types"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"github.com/line/lbm-sdk/x/slashing/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.BytesToAccAddress(priv1.PubKey().Address())

	valKey  = ed25519.GenPrivKey()
	valAddr = sdk.BytesToAccAddress(valKey.PubKey().Address())
)

func checkValidator(t *testing.T, app *simapp.SimApp, _ sdk.AccAddress, expFound bool) stakingtypes.Validator {
	ctxCheck := app.BaseApp.NewContext(true, ocproto.Header{})
	validator, found := app.StakingKeeper.GetValidator(ctxCheck, addr1.ToValAddress())
	require.Equal(t, expFound, found)
	return validator
}

func checkValidatorSigningInfo(t *testing.T, app *simapp.SimApp, addr sdk.ConsAddress, expFound bool) types.ValidatorSigningInfo {
	ctxCheck := app.BaseApp.NewContext(true, ocproto.Header{})
	signingInfo, found := app.SlashingKeeper.GetValidatorSigningInfo(ctxCheck, addr)
	require.Equal(t, expFound, found)
	return signingInfo
}

func TestSlashingMsgs(t *testing.T) {
	genTokens := sdk.TokensFromConsensusPower(42)
	bondTokens := sdk.TokensFromConsensusPower(10)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{
		Address: addr1.String(),
	}
	accs := authtypes.GenesisAccounts{acc1}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.Coins{genCoin},
		},
	}

	app := simapp.SetupWithGenesisAccounts(accs, balances...)
	simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin})

	description := stakingtypes.NewDescription("foo_moniker", "", "", "", "")
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())

	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(
		addr1.ToValAddress(), valKey.PubKey(), bondCoin, description, commission, sdk.OneInt(),
	)
	require.NoError(t, err)

	header := ocproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, "", []uint64{0}, []uint64{0}, true, true, priv1)
	require.NoError(t, err)
	simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin.Sub(bondCoin)})

	header = ocproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	validator := checkValidator(t, app, addr1, true)
	require.Equal(t, addr1.ToValAddress().String(), validator.OperatorAddress)
	require.Equal(t, stakingtypes.Bonded, validator.Status)
	require.True(sdk.IntEq(t, bondTokens, validator.BondedTokens()))
	unjailMsg := &types.MsgUnjail{ValidatorAddr: addr1.ToValAddress().String()}

	checkValidatorSigningInfo(t, app, valAddr.ToConsAddress(), true)

	// unjail should fail with unknown validator
	header = ocproto.Header{Height: app.LastBlockHeight() + 1}
	_, res, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{unjailMsg}, "", []uint64{0}, []uint64{1}, false, false, priv1)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, errors.Is(types.ErrValidatorNotJailed, err))
}
