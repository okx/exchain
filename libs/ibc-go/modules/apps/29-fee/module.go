package fee

import (
	"context"
	"encoding/json"

	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/common"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/client/cli"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/keeper"

	cliCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	anytypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModuleAdapter      = AppModule{}
	_ module.AppModuleBasicAdapter = AppModuleBasic{}
	_ upgrade.UpgradeModule        = AppModule{}
)

// AppModuleBasic is the 29-fee AppModuleBasic
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterLegacyAminoCodec(codec)
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
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(ctx))
}

func (b AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return cli.NewTxCmd(cdc, reg)
}

func (b AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg anytypes.InterfaceRegistry) *cobra.Command {
	return cli.GetQueryCmd(cdc, reg)
}

func (b AppModuleBasic) RegisterRouterForGRPC(cliCtx cliCtx.CLIContext, r *mux.Router) {}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

type AppModule struct {
	*common.Veneus3BaseUpgradeModule
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	ret := AppModule{
		keeper: k,
	}
	ret.Veneus3BaseUpgradeModule = common.NewVeneus3BaseUpgradeModule(ret)
	return ret
}

func (a AppModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (a AppModule) Route() string {
	return types.RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
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

func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), a.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), a.keeper)
}

func (a AppModule) RegisterTask() upgrade.HeightTask {
	return upgrade.NewHeightTask(5, func(ctx sdk.Context) error {
		data := ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
		a.initGenesis(ctx, data)
		return nil
	})
}

func (am AppModule) initGenesis(ctx sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(message, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}
