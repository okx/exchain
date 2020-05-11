package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/saving/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	savingTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		RunE:  client.ValidateCmd,
	}

	savingTxCmd.AddCommand(flags.PostCommands(
		GetCmdDeposit(cdc),
		GetCmdWithdraw(cdc),
	)...)

	return savingTxCmd
}

// GetCmdDeposit is the CLI command for doing deposit
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [amount]",
		Short: "deposit an amount of token to saving module",
		Args:  cobra.ExactArgs(1), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}
			msg := types.MsgDeposit{
				Address: address,
				Amount:  amount,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdWithdraw is the CLI command for doing withdraw
func GetCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [amount]",
		Short: "withdraw an amount of token from saving module",
		Args:  cobra.ExactArgs(1), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// get withdraw address
			address := cliCtx.GetFromAddress()

			// get amount of coins
			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}
			msg := types.MsgWithdraw{
				Address: address,
				Amount:  amount,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
