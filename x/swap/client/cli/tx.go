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
	var maxBaseTokens string
	var quoteTokens string
	var deadline int64
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
			maxBaseTokensDecCoin, err2 := sdk.ParseDecCoin(maxBaseTokens)
			if err2 != nil {
				return err2
			}
			quoteTokensDecCoin, err2 := sdk.ParseDecCoin(quoteTokens)
			if err2 != nil {
				return err2
			}

			msg := types.NewMsgAddLiquidity(minLiquidityDec, maxBaseTokensDecCoin, quoteTokensDecCoin, deadline, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&minLiquidity, "min_liquidity", "l", "", "Minimum number of sender will mint if total pool token supply is greater than 0")
	cmd.Flags().StringVarP(&maxBaseTokens, "max_base_tokens", "", "", "Maximum number of base tokens deposited. Deposits max amount if total pool token supply is 0. For example \"100xxb\"")
	cmd.Flags().StringVarP(&quoteTokens, "quote_tokens", "q", "", "The number of quote tokens. For example \"100okb\"")
	cmd.Flags().Int64VarP(&deadline, "deadline", "d", 0, "Time after which this transaction can no longer be executed.")
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

			msg := types.NewMsgCreateExchange(token,cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "CreateExchange by token name")
	return cmd
}