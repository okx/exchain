package params

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// GenesisState contains all params state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// DefaultGenesisState returns the default genesis state of this module
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	if !data.Params.MinDeposit.IsValid() {
		return fmt.Errorf("params deposit amount must be a valid sdk.Coins amount, is %s",
			data.Params.MinDeposit.String())
	}
	return nil
}

// AppModuleBasic is the struct of app module basics object
type AppModuleBasic struct{}

// Name returns the module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns the default genesis state in json raw message
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis gives a quick validity check for module genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// nolint
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command                 { return nil }
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command              { return nil }

// AppModule is the struct of this app module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Route returns the module route name
func (AppModule) Route() string {
	return RouterKey
}

// InitGenesis initializes the module genesis state
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.SetParams(ctx, genesisState.Params)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the module genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := GenesisState{
		Params: am.keeper.GetParams(ctx),
	}
	return ModuleCdc.MustMarshalJSON(gs)
}

// nolint
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry)        {}
func (AppModule) NewHandler() sdk.Handler                            { return nil }
func (AppModule) QuerierRoute() string                               { return "" }
func (AppModule) NewQuerierHandler() sdk.Querier                     { return nil }
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}
func (AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
