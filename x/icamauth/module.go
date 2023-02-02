package icamauth

import (
	"context"
	"encoding/json"

	"github.com/okex/exchain/x/icamauth/keeper"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/okex/exchain/x/icamauth/client/cli"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clictx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/x/icamauth/types"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModuleAdapter      = AppModule{}
	_ module.AppModuleBasicAdapter = AppModuleBasic{}
)

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc *codec.CodecProxy
}

func NewAppModuleBasic(cdc *codec.CodecProxy) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterCodec(c *codec.Codec) {
	types.RegisterCodec(c)
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (a AppModuleBasic) ValidateGenesis(message json.RawMessage) error {
	return nil
}

func (a AppModuleBasic) RegisterRESTRoutes(context clictx.CLIContext, router *mux.Router) {}

func (a AppModuleBasic) GetTxCmd(c *codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd(c *codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) RegisterInterfaces(registry interfacetypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx clictx.CLIContext, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(ctx))
	if err != nil {
		panic(err)
	}
}

func (a AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return cli.GetTxCmd(cdc, reg)
}

func (a AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return cli.GetQueryCmd(cdc, reg)
}

func (a AppModuleBasic) RegisterRouterForGRPC(cliCtx clictx.CLIContext, r *mux.Router) {}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(cdc *codec.CodecProxy, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

func (a AppModule) Route() string {
	return types.RouterKey
}

func (a AppModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (a AppModule) ExportGenesis(s sdk.Context) json.RawMessage {
	return nil
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(keeper.NewMsgServerImpl(a.keeper))
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (a AppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {}

func (a AppModule) EndBlock(s sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), a.keeper)
}
