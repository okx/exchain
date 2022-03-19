package ibc

import (
	"encoding/json"
	"fmt"
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"

	//ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"

	//ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	simulation2 "github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	//ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/okex/exchain/libs/ibc-go/modules/core/base"
	"github.com/okex/exchain/libs/ibc-go/modules/core/client/cli"
	"github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/core/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModuleAdapter    = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the ibc module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the ibc module's name.
func (AppModuleBasic) Name() string {
	return host.ModuleName
}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	if !tmtypes.UpgradeIBCInRuntime() {
		return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
	}
	return nil
}

// ValidateGenesis performs genesis state validation for the mint module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	if tmtypes.UpgradeIBCInRuntime() {
		if nil == bz {
			return nil
		}
	}
	var data types.GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", "asd", err)
	}

	return data.Validate()
}

// RegisterRESTRoutes does nothing. IBC does not support legacy REST routes.
func (AppModuleBasic) RegisterRESTRoutes(ctx clientCtx.CLIContext, rtr *mux.Router) {}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// TODO
// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the ibc module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(ctx clientCtx.CLIContext, mux *runtime.ServeMux) {
	//clienttypes.RegisterQueryHandlerClient(context.Background(), mux, clienttypes.NewQueryClient(clientCtx))
	//connectiontypes.RegisterQueryHandlerClient(context.Background(), mux, connectiontypes.NewQueryClient(clientCtx))
	//channeltypes.RegisterQueryHandlerClient(context.Background(), mux, channeltypes.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the ibc module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the ibc module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module for the ibc module.
type AppModule struct {
	AppModuleBasic
	*base.BaseIBCUpgradeModule
	keeper *keeper.Keeper

	// create localhost by default
	createLocalhost bool
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper) AppModule {
	ret := AppModule{
		keeper: k,
	}
	ret.BaseIBCUpgradeModule = base.NewBaseIBCUpgradeModule(ret)
	return ret
}

func (am AppModule) Upgrade(req *abci.UpgradeReq) (*abci.ModuleUpgradeResp, error) {
	return nil, nil
}

// TODO
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(*a.keeper)
}

// Name returns the ibc module's name.
func (AppModule) Name() string {
	return host.ModuleName
}

// RegisterInvariants registers the ibc module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO:
}

// Route returns the message routing key for the ibc module.
func (am AppModule) Route() string {
	return host.RouterKey
	//return sdk.NewRoute(host.RouterKey, NewHandler(*am.keeper))
}

// QuerierRoute returns the ibc module's querier route name.
func (AppModule) QuerierRoute() string {
	return host.QuerierRoute
}

// LegacyQuerierHandler returns nil. IBC does not support the legacy querier.
//func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
//	return nil
//}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	clienttypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	connectiontypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	channeltypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryService(cfg.QueryServer(), am.keeper)

	//m := clientkeeper.NewMigrator(am.keeper.ClientKeeper)
	//cfg.RegisterMigration(host.ModuleName, 1, m.Migrate1to2)
}

// InitGenesis performs genesis initialization for the ibc module. It returns
// no validator updates.
//func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, bz json.RawMessage) []abci.ValidatorUpdate {
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	if !tmtypes.UpgradeIBCInRuntime() {
		return am.initGenesis(ctx, data)
	}
	return nil
}

func (am AppModule) initGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	err := ModuleCdc.UnmarshalJSON(data, &gs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", host.ModuleName, err))
	}
	InitGenesis(ctx, *am.keeper, am.createLocalhost, &gs)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the ibc
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	if !tmtypes.UpgradeIBCInRuntime() {
		return am.exportGenesis(ctx)
	}
	return nil
}

func (am AppModule) exportGenesis(ctx sdk.Context) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(ExportGenesis(ctx, *am.keeper))
}

func lazyGenesis() json.RawMessage {
	ret := DefaultGenesisState()
	return ModuleCdc.MustMarshalJSON(&ret)
}

// BeginBlock returns the begin blocker for the ibc module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if !tmtypes.HigherThanIBCHeight(req.Header.Height) {
		return
	}
	ibcclient.BeginBlocker(ctx, am.keeper.ClientKeeper)
}

// EndBlock returns the end blocker for the ibc module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the ibc module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simulation2.WeightedProposalContent {
	return nil
}

// RandomizedParams returns nil since IBC doesn't register parameter changes.
func (AppModule) RandomizedParams(_ *rand.Rand) []simulation2.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for ibc module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[host.StoreKey] = simulation.NewDecodeStore(*am.keeper)
}

// WeightedOperations returns the all the ibc module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simulation2.WeightedOperation {
	return nil
}

func (am AppModule) RegisterTask() upgrade.HeightTask {
	if !tmtypes.UpgradeIBCInRuntime() {
		return nil
	}
	return upgrade.NewHeightTask(4, func(ctx sdk.Context) error {
		data := lazyGenesis()
		logrusplugin.Info("core init genesis")
		am.initGenesis(ctx, data)
		return nil
	})
}

func (am AppModule) RegisterParam() params.ParamSet {
	return nil
}
