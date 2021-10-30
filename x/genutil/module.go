package genutil

import (
	"encoding/json"

	"github.com/okex/exchain/x/genutil/types"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// AppModuleBasic is the struct of app module basics object
type AppModuleBasic struct{}

// Name returns the module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// DefaultGenesis returns the default genesis state in json raw message
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(GenesisState{})
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
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec)                         {}
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command                 { return nil }
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command              { return nil }

// AppModule is the struct of this app module
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	deliverTx     deliverTxfn
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper, deliverTx deliverTxfn) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		deliverTx:      deliverTx,
	})
}

// InitGenesis initializes the module genesis state
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, ModuleCdc, am.stakingKeeper, am.deliverTx, genesisState)
}

// ExportGenesis exports the module genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}
