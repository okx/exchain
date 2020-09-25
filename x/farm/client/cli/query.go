package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/farm/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group farm queries under a subcommand
	farmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryPool(queryRoute, cdc),
			GetCmdQueryPools(queryRoute, cdc),
			GetCmdQueryEarnings(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
		)...,
	)

	return farmQueryCmd
}

// GetCmdQueryPool gets the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool [pool-name]",
		Short: "query a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about the kind of coins as reward, the balance and the amount to farm in one block.

Example:
$ %s query farm pool pool-airtoken1-eth
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// TODO:
			return cliCtx.PrintOutput(newToPrint(args[0]))
		},
	}
}

// GetCmdQueryPools gets the pools query command.
func GetCmdQueryPools(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pools",
		Short: "query for all pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all pools.

Example:
$ %s query farm pools
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// TODO:
			return cliCtx.PrintOutput(newToPrint("all pools"))
		},
	}
}

// GetCmdQueryEarnings gets the earnings query command.
func GetCmdQueryEarnings(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "earnings [pool-name] [address]",
		Short: "query the current earnings",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the amount of locked coins and yield available.

Example:
$ %s query farm earnings pool-airtoken1-eth okexchain1hw4r48aww06ldrfeuq2v438ujnl6alsz0685a0
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			// TODO:
			return cliCtx.PrintOutput(newToPrint(args[0] + " : " + accAddr.String()))
		},
	}
}

// GetCmdQueryParams gets the pools query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "query the current farm parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as farm parameters.

Example:
$ %s query farm params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// TODO:
			return cliCtx.PrintOutput(newToPrint("params"))
		},
	}
}

// TODO: remove it later
type toPrint struct {
	Reminder string
}

func (tp toPrint) String() string {
	return tp.Reminder
}

func newToPrint(s string) toPrint {
	return toPrint{s}
}
