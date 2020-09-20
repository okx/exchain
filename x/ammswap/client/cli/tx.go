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

// flags
const (
	flagMinLiquidity     = "min-liquidity"
	flagMaxAmountTokenA  = "max-amount-token-a"
	flagAmountTokenB     = "amount-token-b"
	flagDeadlineDuration = "deadline-duration"
	flagLiquidity        = "liquidity"
	flagMinAmountTokenA  = "min-amount-token-a"
	flagMinAmountTokenB  = "min-amount-token-b"
	flagSellAmount       = "sell-amount"
	flagTokenRoute       = "token-route"
	flagMinBuyAmount     = "min-buy-amount"
	flagRecipient        = "recipient"
	flagTokenA           = "token-a"
	flagTokenB           = "token-b"
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
	var maxAmountTokenA string
	var amountTokenB string
	var deadlineDuration string
	cmd := &cobra.Command{
		Use:   "add-liquidity",
		Short: "add liquidity",
		Long: strings.TrimSpace(
			fmt.Sprintf(`add liquidity.

Example:
$ okexchaincli tx swap add-liquidity --max-amount-token-a 10eth-355 --amount-token-b 100btc-366 --min-liquidity 0.001

`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			minLiquidityDec, sdkErr := sdk.NewDecFromStr(minLiquidity)
			if sdkErr != nil {
				return sdkErr
			}
			maxAmountTokenADecCoin, err := sdk.ParseDecCoin(maxAmountTokenA)
			if err != nil {
				return err
			}
			amountTokenBDecCoin, err := sdk.ParseDecCoin(amountTokenB)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(deadlineDuration)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(duration).Unix()
			msg := types.NewMsgAddLiquidity(minLiquidityDec, maxAmountTokenADecCoin, amountTokenBDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&minLiquidity, flagMinLiquidity, "l", "", "Minimum number of sender will mint if total pool token supply is greater than 0")
	cmd.Flags().StringVarP(&maxAmountTokenA, flagMaxAmountTokenA, "", "", "Maximum number of amount deposited. Deposits max amount if total pool token supply is 0. For example \"100xxb\"")
	cmd.Flags().StringVarP(&amountTokenB, flagAmountTokenB, "q", "", "The number of amount. For example \"100okb\"")
	cmd.Flags().StringVarP(&deadlineDuration, flagDeadlineDuration, "d", "30s", "Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	cmd.MarkFlagRequired(flagMinLiquidity)
	cmd.MarkFlagRequired(flagMaxAmountTokenA)
	cmd.MarkFlagRequired(flagAmountTokenB)
	return cmd
}

func getCmdRemoveLiquidity(cdc *codec.Codec) *cobra.Command {
	// flags
	var liquidity string
	var minAmountTokenA string
	var minAmountTokenB string
	var deadlineDuration string
	cmd := &cobra.Command{
		Use:   "remove-liquidity",
		Short: "remove liquidity",
		Long: strings.TrimSpace(
			fmt.Sprintf(`remove liquidity.

Example:
$ okexchaincli tx swap remove-liquidity --liquidity 1 --min-amount-token-a 10eth-355 --min-amount-token-b 1btc-366

`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			liquidityDec, sdkErr := sdk.NewDecFromStr(liquidity)
			if sdkErr != nil {
				return sdkErr
			}
			minAmountTokenADecCoin, err := sdk.ParseDecCoin(minAmountTokenA)
			if err != nil {
				return err
			}
			minAmountTokenBDecCoin, err := sdk.ParseDecCoin(minAmountTokenB)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(deadlineDuration)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(duration).Unix()
			msg := types.NewMsgRemoveLiquidity(liquidityDec, minAmountTokenADecCoin, minAmountTokenBDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&liquidity, flagLiquidity, "l", "", "Liquidity amount of sender will burn")
	cmd.Flags().StringVarP(&minAmountTokenA, flagMinAmountTokenA, "", "", "Minimum number of amount withdrawn")
	cmd.Flags().StringVarP(&minAmountTokenB, flagMinAmountTokenB, "q", "", "Minimum number of amount withdrawn")
	cmd.Flags().StringVarP(&deadlineDuration, flagDeadlineDuration, "d", "30s", "Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	cmd.MarkFlagRequired(flagLiquidity)
	cmd.MarkFlagRequired(flagMinAmountTokenA)
	cmd.MarkFlagRequired(flagMinAmountTokenB)
	return cmd
}

func getCmdCreateExchange(cdc *codec.Codec) *cobra.Command {
	// flags
	var nameTokenA string
	var nameTokenB string
	cmd := &cobra.Command{
		Use:   "create-pair",
		Short: "create token pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`create token pair.

Example:
$ okexchaincli tx swap create-pair --token-a eth-355 --token-b btc-366 --fees 0.01okt 

`),
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateExchange(nameTokenA, nameTokenB, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVar(&nameTokenA, flagTokenA,  "", "the token name is required to create an AMM swap pair")
	cmd.Flags().StringVarP(&nameTokenB, flagTokenB, "q", "", "the token name is required to create an AMM swap pair")
	cmd.MarkFlagRequired(flagTokenA)
	cmd.MarkFlagRequired(flagTokenB)
	return cmd
}

func getCmdTokenSwap(cdc *codec.Codec) *cobra.Command {
	// flags
	var soldTokenAmount string
	var minBoughtTokenAmount string
	var deadline string
	var recipient string
	var tokenRoute []string
	cmd := &cobra.Command{
		Use:   "token",
		Short: "swap token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`swap token.

Example:
$ okexchaincli tx swap token --sell-amount 1eth-355 --min-buy-amount 60btc-366

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
				tokenRoute, deadline, recip, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&soldTokenAmount, flagSellAmount, "", "",
		"Amount expected to sell")
	cmd.Flags().StringVarP(&minBoughtTokenAmount, flagMinBuyAmount, "", "",
		"Minimum amount expected to buy")
	cmd.Flags().StringVarP(&recipient, flagRecipient, "", "",
		"The address to receive the amount bought")
	cmd.Flags().StringVarP(&deadline, flagDeadlineDuration, "", "100s",
		"Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	cmd.Flags().StringSliceVarP(&tokenRoute, flagTokenRoute, "", nil,
		"Intermediate route from sold token to bought token, split with \\',\\',for example \"aab,ccb,ddb\"")
	cmd.MarkFlagRequired(flagSellAmount)
	cmd.MarkFlagRequired(flagMinBuyAmount)

	return cmd
}
