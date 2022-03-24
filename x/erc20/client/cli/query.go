package cli

import (
	"fmt"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/x/erc20/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd defines erc20 module queries through the cli
func GetQueryCmd(moduleName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the erc20 module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(flags.GetCommands(
		GetCmdQueryParams(moduleName, cdc),
		GetCmdQueryAllMapping(moduleName, cdc),
	)...)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query all the modifiable parameters of gov proposal",
		Long: strings.TrimSpace(`Query the all the parameters for the governance process:

$ exchaincli query erc20 params
`),
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

//
func GetCmdQueryAllMapping(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all-mapping",
		Short: "Query all mapping of denom and contract",
		Long: strings.TrimSpace(`Query all mapping of denom and contract:

$ exchaincli query erc20 all-mapping
`),
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllMapping)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var mapping []types.TokenMapping
			cdc.MustUnmarshalJSON(bz, &mapping)
			return cliCtx.PrintOutput(mapping)
		},
	}
}
