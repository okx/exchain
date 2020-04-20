package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okchain/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "swap",
		Short: "swap module",
	}

	txCmd.AddCommand(client.PostCommands(
		getCmdAddLiquidity(cdc),

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

			minLiquidityDec, err :=sdk.NewDecFromStr(minLiquidity)
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
