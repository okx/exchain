package wasm

import (
	"encoding/json"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	cdctypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
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
	return nil
}

func (b AppModuleBasic) GetQueryCmdV2(cdc *codec.CodecProxy, reg cdctypes.InterfaceRegistry) *cobra.Command {
	return nil
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
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	validators, err := InitGenesis(ctx, am.keeper, genesisState, am.NewHandler())
	if err != nil {
		panic(err)
	}
	return validators
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
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
