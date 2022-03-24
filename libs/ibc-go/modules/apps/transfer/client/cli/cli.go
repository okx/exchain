package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC connections
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "ibc-transfer",
		Short:                      "IBC fungible token transfer query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenomTrace(cdc, reg),
		GetCmdQueryDenomTraces(cdc, reg),
		GetCmdParams(cdc, reg),
		GetCmdQueryEscrowAddress(cdc, reg),
	)

	return queryCmd
}

// NewTxCmd returns the transaction commands for IBC fungible token transfer
func NewTxCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "ibc-transfer",
		Short:                      "IBC fungible token transfer transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewTransferTxCmd(cdc, reg),
	)

	return txCmd
}
