package cli

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/margin/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	marginTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		RunE:  client.ValidateCmd,
	}

	marginTxCmd.AddCommand(flags.PostCommands(
		GetCmdDexDeposit(cdc),
		GetCmdDexWithdraw(cdc),
		GetCmdDexSet(cdc),
		GetCmdDexSave(cdc),
		GetCmdDexReturn(cdc),
		GetCmdDeposit(cdc),
	)...)

	return marginTxCmd
}

// GetCmdDexDeposit is the CLI command for doing dex-deposit
func GetCmdDexDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dex-deposit [product] [amount]",
		Short: "dex deposits an amount of token for a product",
		Args:  cobra.ExactArgs(2), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()

			product := args[0]
			// Get amount of coins
			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDexDeposit(address, product, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDexWithdraw is the CLI command for doing dex-withdraw
func GetCmdDexWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dex-withdraw [product] [amount]",
		Short: "dex withdraws an amount of token from a product",
		Args:  cobra.ExactArgs(2), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()

			product := args[0]
			// Get amount of coins
			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDexWithdraw(address, product, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDexSet is the CLI command for doing dex-set
func GetCmdDexSet(cdc *codec.Codec) *cobra.Command {
	var maxLeverageStr, borrowRate, maintenanceMarginRatio string
	cmd := &cobra.Command{
		Use:   "dex-set [product]",
		Short: "dex sets params for a product",
		Args:  cobra.ExactArgs(1), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()
			product := args[0]

			maxLeverage, err := sdk.NewDecFromStr(maxLeverageStr)
			if err != nil {
				return err
			}

			if maxLeverage.IsNegative() {
				return errors.New("invalid max-leverage")
			}

			var borrowRateDec sdk.Dec
			if len(borrowRate) > 0 {
				borrowRateDec, err = sdk.NewDecFromStr(borrowRate)
				if err != nil {
					return fmt.Errorf("invalid borrow-rate:%s", err.Error())
				}
			}

			var mmrDec sdk.Dec
			if len(maintenanceMarginRatio) > 0 {
				mmrDec, err = sdk.NewDecFromStr(maintenanceMarginRatio)
				if err != nil {
					return fmt.Errorf("invalid maintenance-margin-ratio:%s", err.Error())
				}
			}

			msg := types.NewMsgDexSet(address, product, maxLeverage, borrowRateDec, mmrDec)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringVarP(&maxLeverageStr, "max-leverage", "ml", "", "max leverage of the product")
	cmd.Flags().StringVarP(&borrowRate, "borrow-rate", "br", "", "interest rate on borrowing")
	cmd.Flags().StringVarP(&maintenanceMarginRatio, "maintenance-margin-ratio", "mmr", "", "when the position Margin Ratio (MR) is lower than the Maintenance Margin Ratio (MMR) , liquidation will be triggered")
	return cmd
}

// GetCmdDexSave is the CLI command for doing dex-save
func GetCmdDexSave(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dex-save [product] [amount]",
		Short: "dex saves an amount of token for borrowing",
		Args:  cobra.ExactArgs(2), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()

			product := args[0]
			// Get amount of coins
			amount, err := sdk.ParseDecCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDexSave(address, product, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDexReturn is the CLI command for doing dex-save
func GetCmdDexReturn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dex-return [product] [amount]",
		Short: "dex returns an amount of token for borrowing",
		Args:  cobra.ExactArgs(2), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get depositor address
			address := cliCtx.GetFromAddress()

			product := args[0]
			// Get amount of coins
			amount, err := sdk.ParseDecCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDexReturn(address, product, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDeposit is the CLI command for doing Deposit
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [product] [amount]",
		Short: "add deposit for margin trade product ",
		Args:  cobra.ExactArgs(2), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			from := cliCtx.GetFromAddress()
			product := args[0]
			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgDeposit(from, product, sdk.NewCoins(amount))
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDeposit is the CLI command for doing Deposit
func GetCmdBorrow(cdc *codec.Codec) *cobra.Command {
	var leverageStr string
	var depositStr string
	cmd := &cobra.Command{
		Use:   "borrow [product] ",
		Short: "add deposit for margin trade product ",
		Args:  cobra.ExactArgs(3), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			from := cliCtx.GetFromAddress()
			product := args[0]
			deposit, err := sdk.ParseDecCoin(depositStr)
			if err != nil {
				return err
			}
			leverageDec, err := sdk.NewDecFromStr(leverageStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgBorrow(from, product, deposit, leverageDec)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringVarP(&depositStr, "leverage", "l", "", "The leverage of the borrow")
	cmd.Flags().StringVarP(&leverageStr, "deposit", "d", "", "The deposit for  borrow token")
	return cmd
}
