package wasm

import (
	"encoding/json"
	store "github.com/okex/exchain/libs/cosmos-sdk/store/types"

	"github.com/gorilla/mux"
	clictx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	cdctypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/wasm/client/cli"
	"github.com/okex/exchain/x/wasm/keeper"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (b AppModuleBasic) RegisterCodec(amino *codec.Codec) {
	RegisterCodec(amino)
}

func (b AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(&GenesisState{
		Params: DefaultParams(),
	})
}

func (b AppModuleBasic) ValidateGenesis(message json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(message, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (b AppModuleBasic) GetTxCmdV2(cdc *codec.CodecProxy, reg cdctypes.InterfaceRegistry) *cobra.Command {
	return cli.NewTxCmd(cdc, reg)
}

func (b AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg cdctypes.InterfaceRegistry) *cobra.Command {
	return cli.NewQueryCmd(cdc, reg)
}

func (b AppModuleBasic) RegisterRouterForGRPC(cliCtx clictx.CLIContext, r *mux.Router) {

}
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(keeper.NewDefaultPermissionKeeper(am.keeper))
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewLegacyQuerier(am.keeper, am.keeper.QueryGasLimit())
}

// InitGenesis performs genesis initialization for the wasm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// Note: use RegisterTask instead

	//var genesisState GenesisState
	//ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	//validators, err := InitGenesis(ctx, am.keeper, genesisState, am.NewHandler())
	//if err != nil {
	//	panic(err)
	//}
	//return validators
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	if !types2.HigherThanVenus2(ctx.BlockHeight()) {
		return nil
	}
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return upgrade.NewHeightTask(
		0, func(ctx sdk.Context) error {
			if am.Sealed() {
				return nil
			}
			_, err := InitGenesis(ctx, am.keeper, GenesisState{Params: DefaultParams()}, am.NewHandler())
			return err
		})
}

var (
	defaultDenyFilter store.StoreFilter = func(module string, h int64, s store.CommitKVStore) bool {
		return module == ModuleName
	}

	defaultCommitFilter store.StoreFilter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}

		if h == types2.GetVenus2Height() {
			if s != nil {
				s.SetUpgradeVersion(h)
			}
			return false
		}

		if types2.HigherThanVenus2(h) {
			return false
		}

		return true
	}
	defaultPruneFilter store.StoreFilter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}

		if types2.HigherThanVenus2(h) {
			return false
		}

		return true
	}
	defaultVersionFilter store.VersionFilter = func(h int64) func(cb func(name string, version int64)) {
		if h < 0 {
			return func(cb func(name string, version int64)) {}
		}

		return func(cb func(name string, version int64)) {
			cb(ModuleName, types2.GetVenus2Height())
		}
	}
)

func (am AppModule) CommitFilter() *store.StoreFilter {
	if am.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultCommitFilter
}

func (am AppModule) PruneFilter() *store.StoreFilter {
	if am.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultPruneFilter
}

func (am AppModule) VersionFilter() *store.VersionFilter {
	return &defaultVersionFilter
}

func (am AppModule) UpgradeHeight() int64 {
	return types2.GetVenus2Height()
}

// ReadWasmConfig reads the wasm specifig configuration
func ReadWasmConfig() (types.WasmConfig, error) {
	cfg := types.DefaultWasmConfig()
	var err error
	if v := viper.Get(flagWasmMemoryCacheSize); v != nil {
		if cfg.MemoryCacheSize, err = cast.ToUint32E(v); err != nil {
			return cfg, err
		}
	}
	if v := viper.Get(flagWasmQueryGasLimit); v != nil {
		if cfg.SmartQueryGasLimit, err = cast.ToUint64E(v); err != nil {
			return cfg, err
		}
	}
	if v := viper.Get(flagWasmSimulationGasLimit); v != nil {
		if raw, ok := v.(string); ok && raw != "" {
			limit, err := cast.ToUint64E(v) // non empty string set
			if err != nil {
				return cfg, err
			}
			cfg.SimulationGasLimit = &limit
		}
	}
	// attach contract debugging to global "trace" flag
	if v := viper.Get(server.FlagTrace); v != nil {
		if cfg.ContractDebugMode, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}
