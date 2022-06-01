package keeper

import (
	"testing"

	authkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/wasm/keeper/wasmtesting"
	"github.com/okex/exchain/x/wasm/types"
)

func TestConstructorOptions(t *testing.T) {
	cfg := MakeEncodingConfig(t)
	specs := map[string]struct {
		srcOpt Option
		verify func(*testing.T, Keeper)
	}{
		"wasm engine": {
			srcOpt: WithWasmEngine(&wasmtesting.MockWasmer{}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockWasmer{}, k.wasmVM)
			},
		},
		"message handler": {
			srcOpt: WithMessageHandler(&wasmtesting.MockMessageHandler{}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockMessageHandler{}, k.messenger)
			},
		},
		"query plugins": {
			srcOpt: WithQueryHandler(&wasmtesting.MockQueryHandler{}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockQueryHandler{}, k.wasmVMQueryHandler)
			},
		},
		"message handler decorator": {
			srcOpt: WithMessageHandlerDecorator(func(old Messenger) Messenger {
				require.IsType(t, &MessageHandlerChain{}, old)
				return &wasmtesting.MockMessageHandler{}
			}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockMessageHandler{}, k.messenger)
			},
		},
		"query plugins decorator": {
			srcOpt: WithQueryHandlerDecorator(func(old WasmVMQueryHandler) WasmVMQueryHandler {
				require.IsType(t, QueryPlugins{}, old)
				return &wasmtesting.MockQueryHandler{}
			}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockQueryHandler{}, k.wasmVMQueryHandler)
			},
		},
		"coin transferrer": {
			srcOpt: WithCoinTransferrer(&wasmtesting.MockCoinTransferrer{}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockCoinTransferrer{}, k.bank)
			},
		},
		"costs": {
			srcOpt: WithGasRegister(&wasmtesting.MockGasRegister{}),
			verify: func(t *testing.T, k Keeper) {
				assert.IsType(t, &wasmtesting.MockGasRegister{}, k.gasRegister)
			},
		},
		"api costs": {
			srcOpt: WithAPICosts(1, 2),
			verify: func(t *testing.T, k Keeper) {
				t.Cleanup(setApiDefaults)
				assert.Equal(t, uint64(1), costHumanize)
				assert.Equal(t, uint64(2), costCanonical)
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			k := NewKeeper(&cfg.Marshaler, nil, params.NewSubspace(nil, nil, nil, ""), authkeeper.AccountKeeper{}, nil, nil, nil, nil, nil, nil, nil, "tempDir", types.DefaultWasmConfig(), SupportedFeatures, spec.srcOpt)
			spec.verify(t, k)
		})
	}
}

func setApiDefaults() {
	costHumanize = DefaultGasCostHumanAddress * DefaultGasMultiplier
	costCanonical = DefaultGasCostCanonicalAddress * DefaultGasMultiplier
}
