package keeper

import (
	"encoding/json"
	"fmt"
	"testing"

	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/keeper/wasmtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkWasmCreate(b *testing.B) {
	ctx, keepers := CreateTestInput(b, false, SupportedFeatures)
	creator := keepers.Faucet.NewFundedAccount(ctx, sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))...)
	ctx.SetEventManager(sdk.NewEventManager())

	//reset timer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		contractID, err := keepers.ContractKeeper.Create(ctx, creator, hackatomWasm, nil)
		require.NoError(b, err)
		require.Equal(b, uint64(i+1), contractID)
	}
}

func BenchmarkWasmInstantiate(b *testing.B) {
	ctx, keepers := CreateTestInput(b, false, SupportedFeatures)
	wasmerMock := &wasmtesting.MockWasmer{
		InstantiateFn: func(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, initMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64, deserCost wasmvmtypes.UFraction) (*wasmvmtypes.Response, uint64, error) {
			return &wasmvmtypes.Response{Data: []byte("my-response-data")}, 0, nil
		},
		AnalyzeCodeFn: wasmtesting.WithoutIBCAnalyzeFn,
		CreateFn:      wasmtesting.NoOpCreateFn,
	}
	example := StoreRandomContract(b, ctx, keepers, wasmerMock)

	//reset timer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, data, err := keepers.ContractKeeper.Instantiate(ctx, example.CodeID, example.CreatorAddr, nil, nil, "test", nil)
		require.NoError(b, err)
		assert.Equal(b, []byte("my-response-data"), data)
	}
}

func BenchmarkWasmExecute(b *testing.B) {
	ctx, keepers := CreateTestInput(b, false, SupportedFeatures)
	_, keeper, _ := keepers.AccountKeeper, keepers.ContractKeeper, keepers.BankKeeper
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit.Add(deposit...)...)
	fred := keepers.Faucet.NewFundedAccount(ctx, topUp...)
	contractID, _ := keeper.Create(ctx, creator, hackatomWasm, nil)
	_, _, bob := keyPubAddr()
	initMsgBz, _ := json.Marshal(HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	})
	addr, _, _ := keepers.ContractKeeper.Instantiate(ctx, contractID, creator, nil, initMsgBz, "demo contract 3", deposit)
	ctx.SetEventManager(sdk.NewEventManager())

	//reset timer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := keepers.ContractKeeper.Execute(ctx, addr, fred, []byte(`{"release":{}}`), topUp)
		if i == 0 {
			require.NoError(b, err)
		} else {
			require.Error(b, err)
		}
	}
}

func BenchmarkWasmQuery(b *testing.B) {
	ctx, keepers := CreateTestInput(b, false, SupportedFeatures)
	example := InstantiateHackatomExampleContract(b, ctx, keepers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := keepers.WasmKeeper.QuerySmart(ctx, example.Contract, []byte(`{"verifier":{}}`))
		require.NoError(b, err)
		require.Equal(b, fmt.Sprintf("{\"verifier\":\"%s\"}", example.VerifierAddr.String()), string(result))
	}
}
