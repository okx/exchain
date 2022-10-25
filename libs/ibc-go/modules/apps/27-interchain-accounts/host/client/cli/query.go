package cli

import (
	"fmt"
	"strconv"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/spf13/cobra"
)

// GetCmdParams returns the command handler for the host submodule parameter querying.
func GetCmdParams(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current interchain-accounts host submodule parameters",
		Long:    "Query the current interchain-accounts host submodule parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query interchain-accounts host params", version.ServerName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := context.NewCLIContext().WithProxy(cdc).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdPacketEvents returns the command handler for the host packet events querying.
func GetCmdPacketEvents(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "packet-events [channel-id] [sequence]",
		Short:   "Query the interchain-accounts host submodule packet events",
		Long:    "Query the interchain-accounts host submodule packet events for a particular channel and sequence",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query interchain-accounts host packet-events channel-0 100", version.ServerName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(cdc).WithInterfaceRegistry(reg)

			channelID, portID := args[0], icatypes.PortID
			if err := host.ChannelIdentifierValidator(channelID); err != nil {
				return err
			}

			seq, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			searchEvents := []string{
				fmt.Sprintf("%s.%s='%s'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeyDstChannel, channelID),
				fmt.Sprintf("%s.%s='%s'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeyDstPort, portID),
				fmt.Sprintf("%s.%s='%d'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeySequence, seq),
			}

			result, err := utils.Query40TxsByEvents(clientCtx, searchEvents, 1, 1)
			if err != nil {
				return err
			}

			var resEvents []sdk.Event
			for _, r := range result.Txs {
				for _, v := range r.Events {
					eve := sdk.Event{
						Type:       v.Type,
						Attributes: v.Attributes,
					}
					resEvents = append(resEvents, eve)
				}
			}

			return clientCtx.PrintOutput(sdk.StringifyEvents(resEvents).String())
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
