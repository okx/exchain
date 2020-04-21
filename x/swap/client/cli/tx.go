package cli

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/okex/okchain/x/swap/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "swap",
		Short: "swap module",
	}

	txCmd.AddCommand(client.PostCommands(
		getCmdAddLiquidity(cdc),
		getCmdCreateExchange(cdc),
		getCmdTokenOKTSwap(cdc),
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			minLiquidityDec, err := sdk.NewDecFromStr(minLiquidity)
			if err != nil {
				return err
			}
			maxBaseAmountDecCoin, err2 := sdk.ParseDecCoin(maxBaseAmount)
			if err2 != nil {
				return err2
			}
			quoteAmountDecCoin, err2 := sdk.ParseDecCoin(quoteAmount)
			if err2 != nil {
				return err2
			}
			duration, err3 := time.ParseDuration(deadlineDuration)
			if err3 != nil {
				return err3
			}
			deadline := time.Now().Add(duration).Unix()
			msg := types.NewMsgAddLiquidity(minLiquidityDec, maxBaseAmountDecCoin, quoteAmountDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&minLiquidity, "min-liquidity", "l", "", "Minimum number of sender will mint if total pool token supply is greater than 0")
	cmd.Flags().StringVarP(&maxBaseAmount, "max-base-tokens", "", "", "Maximum number of base tokens deposited. Deposits max amount if total pool token supply is 0. For example \"100xxb\"")
	cmd.Flags().StringVarP(&quoteAmount, "quote-tokens", "q", "", "The number of quote tokens. For example \"100okb\"")
	cmd.Flags().StringVarP(&deadlineDuration, "deadline-duration", "d", "30s", "Duration after which this transaction can no longer be executed. such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	return cmd
}

func getCmdCreateExchange(cdc *codec.Codec) *cobra.Command {
	// flags
	var token string
	cmd := &cobra.Command{
		Use:   "create_exchange",
		Short: "create exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgCreateExchange(token, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "CreateExchange by token name")
	return cmd
}

func getCmdTokenOKTSwap(cdc *codec.Codec) *cobra.Command {
	// flags
	var soldTokenAmount string
	var minBoughtTokenAmount string
	var deadline string
	var recipient string
	cmd := &cobra.Command{
		Use:   "token-okt",
		Short: "swap between token and okt",
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

			msg := types.NewMsgTokenOKTSwap(soldTokenAmount, minBoughtTokenAmount,
				deadline, recip, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&soldTokenAmount, "sell-amount", "", "",
		"amount expected to sell")
	cmd.Flags().StringVarP(&minBoughtTokenAmount, "min-buy-amount", "", "",
		"minimum amount expected to buy.")
	cmd.Flags().StringVarP(&recipient, "recipient", "", "",
		"address to receive the amount bought")
	cmd.Flags().StringVarP(&deadline, "deadline", "", "0s",
		"time after which this transaction can no longer be executed.")
	return cmd
}
