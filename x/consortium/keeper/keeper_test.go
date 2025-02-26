package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
)

var (
	delPk   = ed25519.GenPrivKey().PubKey()
	delAddr = sdk.BytesToAccAddress(delPk.Address())
	valAddr = delAddr.ToValAddress()
)

func TestCleanup(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.ConsortiumKeeper

	// add auths
	auth :=	&types.ValidatorAuth{
		OperatorAddress: string(valAddr),
		CreationAllowed: true,
	}
	require.NoError(t, k.SetValidatorAuth(ctx, auth))

	// cleanup
	k.Cleanup(ctx)
	require.Equal(t, []*types.ValidatorAuth{}, k.GetValidatorAuths(ctx))
}
