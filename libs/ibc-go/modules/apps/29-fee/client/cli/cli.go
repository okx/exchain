package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "ibc-fee",
		Short:                      "IBC relayer incentivization query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdIncentivizedPacket(cdc, reg),
		GetCmdIncentivizedPackets(cdc, reg),
		GetCmdTotalRecvFees(cdc, reg),
		GetCmdTotalAckFees(cdc, reg),
		GetCmdTotalTimeoutFees(cdc, reg),
		GetCmdIncentivizedPacketsForChannel(cdc, reg),
		GetCmdPayee(cdc, reg),
		GetCmdCounterpartyPayee(cdc, reg),
		GetCmdFeeEnabledChannel(cdc, reg),
		GetCmdFeeEnabledChannels(cdc, reg),
	)

	return queryCmd
}

// NewTxCmd returns the transaction commands for 29-fee
func NewTxCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "ibc-fee",
		Short:                      "IBC relayer incentivization transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterPayeeCmd(cdc, reg),
		NewRegisterCounterpartyPayeeCmd(cdc, reg),
		NewPayPacketFeeAsyncTxCmd(cdc, reg),
	)

	return txCmd
}
