package staking

import (
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	anytypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"

	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	params2 "github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/params"
	"github.com/okex/exchain/x/staking/client/rest"
	"github.com/okex/exchain/x/staking/types"
	"github.com/spf13/cobra"
)

var (
	_ upgrade.UpgradeModule        = AppModule{}
	_ module.AppModuleAdapter      = AppModule{}
	_ module.AppModuleBasicAdapter = AppModuleBasic{}
)

// appmoduleBasic
func (am AppModuleBasic) RegisterRouterForGRPC(cliCtx context.CLIContext, r *mux.Router) {
	rest.RegisterOriginRPCRoutersForGRPC(cliCtx, r)
}

func (am AppModuleBasic) RegisterInterfaces(registry anytypes.InterfaceRegistry) {}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(cliContext context.CLIContext, serveMux *runtime.ServeMux) {
}

func (am AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return nil
}

func (am AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return nil
}

/// appmodule
func (am AppModule) RegisterServices(configurator module.Configurator) {}

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return nil
}

func (am AppModule) UpgradeHeight() int64 {
	return -1
}

func (am AppModule) RegisterParam() params.ParamSet {
	v := types.KeyHistoricalEntriesParams(7)
	return params2.ParamSet(v)
}

func (am AppModule) ModuleName() string {
	return ModuleName
}

func (am AppModule) CommitFilter() *cosmost.StoreFilter {
	return nil
}

func (am AppModule) PruneFilter() *cosmost.StoreFilter {
	return nil
}

func (am AppModule) VersionFilter() *cosmost.VersionFilter {
	return nil
}
