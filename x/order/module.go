package order

import (
	"encoding/json"
	auth "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"

	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/spf13/cobra"

	"github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/order/client/cli"
	"github.com/okex/exchain/x/order/client/rest"
	"github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/order/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic : app module basics object
type AppModuleBasic struct{}

// Name : module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec : register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis : default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis : module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes : register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd : get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd : get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.QuerierRoute, cdc)
}

// AppModule : app module
type AppModule struct {
	AppModuleBasic
	keeper       keeper.Keeper
	supplyKeeper auth.SupplyKeeper
	version      version.ProtocolVersionType
}

// NewAppModule : creates a new AppModule object
func NewAppModule(v version.ProtocolVersionType, keeper keeper.Keeper, supplyKeeper auth.SupplyKeeper) AppModule {

	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
		version:        v,
	}
}

// RegisterInvariants : register invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.keeper)
}

// Route : module message route name
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler : module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewOrderHandler(am.keeper)
}

// QuerierRoute : module querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler : module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// InitGenesis : module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return nil
}

// ExportGenesis : module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock : module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock : module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return nil
}
