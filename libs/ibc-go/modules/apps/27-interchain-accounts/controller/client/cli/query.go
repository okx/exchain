package cli

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/version"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/controller/types"
	"github.com/spf13/cobra"
)

// GetCmdParams returns the command handler for the controller submodule parameter querying.
func GetCmdParams(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current interchain-accounts controller submodule parameters",
		Long:    "Query the current interchain-accounts controller submodule parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query interchain-accounts controller params", version.ServerName),
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
