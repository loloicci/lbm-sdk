package keeper_test

import (
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/slashing/testslashing"
	"github.com/line/lbm-sdk/x/staking"
	"github.com/line/lbm-sdk/x/staking/teststaking"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func TestUnJailNotBonded(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	p := app.StakingKeeper.GetParams(ctx)
	p.MaxValidators = 5
	app.StakingKeeper.SetParams(ctx, p)

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 6, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(6)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create max (5) validators all with the same power
	for i := uint32(0); i < p.MaxValidators; i++ {
		addr, val := valAddrs[i], pks[i]
		tstaking.CreateValidatorWithValPower(addr, val, 100, true)
	}

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// create a 6th validator with less power than the cliff validator (won't be bonded)
	addr, val := valAddrs[5], pks[5]
	amt := sdk.TokensFromConsensusPower(50)
	msg := tstaking.CreateValidatorMsg(addr, val, amt)
	msg.MinSelfDelegation = amt
	tstaking.Handle(msg, true)

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	tstaking.CheckValidator(addr, stakingtypes.Unbonded, false)

	// unbond below minimum self-delegation
	require.Equal(t, p.BondDenom, tstaking.Denom)
	tstaking.Undelegate(addr.ToAccAddress(), addr, sdk.TokensFromConsensusPower(1), true)

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// verify that validator is jailed
	tstaking.CheckValidator(addr, -1, true)

	// verify we cannot unjail (yet)
	require.Error(t, app.SlashingKeeper.Unjail(ctx, addr))

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// bond to meet minimum self-delegation
	tstaking.DelegateWithPower(addr.ToAccAddress(), addr, 1)

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// verify we can immediately unjail
	require.NoError(t, app.SlashingKeeper.Unjail(ctx, addr))

	tstaking.CheckValidator(addr, -1, false)
}

// Test a new validator entering the validator set
// Ensure that SigningInfo.VoterSetCounter is set correctly
// and that they are not immediately jailed
func TestHandleNewValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(1)
	addr, val := valAddrs[0], pks[0]
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(app.SlashingKeeper.SignedBlocksWindow(ctx) + 1)

	// Validator created
	amt := tstaking.CreateValidatorWithValPower(addr, val, 100, true)

	staking.EndBlocker(ctx, app.StakingKeeper)
	require.Equal(
		t, app.BankKeeper.GetAllBalances(ctx, addr.ToAccAddress()),
		sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, app.StakingKeeper.Validator(ctx, addr).GetBondedTokens())

	// Now a validator, for two blocks
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	ctx = ctx.WithBlockHeight(app.SlashingKeeper.SignedBlocksWindow(ctx) + 2)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)

	info, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.BytesToConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(2), info.VoterSetCounter)
	require.Equal(t, int64(1), info.MissedBlocksCounter)
	require.Equal(t, time.Unix(0, 0).UTC(), info.JailedUntil)

	// validator should be bonded still, should not have been jailed or slashed
	validator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Bonded, validator.GetStatus())
	bondPool := app.StakingKeeper.GetBondedPool(ctx)
	expTokens := sdk.TokensFromConsensusPower(100)
	require.True(t, expTokens.Equal(app.BankKeeper.GetBalance(ctx, bondPool.GetAddress(), app.StakingKeeper.BondDenom(ctx)).Amount))
}

// Test a jailed validator being "down" twice
// Ensure that they're only slashed once
func TestHandleAlreadyJailed(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(1)
	addr, val := valAddrs[0], pks[0]
	power := int64(100)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	amt := tstaking.CreateValidatorWithValPower(addr, val, power, true)

	staking.EndBlocker(ctx, app.StakingKeeper)

	// 1000 first blocks OK
	height := int64(0)
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// 501 blocks missed
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx)+(app.SlashingKeeper.SignedBlocksWindow(ctx)-app.SlashingKeeper.MinSignedPerWindow(ctx))+1; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	// validator should have been jailed and slashed
	validator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Unbonding, validator.GetStatus())

	// validator should have been slashed
	resultingTokens := amt.Sub(sdk.TokensFromConsensusPower(1))
	require.Equal(t, resultingTokens, validator.GetTokens())

	// another block missed
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)

	// validator should not have been slashed twice
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, resultingTokens, validator.GetTokens())
}

// Test a validator dipping in and out of the validator set
// Ensure that missed blocks are tracked correctly and that
// the voter set counter of the signing info is reset correctly
func TestValidatorDippingInAndOut(t *testing.T) {

	// initial setup
	// TestParams set the SignedBlocksWindow to 1000 and MaxMissedBlocksPerWindow to 500
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 1
	app.StakingKeeper.SetParams(ctx, params)
	power := int64(100)

	pks := simapp.CreateTestPubKeys(3)
	simapp.AddTestAddrsFromPubKeys(app, ctx, pks, sdk.TokensFromConsensusPower(200))

	addr, val := pks[0].Address(), pks[0]
	consAddr := sdk.BytesToConsAddress(addr)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	valAddr := sdk.BytesToValAddress(addr)

	tstaking.CreateValidatorWithValPower(valAddr, val, power, true)
	staking.EndBlocker(ctx, app.StakingKeeper)

	// 100 first blocks OK
	height := int64(0)
	for ; height < int64(100); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// kick first validator out of validator set
	tstaking.CreateValidatorWithValPower(sdk.BytesToValAddress(pks[1].Address()), pks[1], 101, true)
	validatorUpdates := staking.EndBlocker(ctx, app.StakingKeeper)
	require.Equal(t, 2, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Unbonding, false)

	// 600 more blocks happened
	height = int64(700)
	ctx = ctx.WithBlockHeight(height)

	// validator added back in
	tstaking.DelegateWithPower(sdk.BytesToAccAddress(pks[2].Address()), sdk.BytesToValAddress(pks[0].Address()), 50)

	validatorUpdates = staking.EndBlocker(ctx, app.StakingKeeper)
	require.Equal(t, 2, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Bonded, false)
	newPower := int64(150)

	// validator misses 501 blocks exceeding the liveness threshold
	latest := height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, false)
	}

	// 398 more blocks happened
	latest = height
	for ; height < latest+398; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, true)
	}

	// shouldn't be jailed/kicked yet because it have not joined to vote set 1000 times
	// 100 times + (kicked) + 501 times + 398 times = 999 times
	tstaking.CheckValidator(valAddr, stakingtypes.Bonded, false)

	// another block happened
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, true)
	height++

	// should now be jailed & kicked
	staking.EndBlocker(ctx, app.StakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Unbonding, true)

	// check all the signing information
	signInfo, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, int64(0), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(1000), signInfo.VoterSetCounter)
	// array should be cleared
	for offset := int64(0); offset < app.SlashingKeeper.SignedBlocksWindow(ctx); offset++ {
		missed := app.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, consAddr, offset)
		require.False(t, missed)
	}

	// some blocks pass
	height = int64(5000)
	ctx = ctx.WithBlockHeight(height)

	// validator rejoins and starts signing again
	app.StakingKeeper.Unjail(ctx, consAddr)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, true)
	height++

	// validator should not be kicked since we reset counter/array when it was jailed
	staking.EndBlocker(ctx, app.StakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Bonded, false)

	// validator misses 501 blocks
	latest = height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, false)
	}

	// validator should now be jailed & kicked
	staking.EndBlocker(ctx, app.StakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Unbonding, true)

	// some blocks pass
	height = int64(10000)
	ctx = ctx.WithBlockHeight(height)

	// validator rejoins and starts signing again
	app.StakingKeeper.Unjail(ctx, consAddr)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, true)
	height++

	// validator should not be kicked since we reset counter/array when it was jailed
	staking.EndBlocker(ctx, app.StakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Bonded, false)

	// 1000 blocks happened
	latest = height
	for ; height < latest+1000; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, true)
	}

	// validator misses 500 blocks
	latest = height
	for ; height < latest+500; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, false)
	}
	tstaking.CheckValidator(valAddr, stakingtypes.Bonded, false)

	// validator misses another block
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), newPower, false)
	height++

	// validator should now be jailed & kicked
	staking.EndBlocker(ctx, app.StakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Unbonding, true)

}
