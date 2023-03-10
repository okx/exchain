package wasm

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	store "github.com/okx/okbchain/libs/cosmos-sdk/store/types"

	"github.com/gorilla/mux"
	clictx "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	cdctypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/upgrade"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	tmcli "github.com/okx/okbchain/libs/tendermint/libs/cli"
	types2 "github.com/okx/okbchain/libs/tendermint/types"
	"github.com/okx/okbchain/x/wasm/client/cli"
	"github.com/okx/okbchain/x/wasm/keeper"
	"github.com/okx/okbchain/x/wasm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const SupportedFeatures = keeper.SupportedFeatures

func (b AppModuleBasic) RegisterCodec(amino *codec.Codec) {
	RegisterCodec(amino)
}

func (b AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (b AppModuleBasic) ValidateGenesis(message json.RawMessage) error {
	return nil
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
	if !types2.HigherThanEarth(ctx.BlockHeight()) {
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
	defaultVersionFilter store.VersionFilter = func(h int64) func(cb func(name string, version int64)) {
		if h < 0 {
			return func(cb func(name string, version int64)) {}
		}

		return func(cb func(name string, version int64)) {
			cb(ModuleName, types2.GetEarthHeight())
		}
	}
)

func (am AppModule) CommitFilter() *store.StoreFilter {
	var filter store.StoreFilter
	// return false:
	//    a. module name mismatch, no processing required
	//    b. module names match and reach the upgrade height
	// return true:
	//    a. the upgrade height is 0, the module is disabled
	//    b. not reach the upgrade height
	filter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}

		if am.UpgradeHeight() == 0 {
			return true
		}

		if h == types2.GetEarthHeight() {
			if s != nil {
				s.SetUpgradeVersion(h)
			}
			return false
		}

		if types2.HigherThanEarth(h) {
			return false
		}

		return true
	}

	return &filter
}

func (am AppModule) PruneFilter() *store.StoreFilter {
	var filter store.StoreFilter
	filter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}

		if am.UpgradeHeight() == 0 {
			return true
		}

		if types2.HigherThanEarth(h) {
			return false
		}

		return true
	}
	return &filter
}

func (am AppModule) VersionFilter() *store.VersionFilter {
	return &defaultVersionFilter
}

func (am AppModule) UpgradeHeight() int64 {
	return types2.GetEarthHeight()
}

var (
	once        sync.Once
	gWasmConfig types.WasmConfig
	gWasmDir    string
)

func WasmDir() string {
	once.Do(Init)
	return gWasmDir
}

func WasmConfig() types.WasmConfig {
	once.Do(Init)
	return gWasmConfig
}

func Init() {
	wasmConfig, err := ReadWasmConfig()
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	gWasmConfig = wasmConfig
	gWasmDir = filepath.Join(viper.GetString(tmcli.HomeFlag), "data")
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
