package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/x/icamauth/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd creates and returns the icamauth query command
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the icamauth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(getInterchainAccountCmd(cdc, reg))

	return cmd
}

func getInterchainAccountCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "interchainaccounts [connection-id] [owner-account]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := context.NewCLIContext().WithProxy(cdc).WithInterfaceRegistry(reg)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.InterchainAccount(cmd.Context(), types.NewQueryInterchainAccountRequest(args[0], args[1]))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
