package cli

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/spf13/cobra"
)

// GetCmdQueryDenomTrace defines the command to query a a denomination trace from a given hash.
func GetCmdQueryDenomTrace(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-trace [hash]",
		Short:   "Query the denom trace info from a given trace hash",
		Long:    "Query the denom trace info from a given trace hash",
		Example: fmt.Sprintf("%s query ibc-transfer denom-trace [hash]", version.ServerName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryDenomTraceRequest{
				Hash: args[0],
			}

			res, err := queryClient.DenomTrace(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryDenomTraces defines the command to query all the denomination trace infos
// that this chain mantains.
func GetCmdQueryDenomTraces(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-traces",
		Short:   "Query the trace info for all token denominations",
		Long:    "Query the trace info for all token denominations",
		Example: fmt.Sprintf("%s query ibc-transfer denom-traces", version.ServerName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryDenomTracesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DenomTraces(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "denominations trace")

	return cmd
}

// GetCmdParams returns the command handler for ibc-transfer parameter querying.
func GetCmdParams(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current ibc-transfer parameters",
		Long:    "Query the current ibc-transfer parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query ibc-transfer params", version.ServerName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)

			res, _ := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdParams returns the command handler for ibc-transfer parameter querying.
func GetCmdQueryEscrowAddress(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "escrow-address",
		Short:   "Get the escrow address for a channel",
		Long:    "Get the escrow address for a channel",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query ibc-transfer escrow-address [port] [channel-id]", version.ServerName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			port := args[0]
			channel := args[1]
			addr := types.GetEscrowAddress(port, channel)
			return clientCtx.PrintOutput(fmt.Sprintf("%s\n", addr.String()))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
