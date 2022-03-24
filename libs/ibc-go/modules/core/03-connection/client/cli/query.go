package cli

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	utils "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/client/utils"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/spf13/cobra"
)

// GetCmdQueryConnections defines the command to query all the connection ends
// that this chain mantains.
func GetCmdQueryConnections(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connections",
		Short:   "Query all connections",
		Long:    "Query all connections ends from a chain",
		Example: fmt.Sprintf("%s query %s %s connections", version.ServerName, host.ModuleName, types.SubModuleName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryConnectionsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Connections(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "connection ends")

	return cmd
}

// GetCmdQueryConnection defines the command to query a connection end
func GetCmdQueryConnection(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "end [connection-id]",
		Short:   "Query stored connection end",
		Long:    "Query stored connection end",
		Example: fmt.Sprintf("%s query %s %s end [connection-id]", version.ServerName, host.ModuleName, types.SubModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			connectionID := args[0]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)
			connRes, err := utils.QueryConnection(clientCtx, connectionID, prove)
			if err != nil {
				return err
			}

			clientCtx = clientCtx.WithHeight(int64(connRes.ProofHeight.RevisionHeight))
			return clientCtx.PrintProto(connRes)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryClientConnections defines the command to query a client connections
func GetCmdQueryClientConnections(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "path [client-id]",
		Short:   "Query stored client connection paths",
		Long:    "Query stored client connection paths",
		Example: fmt.Sprintf("%s query  %s %s path [client-id]", version.ServerName, host.ModuleName, types.SubModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			clientID := args[0]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)

			connPathsRes, err := utils.QueryClientConnections(clientCtx, clientID, prove)
			if err != nil {
				return err
			}

			clientCtx = clientCtx.WithHeight(int64(connPathsRes.ProofHeight.RevisionHeight))
			return clientCtx.PrintProto(connPathsRes)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
