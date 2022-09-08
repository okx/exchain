package ica

import (
	"encoding/json"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
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

// TODO,这里需要进行整个适配

// AppModuleBasic is the IBC interchain accounts AppModuleBasic
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) DefaultGenesis() json.RawMessage {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) ValidateGenesis(message json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterRESTRoutes(context context.CLIContext, router *mux.Router) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetTxCmd(codec *codec.Codec) *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetQueryCmd(codec *codec.Codec) *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterInterfaces(registry anytypes.InterfaceRegistry) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(context context.CLIContext, mux *runtime.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterRouterForGRPC(cliCtx context.CLIContext, r *mux.Router) {
	//TODO implement me
	panic("implement me")
}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// AppModule is the application module for the IBC interchain accounts module
type AppModule struct {
	AppModuleBasic
	controllerKeeper *controllerkeeper.Keeper
	hostKeeper       *hostkeeper.Keeper
}

func (a AppModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) ExportGenesis(s sdk.Context) json.RawMessage {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) Route() string {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) NewHandler() sdk.Handler {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) QuerierRoute() string {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) EndBlock(s sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterServices(configurator module.Configurator) {
	//TODO implement me
	panic("implement me")
}
