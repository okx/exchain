package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
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
			GetCmdQueryBuyAmount(queryRoute, cdc),
		)...,
	)

	return swapQueryCmd
}

//GetCmdSwapTokenPair query exchange with token name
func GetCmdSwapTokenPair(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool [base-token] [quote-token]",
		Short: "Query pool info by token name",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query pool info by token name.

Example:
$ okexchaincli query swap pool eth-355

`),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			baseToken := args[0]
			quoteToken := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QuerySwapTokenPair, baseToken, quoteToken), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdQueryBuyAmount queries amount of base/quote token by the given amount of quote/base token
func GetCmdQueryBuyAmount(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "amount [token-to-sell] [token-name-to-buy]",
		Short: "Query how many token returned by the given amount of token to sell",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query how many base token returned by the given amount of quote token.

Example:
$ %s query swap amount 100eth-245 xxb`, version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sellToken, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}
			params := types.QueryBuyAmountParams{
				SoldToken:    sellToken,
				TokenToBuy: args[1],
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			tp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBuyAmount), bz)
			if err != nil {
				return err
			}

			var buyAmt sdk.Dec
			cdc.MustUnmarshalJSON(tp, &buyAmt)

			return cliCtx.PrintOutput(buyAmt)
		},
	}
}

// GetCmdQueryParams queries the parameters of the AMM swap system
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
		Use:   "pools",
		Short: "Query infomation of all pools",
		Long: 	strings.TrimSpace(
			fmt.Sprintf(`Query infomation of all pools.
Example:
$ okexchaincli query swap pools
`),
		),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySwapTokenPairs), nil)
			if err != nil {
				return err
			}
			if res == nil || len(res) == 0 || string(res) == "null" {
				fmt.Println("empty SwapTokenPairs")
			}else {
				fmt.Println(string(res))
			}

			return nil
		},
	}
}


//GetCmdRedeemableAssets query redeemable assets by specifying the number of lpt
func GetCmdRedeemableAssets(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "redeemable-assets [base-token] [quote-token] [pool-token-amount]",
		Short: "Query redeemable assets by specifying pool token amount",
		Long: 	strings.TrimSpace(
			fmt.Sprintf(`Query redeemable assets by specifying pool token amount.
Example:
$ okexchaincli query swap redeemable-assets eth xxb 1
`),
		),
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			baseTokenName := args[0]
			quoteTokenName := args[1]
			liquidity := args[2]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s/%s", queryRoute, types.QueryRedeemableAssets, baseTokenName, quoteTokenName, liquidity), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}