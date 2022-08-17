package gov

import (
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/okex/exchain/x/gov/types"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	anytypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	GovCli "github.com/okex/exchain/x/gov/client/cli"
)

var (
	_ module.AppModuleBasicAdapter = AppModuleBasic{}
)

func (a AppModuleBasic) RegisterInterfaces(registry anytypes.InterfaceRegistry) {
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(cliContext context.CLIContext, serveMux *runtime.ServeMux) {
}

func (a AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	proposalCLIHandlers := make([]*cobra.Command, len(a.proposalHandlers))
	for i, proposalHandler := range a.proposalHandlers {
		proposalCLIHandlers[i] = proposalHandler.CLIHandler(cdc, reg)
	}

	return GovCli.GetTxCmd(types.StoreKey, cdc.GetCdc(), proposalCLIHandlers)
}

func (a AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return nil
}

func (a AppModuleBasic) RegisterRouterForGRPC(cliCtx context.CLIContext, r *mux.Router) {}
