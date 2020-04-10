package token

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/common/version"
	tokenTypes "github.com/okex/okchain/x/token/types"
)

var (
	_ module.AppModule = AppModule{}
)

// AppModule app module
type AppModule struct {
	AppModuleBasic
	keeper       Keeper
	supplyKeeper authTypes.SupplyKeeper
	version      version.ProtocolVersionType
}

// NewAppModule creates a new AppModule object
func NewAppModule(v version.ProtocolVersionType, keeper Keeper, supplyKeeper authTypes.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
		version:        v,
	}
}

// Name module name
func (AppModule) Name() string {
	return tokenTypes.ModuleName
}

// RegisterInvariants register invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	//RegisterInvariants(ir, am.keeper)
}

// Route module message route name
func (AppModule) Route() string {
	return tokenTypes.RouterKey
}

// NewHandler module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewTokenHandler(am.keeper, am.version)
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return tokenTypes.QuerierRoute
}

// NewQuerierHandler module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	tokenTypes.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return tokenTypes.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
