package capability

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/capability/simulation"
	"github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	simulation2 "github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/core/base"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
	"math/rand"
)

var (
	_ module.AppModuleAdapter      = AppModule{}
	_ module.AppModuleBasicAdapter = AppModuleBasic{}
	_ module.AppModuleSimulation   = AppModule{}
	_ upgrade.UpgradeModule        = AppModule{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc *codec.CodecProxy
}

func NewAppModuleBasic(cdc *codec.CodecProxy) AppModuleBasic {
	ret := AppModuleBasic{cdc: cdc}
	return ret
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// Name returns the capability module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec does nothing. Capability does not support amino.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.Codec) {}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(_ types2.InterfaceRegistry) {}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the capability module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the capability module's REST service handlers.
func (a AppModuleBasic) RegisterRESTRoutes(ctx clientCtx.CLIContext, rtr *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the capability module.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx clientCtx.CLIContext, mux *runtime.ServeMux) {
}

// GetTxCmd returns the capability module's root tx command.
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command { return nil }

func (am AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg types2.InterfaceRegistry) *cobra.Command {
	return nil
}

// GetQueryCmd returns the capability module's root query command.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg types2.InterfaceRegistry) *cobra.Command {
	return nil
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the capability module.
type AppModule struct {
	AppModuleBasic
	*base.BaseIBCUpgradeModule
	keeper keeper.Keeper
}

func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func NewAppModule(cdc *codec.CodecProxy, keeper keeper.Keeper) AppModule {
	ret := AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
	ret.BaseIBCUpgradeModule = base.NewBaseIBCUpgradeModule(ret)
	return ret
}

// Name returns the capability module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// Route returns the capability module's message routing key.
func (AppModule) Route() string { return types.ModuleName }

// QuerierRoute returns the capability module's query routing key.
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler returns the capability module's Querier.
func (am AppModule) LegacyQuerierHandler(codec2 *codec.Codec) sdk.Querier { return nil }

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(module.Configurator) {}

// RegisterInvariants registers the capability module's invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the capability module's genesis initialization It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return nil
}
func (am AppModule) initGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	ModuleCdc.MustUnmarshalJSON(data, &genState)

	InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the capability module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

func (am AppModule) exportGenesis(ctx sdk.Context) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(genState)
}

// BeginBlock executes all ABCI BeginBlock logic respective to the capability module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	if tmtypes.HigherThanVenus1(ctx.BlockHeight()) {
		am.keeper.InitMemStore(ctx)
	}
}

// EndBlock executes all ABCI EndBlock logic respective to the capability module. It
// returns no validator updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState creates a randomized GenState of the capability module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents performs a no-op
func (am AppModule) ProposalContents(simState module.SimulationState) []simulation2.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized capability param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simulation2.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for capability module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore()
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simulation2.WeightedOperation {
	return nil
}

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return upgrade.NewHeightTask(
		0, func(ctx sdk.Context) error {
			if am.Sealed() {
				return nil
			}
			data := ModuleCdc.MustMarshalJSON(types.DefaultGenesis())
			am.initGenesis(ctx, data)
			return nil
		})
}
