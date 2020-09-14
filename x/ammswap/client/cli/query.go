package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/spf13/cobra"
	"strings"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group swap queries under a subcommand
	swapQueryCmd := &cobra.Command{
		Use:                        "swap",
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	swapQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdSwapTokenPair(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
			GetCmdAllSwapTokenPairs(queryRoute, cdc),
			GetCmdRedeemableAssets(queryRoute, cdc),
		)...,
	)

	return swapQueryCmd
}

//GetCmdSwapTokenPair query exchange with token name
func GetCmdSwapTokenPair(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-info [token]",
		Short: "Query pool info by token name",
		Long: 	strings.TrimSpace(
			fmt.Sprintf(`Query pool info by token name.

Example:
$ okexchaincli query swap pool-info eth-355

`),
		),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenName := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QuerySwapTokenPair, tokenName), nil)
			if err != nil {
				fmt.Printf("token pool - %s doesn't exist. error:%s \n", tokenName, err.Error())
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the AMM swap system",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the parameters of the AMM swap system.

Example:
$ %s query swap params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params", queryRoute), nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(tp, &params)

			return cliCtx.PrintOutput(params)
		},
	}
}


//GetCmdAllSwapTokenPairs lists all info of pools
func GetCmdAllSwapTokenPairs(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all-pool-info [token]",
		Short: "List all info of pools",
		Long: 	strings.TrimSpace(
			fmt.Sprintf(`List all info of pools.
Example:
$ okexchaincli query swap all-pool-info
`),
		),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySwapTokenPairs), nil)
			if err != nil {
				fmt.Printf("query all SwapTokenPairs failed. error:%s \n", err.Error())
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
}


//GetCmdRedeemableAssets query redeemable assets by specifying the number of lpt
func GetCmdRedeemableAssets(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "redeemable-assets [the name of base token] [the number ofliquidity pool token]",
		Short: "Query redeemable assets by specifying the number of lpt",
		Long: 	strings.TrimSpace(
			fmt.Sprintf(`Query redeemable assets by specifying the number of lpt.
Example:
$ okexchaincli query swap redeemable-assets eth-355 1
`),
		),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			baseTokenName := args[0]
			liquidity := args[1]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryRedeemableAssets, baseTokenName, liquidity), nil)
			if err != nil {
				fmt.Printf("query redeemable assets failed. error:%s \n", err.Error())
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
}