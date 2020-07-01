package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/token/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the token module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(client.GetCommands(
		getCmdQueryParams(queryRoute, cdc),
		getCmdTokenInfo(queryRoute, cdc),
		//getAccountCmd(queryRoute, cdc),
	)...)

	return queryCmd
}

// getCmdTokenInfo queries token info by address
func getCmdTokenInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	var owner string
	cmd := &cobra.Command{
		Use:   "info [<symbol>]",
		Short: "query token info by symbol",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			switch {
			case len(args) == 1:
				symbol := args[0]
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/info/%s", queryRoute, symbol), nil)
				if err != nil {
					fmt.Printf("token does not exist - %s %s\n", symbol, err.Error())
					return nil
				}

				var token types.TokenInfo
				cdc.MustUnmarshalJSON(res, &token)
				return cliCtx.PrintOutput(token)
			case owner != "":
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tokens/%s", queryRoute, owner), nil)
				if err != nil {
					fmt.Printf("Invalid owner address - %s %s\n", owner, err.Error())
					return nil
				}

				var tokens types.Tokens
				cdc.MustUnmarshalJSON(res, &tokens)
				return cliCtx.PrintOutput(tokens)
			default:
				fmt.Println("At least a [<symbol>] arg or an --owner flag need to be set")
				err := cmd.Help()
				if err != nil {
					return err
				}
				return nil
			}
		},
	}
	cmd.Flags().StringVarP(&owner, "owner", "", "", "Get all the tokens information belong to")

	return cmd
}

// getCmdQueryParams implements the query params command.
func getCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the token process",
		Long: strings.TrimSpace(`Query the all the parameters for the governance process:

$ okchaincli query token params
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

// just for the object of []string could be inputted into cliCtx.PrintOutput(...)
type Strings []string

func (strs Strings) String() string {
	return strings.Join(strs, "\n")
}

func getAccount(ctx *context.CLIContext, address []byte) (auth.Account, error) {
	bz, err := ctx.Codec.MarshalJSON(auth.NewQueryAccountParams(address))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", auth.StoreKey, auth.QueryAccount)
	res, _, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	var account auth.Account
	if err := ctx.Codec.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}
	return account, nil
}
