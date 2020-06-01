package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/margin/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group margin queries under a subcommand
	marginQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marginQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdAccountDeposit(queryRoute, cdc),
			GetCmdMarginProducts(queryRoute, cdc),
			GetCmdSaving(queryRoute, cdc),
			GetCmdBorrowing(queryRoute, cdc),
			// TODO: Add query Cmds
		)...,
	)

	return marginQueryCmd
}

func GetCmdAccountDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Query the margin account deposits",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryAccount, args[0]), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// TODO: Add Query Commands
func GetCmdMarginProducts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Query the margin products",
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProducts), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

func GetCmdSaving(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saving",
		Short: "Query the margin saving",
		RunE: func(cmd *cobra.Command, _ []string) error {
			product, err := cmd.Flags().GetString("product")
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QuerySaving, product), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().String("product", "",
		" Trading pair in full name of the tokens: ${baseAssetSymbol}_${quoteAssetSymbol}, for example "+
			"\"mycoin_"+common.NativeToken+"\"")
	return cmd
}

func GetCmdBorrowing(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrowing",
		Short: "Query the margin borrowing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}

			product, err := cmd.Flags().GetString("product")
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryBorrowing, product, args[0]), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().String("product", "",
		" Trading pair in full name of the tokens: ${baseAssetSymbol}_${quoteAssetSymbol}, for example "+
			"\"mycoin_"+common.NativeToken+"\"")
	return cmd
}
