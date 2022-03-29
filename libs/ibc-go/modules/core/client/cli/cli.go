package cli

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	connection "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection"
	channel "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(
		ibcclient.GetTxCmd(cdc, reg),
		channel.GetTxCmd(),
	)

	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(codec *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	// Group ibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "Querying commands for the IBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(
		ibcclient.GetQueryCmd(codec, reg),
		connection.GetQueryCmd(codec, reg),
		channel.GetQueryCmd(codec, reg),
		GetCmdParams(codec, reg),
	)

	return ibcQueryCmd
}

// GetCmdParams returns the command handler for ibc parameter querying.
func GetCmdParams(m *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current ibc parameters",
		Long:    "Query the current ibc parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query %s params", version.ServerName, host.ModuleName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := context.NewCLIContext().WithProxy(m).WithInterfaceRegistry(reg)
			queryClient := types.NewQueryClient(clientCtx)

			res, _ := queryClient.IbcParams(cmd.Context(), &types.QueryIbcParamsRequest{})
			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
