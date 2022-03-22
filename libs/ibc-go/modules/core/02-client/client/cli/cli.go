package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC clients
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryClientStates(cdc, reg),
		GetCmdQueryClientState(cdc, reg),
		GetCmdQueryConsensusStates(cdc, reg),
		GetCmdQueryConsensusState(cdc, reg),
		GetCmdQueryHeader(cdc, reg),
		GetCmdSelfConsensusState(cdc, reg),
		GetCmdParams(cdc, reg),
	)

	return queryCmd
}

// NewTxCmd returns the command to create and handle IBC clients
func NewTxCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC client transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewCreateClientCmd(cdc, reg),
		NewUpdateClientCmd(cdc, reg),
		NewSubmitMisbehaviourCmd(cdc, reg),
		NewUpgradeClientCmd(cdc, reg),
	)

	return txCmd
}
