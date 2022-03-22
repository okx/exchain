package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC channels
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC channel query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryChannels(cdc, reg),
		GetCmdQueryChannel(cdc, reg),
		GetCmdQueryConnectionChannels(cdc, reg),
		GetCmdQueryChannelClientState(cdc, reg),
		GetCmdQueryPacketCommitment(cdc, reg),
		GetCmdQueryPacketCommitments(cdc, reg),
		GetCmdQueryPacketReceipt(cdc, reg),
		GetCmdQueryPacketAcknowledgement(cdc, reg),
		GetCmdQueryUnreceivedPackets(cdc, reg),
		GetCmdQueryUnreceivedAcks(cdc, reg),
		GetCmdQueryNextSequenceReceive(cdc, reg),
		//// TODO: next sequence Send ?
	)

	return queryCmd
}

// NewTxCmd returns a CLI command handler for all x/ibc channel transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC channel transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()

	return txCmd
}
