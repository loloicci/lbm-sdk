//+build ledger test_ledger_mock

package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/crypto/keyring"
	"github.com/line/lbm-sdk/testutil"
	sdk "github.com/line/lbm-sdk/types"
)

func Test_runAddCmdLedgerWithCustomCoinType(t *testing.T) {
	config := sdk.GetConfig()

	bech32PrefixAccAddr := "terra"
	bech32PrefixAccPub := "terrapub"
	bech32PrefixValAddr := "terravaloper"
	bech32PrefixValPub := "terravaloperpub"
	bech32PrefixConsAddr := "terravalcons"
	bech32PrefixConsPub := "terravalconspub"

	config.SetCoinType(330)
	config.SetFullFundraiserPath("44'/330'/0'/0/0")
	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)

	cmd := AddKeyCommand()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())

	// Prepare a keybase
	kbHome := t.TempDir()

	clientCtx := client.Context{}.WithKeyringDir(kbHome)
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"keyname1",
		fmt.Sprintf("--%s=true", flags.FlagUseLedger),
		fmt.Sprintf("--%s=0", flagAccount),
		fmt.Sprintf("--%s=0", flagIndex),
		fmt.Sprintf("--%s=330", flagCoinType),
		fmt.Sprintf("--%s=%s", cli.OutputFlag, OutputFormatText),
		fmt.Sprintf("--%s=%s", flags.FlagKeyAlgorithm, string(hd.Secp256k1Type)),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
	})

	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)
	mockIn.Reset("test1234\ntest1234\n")
	require.NoError(t, cmd.ExecuteContext(ctx))

	// Now check that it has been stored properly
	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
	require.NoError(t, err)
	require.NotNil(t, kb)
	t.Cleanup(func() {
		_ = kb.Delete("keyname1")
	})
	mockIn.Reset("test1234\n")
	key1, err := kb.Key("keyname1")
	require.NoError(t, err)
	require.NotNil(t, key1)

	require.Equal(t, "keyname1", key1.GetName())
	require.Equal(t, keyring.TypeLedger, key1.GetType())
	require.Equal(t,
		"terrapub1cqmsrdepqvpg7r26nl2pvqqern00m6s9uaax3hauu2rzg8qpjzq9hy6xve7swm8dz8g",
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, key1.GetPubKey()))

	config.SetCoinType(438)
	config.SetFullFundraiserPath("44'/438'/0'/0/0")
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
}

func Test_runAddCmdLedger(t *testing.T) {
	cmd := AddKeyCommand()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())

	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)
	kbHome := t.TempDir()

	clientCtx := client.Context{}.WithKeyringDir(kbHome)
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"keyname1",
		fmt.Sprintf("--%s=true", flags.FlagUseLedger),
		fmt.Sprintf("--%s=%s", cli.OutputFlag, OutputFormatText),
		fmt.Sprintf("--%s=%s", flags.FlagKeyAlgorithm, string(hd.Secp256k1Type)),
		fmt.Sprintf("--%s=%d", flagCoinType, sdk.CoinType),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
	})
	mockIn.Reset("test1234\ntest1234\n")

	require.NoError(t, cmd.ExecuteContext(ctx))

	// Now check that it has been stored properly
	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
	require.NoError(t, err)
	require.NotNil(t, kb)
	t.Cleanup(func() {
		_ = kb.Delete("keyname1")
	})

	mockIn.Reset("test1234\n")
	key1, err := kb.Key("keyname1")
	require.NoError(t, err)
	require.NotNil(t, key1)

	require.Equal(t, "keyname1", key1.GetName())
	require.Equal(t, keyring.TypeLedger, key1.GetType())
	require.Equal(t,
		"linkpub1cqmsrdepq27djm9tzq3sftqsayx95refxk8r5jn0kyshhql9mdjhjx829zlvzygzwr2",
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, key1.GetPubKey()))
}
