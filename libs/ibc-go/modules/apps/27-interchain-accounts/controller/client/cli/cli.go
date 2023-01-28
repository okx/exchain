package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for the ICA controller submodule
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "controller",
		Short:                      "interchain-accounts controller subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdParams(cdc, reg),
	)

	return queryCmd
}
