package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/version"
	"github.com/line/lbm-sdk/x/authz/types"
	bank "github.com/line/lbm-sdk/x/bank/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	authorizationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the authz module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authorizationQueryCmd.AddCommand(
		GetCmdQueryAuthorization(),
		GetCmdQueryAuthorizations(),
	)

	return authorizationQueryCmd
}

// GetCmdQueryAuthorizations implements the query authorizations command.
func GetCmdQueryAuthorizations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorizations [granter-addr] [grantee-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "query list of authorizations for a granter-grantee pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query list of authorizations for a granter-grantee pair:
Example:
$ %s query %s authorizations link1skj.. link1skjwj..
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			err = sdk.ValidateAccAddress(args[0])
			if err != nil {
				return err
			}
			granterAddr := sdk.AccAddress(args[0])

			err = sdk.ValidateAccAddress(args[1])
			if err != nil {
				return err
			}
			granteeAddr := sdk.AccAddress(args[1])

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Authorizations(
				context.Background(),
				&types.QueryAuthorizationsRequest{
					Granter:    granterAddr.String(),
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
	flags.AddPaginationFlagsToCmd(cmd, "authorizations")
	return cmd
}

// GetCmdQueryAuthorization implements the query authorization command.
func GetCmdQueryAuthorization() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorization [granter-addr] [grantee-addr] [msg-type]",
		Args:  cobra.ExactArgs(3),
		Short: "query authorization for a granter-grantee pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query authorization for a granter-grantee pair that matches the given msg-type:
Example:
$ %s query %s authorization link1skjw.. link1skjwj.. %s
`, version.AppName, types.ModuleName, bank.SendAuthorization{}.MethodName()),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			err = sdk.ValidateAccAddress(args[0])
			if err != nil {
				return err
			}
			granter := sdk.AccAddress(args[0])

			err = sdk.ValidateAccAddress(args[1])
			if err != nil {
				return err
			}
			grantee := sdk.AccAddress(args[1])

			msgAuthorized := args[2]

			res, err := queryClient.Authorization(
				context.Background(),
				&types.QueryAuthorizationRequest{
					Granter:    granter.String(),
					Grantee:    grantee.String(),
					MethodName: msgAuthorized,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Authorization)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
