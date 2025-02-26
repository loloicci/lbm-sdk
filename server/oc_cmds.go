package server

// DONTCOVER

import (
	"fmt"
	"strings"

	ostcmd "github.com/line/ostracon/cmd/ostracon/commands"
	"github.com/line/ostracon/libs/cli"
	"github.com/line/ostracon/p2p"
	pvm "github.com/line/ostracon/privval"
	ostversion "github.com/line/ostracon/version"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/line/lbm-sdk/codec"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
)

// ShowNodeIDCmd - ported from Ostracon, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}

			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidatorCmd - ported from Ostracon, show this node's validator info
func ShowValidatorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's ostracon validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valPubKey, err := privValidator.GetPubKey()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString(cli.OutputFlag)
			if strings.ToLower(output) == "json" {
				return printlnJSON(valPubKey)
			}

			pubkey, err := cryptocodec.FromOcPubKeyInterface(valPubKey)
			if err != nil {
				return err
			}
			pubkeyBech32, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey)
			if err != nil {
				return err
			}

			fmt.Println(pubkeyBech32)
			return nil
		},
	}

	cmd.Flags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")
	return &cmd
}

// ShowAddressCmd - show this node's validator address
func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's ostracon validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valConsAddr := sdk.BytesToConsAddress(privValidator.GetAddress())

			output, _ := cmd.Flags().GetString(cli.OutputFlag)
			if strings.ToLower(output) == "json" {
				return printlnJSON(valConsAddr)
			}

			fmt.Println(valConsAddr.String())
			return nil
		},
	}

	cmd.Flags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

// VersionCmd prints ostracon and ABCI version numbers.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print ostracon libraries' version",
		Long: `Print protocols' and libraries' version numbers
against which this app has been compiled.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(&struct {
				Ostracon      string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				Ostracon:      ostversion.OCCoreSemVer,
				ABCI:          ostversion.ABCIVersion,
				BlockProtocol: ostversion.BlockProtocol,
				P2PProtocol:   ostversion.P2PProtocol,
			})
			if err != nil {
				return err
			}

			fmt.Println(string(bs))
			return nil
		},
	}
}

func printlnJSON(v interface{}) error {
	cdc := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(cdc)

	marshalled, err := cdc.MarshalJSON(v)
	if err != nil {
		return err
	}

	fmt.Println(string(marshalled))
	return nil
}

// UnsafeResetAllCmd - extension of the ostracon command, resets initialization
func UnsafeResetAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-reset-all",
		Short: "Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			ostcmd.ResetAll(cfg.DBDir(), cfg.P2P.AddrBookFile(), cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile(),
				cfg.PrivKeyType, serverCtx.Logger)
			return nil
		},
	}
}
