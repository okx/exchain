package cli

import (
	"fmt"
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
			return cliCtx.PrintOutput(newToPrint("pool"))
		},
	}
}

// TODO: remove it later
type toPrint struct {
	string
}

func (tp toPrint) String() string {
	return tp.string
}

func newToPrint(s string) toPrint {
	return toPrint{s}
}
