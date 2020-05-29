package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryMarginAccount, args[0]), nil)
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