package ica

import (
	"context"
	"encoding/json"

	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/common"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/client/cli"

	controllertypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/controller/types"
	hosttypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	cliCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	anytypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/spf13/cobra"

	controllerkeeper "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/controller/keeper"
	hostkeeper "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/keeper"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

var (
	_ module.AppModuleAdapter      = AppModule{}
	_ module.AppModuleBasicAdapter = AppModuleBasic{}

	_ porttypes.IBCModule = host.IBCModule{}
)

// AppModuleBasic is the IBC interchain accounts AppModuleBasic
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (b AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (b AppModuleBasic) ValidateGenesis(message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context cliCtx.CLIContext, router *mux.Router) {}

func (b AppModuleBasic) GetTxCmd(codec *codec.Codec) *cobra.Command {
	return nil
}

func (b AppModuleBasic) GetQueryCmd(codec *codec.Codec) *cobra.Command {
	return nil
}

func (b AppModuleBasic) RegisterInterfaces(registry anytypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(ctx cliCtx.CLIContext, mux *runtime.ServeMux) {
	controllertypes.RegisterQueryHandlerClient(context.Background(), mux, controllertypes.NewQueryClient(ctx))
	hosttypes.RegisterQueryHandlerClient(context.Background(), mux, hosttypes.NewQueryClient(ctx))
}

func (b AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return nil
}

func (b AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return cli.GetQueryCmd(cdc, reg)
}

func (b AppModuleBasic) RegisterRouterForGRPC(cliCtx cliCtx.CLIContext, r *mux.Router) {}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// AppModule is the application module for the IBC interchain accounts module
type AppModule struct {
	*common.Veneus3BaseUpgradeModule
	AppModuleBasic
	controllerKeeper *controllerkeeper.Keeper
	hostKeeper       *hostkeeper.Keeper
}

// NewAppModule creates a new 20-transfer module
func NewAppModule(m *codec.CodecProxy, ck *controllerkeeper.Keeper, hk *hostkeeper.Keeper) AppModule {
	ret := AppModule{
		controllerKeeper: ck,
		hostKeeper:       hk,
	}
	ret.Veneus3BaseUpgradeModule = common.NewVeneus3BaseUpgradeModule(ret)
	return ret
}

func (am AppModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (a AppModule) ExportGenesis(s sdk.Context) json.RawMessage {
	return nil
}

func (a AppModule) Route() string {
	return types.RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.hostKeeper, a.controllerKeeper)
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (a AppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {}

func (a AppModule) EndBlock(s sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (a AppModule) RegisterServices(cfg module.Configurator) {
	controllertypes.RegisterQueryServer(cfg.QueryServer(), a.controllerKeeper)
	hosttypes.RegisterQueryServer(cfg.QueryServer(), a.hostKeeper)
}

func (a AppModule) RegisterTask() upgrade.HeightTask {
	return upgrade.NewHeightTask(6, func(ctx sdk.Context) error {
		ret := types.DefaultGenesis()
		data := ModuleCdc.MustMarshalJSON(ret)
		a.initGenesis(ctx, data)
		return nil
	})
}

func (am AppModule) initGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(message, &genesisState)

	if am.controllerKeeper != nil {
		controllerkeeper.InitGenesis(s, *am.controllerKeeper, genesisState.ControllerGenesisState)
	}

	if am.hostKeeper != nil {
		hostkeeper.InitGenesis(s, *am.hostKeeper, genesisState.HostGenesisState)
	}

	return nil
}
