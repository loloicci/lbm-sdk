package cli

import (
	"fmt"
	"strings"

	"github.com/line/lbm-sdk/x/feegrant/types"
	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	feegrantQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the feegrant module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	feegrantQueryCmd.AddCommand(
		GetCmdQueryFeeGrant(),
		GetCmdQueryFeeGrants(),
	)

	return feegrantQueryCmd
}

// GetCmdQueryFeeGrant returns cmd to query for a grant between granter and grantee.
func GetCmdQueryFeeGrant() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [granter] [grantee]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single grant",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a grant. 
You can find the fee-grant of a granter and grantee.

Example:
$ %s query feegrant grant [granter] [grantee]
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			err := sdk.ValidateAccAddress(args[0])
			if err != nil {
				return err
			}
			granterAddr := sdk.AccAddress(args[0])

			err = sdk.ValidateAccAddress(args[1])
			if err != nil {
				return err
			}
			granteeAddr := sdk.AccAddress(args[1])

			res, err := queryClient.Allowance(
				cmd.Context(),
				&types.QueryAllowanceRequest{
					Granter: granterAddr.String(),
					Grantee: granteeAddr.String(),
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Allowance)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryFeeGrants returns cmd to query for all grants for a grantee.
func GetCmdQueryFeeGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grants [grantee]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all grants of a grantee",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Queries all the grants for a grantee address.

Example:
$ %s query feegrant grants [grantee]
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			err := sdk.ValidateAccAddress(args[0])
			if err != nil {
				return err
			}
			granteeAddr := sdk.AccAddress(args[0])

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Allowances(
				cmd.Context(),
				&types.QueryAllowancesRequest{
					Grantee:    granteeAddr.String(),
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "grants")

	return cmd
}
