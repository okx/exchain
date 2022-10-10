package feesplit

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/ibc-go/modules/core/base"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/feesplit/client/cli"
	"github.com/okex/exchain/x/feesplit/keeper"
	"github.com/okex/exchain/x/feesplit/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ upgrade.UpgradeModule = AppModule{}
)

// AppModuleBasic type for the fees module
type AppModuleBasic struct{}

// Name returns the fees module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers types for module
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis is json default structure
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	//return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
	return nil
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	if len(bz) > 0 {
		var genesisState types.GenesisState
		err := types.ModuleCdc.UnmarshalJSON(bz, &genesisState)
		if err != nil {
			return err
		}

		return genesisState.Validate()
	}
	return nil
}

// RegisterRESTRoutes Registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

// GetQueryCmd Gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.ModuleName, cdc)
}

// GetTxCmd returns the root tx command for the swap module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// ___________________________________________________________________________

// AppModule implements the AppModule interface for the fees module.
type AppModule struct {
	AppModuleBasic
	*base.BaseIBCUpgradeModule
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k keeper.Keeper) AppModule {
	m := AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
	m.BaseIBCUpgradeModule = base.NewBaseIBCUpgradeModule(m)
	return m
}

// Name returns the fees module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the fees module's invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// NewHandler returns nil - fees module doesn't expose tx gRPC endpoints
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// Route returns the fees module's message routing key.
func (am AppModule) Route() string {
	return types.RouterKey
}

// QuerierRoute returns the claim module's query routing key.
func (am AppModule) QuerierRoute() string {
	return types.RouterKey
}

// NewQuerierHandler sets up new querier handler for module
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// BeginBlock executes all ABCI BeginBlock logic respective to the fees module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
	if tmtypes.DownloadDelta {
		types.GetParamsCache().SetNeedParamsUpdate()
	}
}

// EndBlock executes all ABCI EndBlock logic respective to the fees module. It
// returns no validator updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis performs the fees module's genesis initialization. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//var genesisState types.GenesisState
	//
	//types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	//InitGenesis(ctx, am.keeper, genesisState)
	//return []abci.ValidatorUpdate{}
	return nil
}

// ExportGenesis returns the fees module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	if !tmtypes.HigherThanVenus3(ctx.BlockHeight()) {
		return nil
	}
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
	return nil
}
