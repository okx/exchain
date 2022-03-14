package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(
	//solomachine.GetTxCmd(),
	//tendermint.GetTxCmd(),
	//connection.GetTxCmd(),
	//channel.GetTxCmd(),
	)

	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group ibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "Querying commands for the IBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(
	//ibcclient.GetQueryCmd(),
	//connection.GetQueryCmd(),
	//channel.GetQueryCmd(),
	)

	return ibcQueryCmd
}
