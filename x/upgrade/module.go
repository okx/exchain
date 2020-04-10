package upgrade

import (
	"encoding/json"
	"github.com/okex/okchain/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/okex/okchain/x/upgrade/client/cli"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// check the implementation of the interface
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	// just 4 test
	types.RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis(nothing 2 do)
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error {
	return nil
}

// register rest routes(nothing 2 do)
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {
	return
}

// get the root tx command of this module(undone)
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	// just 4 test
	return cli.GetTxCmd(cdc)
	//return nil
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// app module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// creates a new AppModule object for upgrade module
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic{},
		keeper,
	}
}

// module name
func (AppModule) Name() string {
	return ModuleName
}

//nothing 2 do
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// nothing 2 do
func (AppModule) ExportGenesis(_ sdk.Context) json.RawMessage {
	return nil
}

// nothing 2 do
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
	return
}

// module message route name
func (AppModule) Route() string {
	return RouterKey
}

// nothing 2 do
func (am AppModule) NewHandler() sdk.Handler {
	// just 4 test
	return NewHandler(am.keeper)
	//return nil
}

// module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// nothing 2 do
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return nil
}
