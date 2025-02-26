package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/crypto/keyring"
	"github.com/line/lbm-sdk/crypto/keys/multisig"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	"github.com/line/lbm-sdk/testutil"
	sdk "github.com/line/lbm-sdk/types"
)

func Test_multiSigKey_Properties(t *testing.T) {
	tmpKey1 := secp256k1.GenPrivKeyFromSecret([]byte("mySecret"))
	pk := multisig.NewLegacyAminoPubKey(
		1,
		[]cryptotypes.PubKey{tmpKey1.PubKey()},
	)
	tmp := keyring.NewMultiInfo("myMultisig", pk)

	require.Equal(t, "myMultisig", tmp.GetName())
	require.Equal(t, keyring.TypeMulti, tmp.GetType())
	require.Equal(t, "BDF0C827D34CA39919C7688EB5A95383C60B3471", tmp.GetPubKey().Address().String())
	acc := tmp.GetAddress()
	addrBytes, _ := sdk.AccAddressToBytes(acc.String())
	require.Equal(t, "link1hhcvsf7nfj3ejxw8dz8tt22ns0rqkdr3rrh7xy", sdk.MustBech32ifyAddressBytes("link", addrBytes))
}

func Test_showKeysCmd(t *testing.T) {
	cmd := ShowKeysCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "false", cmd.Flag(FlagAddress).DefValue)
	require.Equal(t, "false", cmd.Flag(FlagPublicKey).DefValue)
}

func Test_runShowCmd(t *testing.T) {
	cmd := ShowKeysCmd()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())
	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)

	kbHome := t.TempDir()
	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
	require.NoError(t, err)

	clientCtx := client.Context{}.WithKeyring(kb)
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{"invalid"})
	require.EqualError(t, cmd.ExecuteContext(ctx), "invalid is not a valid name or address: decoding bech32 failed: invalid bech32 string length 7")

	cmd.SetArgs([]string{"invalid1", "invalid2"})
	require.EqualError(t, cmd.ExecuteContext(ctx), "invalid1 is not a valid name or address: decoding bech32 failed: invalid index of 1")

	fakeKeyName1 := "runShowCmd_Key1"
	fakeKeyName2 := "runShowCmd_Key2"

	t.Cleanup(func() {
		kb.Delete("runShowCmd_Key1")
		kb.Delete("runShowCmd_Key2")
	})

	path := hd.NewFundraiserParams(1, sdk.CoinType, 0).String()
	_, err = kb.NewAccount(fakeKeyName1, testutil.TestMnemonic, "", path, hd.Secp256k1)
	require.NoError(t, err)

	path2 := hd.NewFundraiserParams(1, sdk.CoinType, 1).String()
	_, err = kb.NewAccount(fakeKeyName2, testutil.TestMnemonic, "", path2, hd.Secp256k1)
	require.NoError(t, err)

	// Now try single key
	cmd.SetArgs([]string{
		fakeKeyName1,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=", FlagBechPrefix),
	})
	require.EqualError(t, cmd.ExecuteContext(ctx), "invalid Bech32 prefix encoding provided: ")

	cmd.SetArgs([]string{
		fakeKeyName1,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=%s", FlagBechPrefix, sdk.PrefixAccount),
	})

	// try fetch by name
	require.NoError(t, cmd.ExecuteContext(ctx))

	// try fetch by addr
	info, err := kb.Key(fakeKeyName1)
	cmd.SetArgs([]string{
		info.GetAddress().String(),
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=%s", FlagBechPrefix, sdk.PrefixAccount),
	})

	require.NoError(t, err)
	require.NoError(t, cmd.ExecuteContext(ctx))

	// Now try multisig key - set bech to acc
	cmd.SetArgs([]string{
		fakeKeyName1, fakeKeyName2,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=%s", FlagBechPrefix, sdk.PrefixAccount),
		fmt.Sprintf("--%s=0", flagMultiSigThreshold),
	})
	require.EqualError(t, cmd.ExecuteContext(ctx), "threshold must be a positive integer")

	cmd.SetArgs([]string{
		fakeKeyName1, fakeKeyName2,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=%s", FlagBechPrefix, sdk.PrefixAccount),
		fmt.Sprintf("--%s=2", flagMultiSigThreshold),
	})
	require.NoError(t, cmd.ExecuteContext(ctx))

	// Now try multisig key - set bech to acc + threshold=2
	cmd.SetArgs([]string{
		fakeKeyName1, fakeKeyName2,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=acc", FlagBechPrefix),
		fmt.Sprintf("--%s=true", FlagDevice),
		fmt.Sprintf("--%s=2", flagMultiSigThreshold),
	})
	require.EqualError(t, cmd.ExecuteContext(ctx), "the device flag (-d) can only be used for accounts stored in devices")

	cmd.SetArgs([]string{
		fakeKeyName1, fakeKeyName2,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=val", FlagBechPrefix),
		fmt.Sprintf("--%s=true", FlagDevice),
		fmt.Sprintf("--%s=2", flagMultiSigThreshold),
	})
	require.EqualError(t, cmd.ExecuteContext(ctx), "the device flag (-d) can only be used for accounts")

	cmd.SetArgs([]string{
		fakeKeyName1, fakeKeyName2,
		fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--%s=val", FlagBechPrefix),
		fmt.Sprintf("--%s=true", FlagDevice),
		fmt.Sprintf("--%s=2", flagMultiSigThreshold),
		fmt.Sprintf("--%s=true", FlagPublicKey),
	})
	require.EqualError(t, cmd.ExecuteContext(ctx), "the device flag (-d) can only be used for addresses not pubkeys")
}

func Test_validateMultisigThreshold(t *testing.T) {
	type args struct {
		k     int
		nKeys int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"zeros", args{0, 0}, true},
		{"1-0", args{1, 0}, true},
		{"1-1", args{1, 1}, false},
		{"1-2", args{1, 1}, false},
		{"1-2", args{2, 1}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMultisigThreshold(tt.args.k, tt.args.nKeys); (err != nil) != tt.wantErr {
				t.Errorf("validateMultisigThreshold() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getBechKeyOut(t *testing.T) {
	type args struct {
		bechPrefix string
	}
	tests := []struct {
		name    string
		args    args
		want    bechKeyOutFn
		wantErr bool
	}{
		{"empty", args{""}, nil, true},
		{"wrong", args{"???"}, nil, true},
		{"acc", args{sdk.PrefixAccount}, keyring.Bech32KeyOutput, false},
		{"val", args{sdk.PrefixValidator}, keyring.Bech32ValKeyOutput, false},
		{"cons", args{sdk.PrefixConsensus}, keyring.Bech32ConsKeyOutput, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBechKeyOut(tt.args.bechPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBechKeyOut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				require.NotNil(t, got)
			}

			// TODO: Still not possible to compare functions
			// Maybe in next release: https://github.com/stretchr/testify/issues/182
			//if &got != &tt.want {
			//	t.Errorf("getBechKeyOut() = %v, want %v", got, tt.want)
			//}
		})
	}
}
