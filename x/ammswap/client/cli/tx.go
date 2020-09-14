package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/okex/okexchain/x/ammswap/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "swap",
		Short: "Swap transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		getCmdAddLiquidity(cdc),
		getCmdRemoveLiquidity(cdc),
		getCmdCreateExchange(cdc),
		getCmdTokenSwap(cdc),
	)...)

	return txCmd
}

func getCmdAddLiquidity(cdc *codec.Codec) *cobra.Command {
	// flags
	var minLiquidity string
	var maxBaseAmount string
	var quoteAmount string
	var deadlineDuration string
	cmd := &cobra.Command{
		Use:   "add-liquidity",
		Short: "add liquidity",
		Long: strings.TrimSpace(
			fmt.Sprintf(`add liquidity.

Example:
$ okexchaincli tx swap add-liquidity --max-base-amount 10eth-355 --quote-amount 100okt --min-liquidity 0.001

`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			minLiquidityDec, sdkErr := sdk.NewDecFromStr(minLiquidity)
			if sdkErr != nil {
				return sdkErr
			}
			maxBaseAmountDecCoin, err := sdk.ParseDecCoin(maxBaseAmount)
			if err != nil {
				return err
			}
			quoteAmountDecCoin, err := sdk.ParseDecCoin(quoteAmount)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(deadlineDuration)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(duration).Unix()
			msg := types.NewMsgAddLiquidity(minLiquidityDec, maxBaseAmountDecCoin, quoteAmountDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&minLiquidity, "min-liquidity", "l", "", "Minimum number of sender will mint if total pool token supply is greater than 0")
	cmd.Flags().StringVarP(&maxBaseAmount, "max-base-amount", "", "", "Maximum number of base amount deposited. Deposits max amount if total pool token supply is 0. For example \"100xxb\"")
	cmd.Flags().StringVarP(&quoteAmount, "quote-amount", "q", "", "The number of quote amount. For example \"100okb\"")
	cmd.Flags().StringVarP(&deadlineDuration, "deadline-duration", "d", "30s", "Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	return cmd
}

func getCmdRemoveLiquidity(cdc *codec.Codec) *cobra.Command {
	// flags
	var liquidity string
	var minBaseAmount string
	var minQuoteAmount string
	var deadlineDuration string
	cmd := &cobra.Command{
		Use:   "remove-liquidity",
		Short: "remove liquidity",
		Long: strings.TrimSpace(
			fmt.Sprintf(`remove liquidity.

Example:
$ okexchaincli tx swap remove-liquidity --liquidity 1 --min-base-amount 10eth-355 --min-quote-amount 1okt

`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			liquidityDec, sdkErr := sdk.NewDecFromStr(liquidity)
			if sdkErr != nil {
				return sdkErr
			}
			minBaseAmountDecCoin, err := sdk.ParseDecCoin(minBaseAmount)
			if err != nil {
				return err
			}
			minQuoteAmountDecCoin, err := sdk.ParseDecCoin(minQuoteAmount)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(deadlineDuration)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(duration).Unix()
			msg := types.NewMsgRemoveLiquidity(liquidityDec, minBaseAmountDecCoin, minQuoteAmountDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&liquidity, "liquidity", "l", "", "Liquidity amount of sender will burn")
	cmd.Flags().StringVarP(&minBaseAmount, "min-base-amount", "", "", "Minimum number of base amount withdrawn")
	cmd.Flags().StringVarP(&minQuoteAmount, "min-quote-amount", "q", "", "Minimum number of quote amount withdrawn")
	cmd.Flags().StringVarP(&deadlineDuration, "deadline-duration", "d", "30s", "Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	return cmd
}

func getCmdCreateExchange(cdc *codec.Codec) *cobra.Command {
	// flags
	var token string
	cmd := &cobra.Command{
		Use:   "create-pair",
		Short: "create token pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`create token pair.

Example:
$ okexchaincli tx swap create-pair --token eth-355 --fees 0.01okt 

`),
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateExchange(token, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Create an AMM swap pair by token name")
	return cmd
}

func getCmdTokenSwap(cdc *codec.Codec) *cobra.Command {
	// flags
	var soldTokenAmount string
	var minBoughtTokenAmount string
	var deadline string
	var recipient string
	cmd := &cobra.Command{
		Use:   "token",
		Short: "swap token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`swap token.

Example:
$ okexchaincli tx swap token --sell-amount 1eth-355 --min-buy-amount 60okt

`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			soldTokenAmount, err := sdk.ParseDecCoin(soldTokenAmount)
			if err != nil {
				return err
			}
			minBoughtTokenAmount, err := sdk.ParseDecCoin(minBoughtTokenAmount)
			if err != nil {
				return err
			}
			dur, err := time.ParseDuration(deadline)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(dur).Unix()
			var recip sdk.AccAddress
			if recipient == "" {
				recip = cliCtx.FromAddress
			} else {
				recip, err = sdk.AccAddressFromBech32(recipient)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgTokenToToken(soldTokenAmount, minBoughtTokenAmount,
				deadline, recip, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&soldTokenAmount, "sell-amount", "", "",
		"Amount expected to sell")
	cmd.Flags().StringVarP(&minBoughtTokenAmount, "min-buy-amount", "", "",
		"Minimum amount expected to buy")
	cmd.Flags().StringVarP(&recipient, "recipient", "", "",
		"The address to receive the amount bought")
	cmd.Flags().StringVarP(&deadline, "deadline", "", "100s",
		"Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	return cmd
}
