package keeper

import (
	"encoding/json"
	wasmvm "github.com/CosmWasm/wasmvm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	wasmCache wasmvm.Cache
)

func SetWasmCache(cache wasmvm.Cache) {
	wasmCache = cache
}

func GetWasmCacheInfo() wasmvm.Cache {
	return wasmCache
}

func getCallerInfoFunc(ctx sdk.Context, keeper Keeper) func(contractAddress, storeAddress string) ([]byte, uint64, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, error) {
	return func(contractAddress, storeAddress string) ([]byte, uint64, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, error) {
		gasBefore := ctx.GasMeter().GasConsumed()
		codeHash, store, querier, gasMeter, err := getCallerInfo(ctx, keeper, contractAddress, storeAddress)
		gasAfter := ctx.GasMeter().GasConsumed()
		return codeHash, keeper.gasRegister.ToWasmVMGas(gasAfter - gasBefore), store, querier, gasMeter, err
	}
}

func getCallerInfo(ctx sdk.Context, keeper Keeper, contractAddress, storeAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, error) {
	cAddr, err := sdk.WasmAddressFromBech32(contractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// 1. get wasm code from contractAddress
	_, codeInfo, prefixStore, err := keeper.contractInstance(ctx, cAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// 2. contractAddress == storeAddress and direct return
	if contractAddress == storeAddress {
		queryHandler := keeper.newQueryHandler(ctx, cAddr)
		return codeInfo.CodeHash, prefixStore, queryHandler, keeper.gasMeter(ctx), nil
	}
	// 3. get store from storeaddress
	sAddr, err := sdk.WasmAddressFromBech32(storeAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	_, _, prefixStore, err = keeper.contractInstance(ctx, sAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queryHandler := keeper.newQueryHandler(ctx, sAddr)
	return codeInfo.CodeHash, prefixStore, queryHandler, keeper.gasMeter(ctx), nil
}

func transferCoinsFunc(ctx sdk.Context, keeper Keeper) func(contractAddress, caller string, coinsData []byte) (uint64, error) {
	return func(contractAddress, caller string, coinsData []byte) (uint64, error) {
		var coins sdk.Coins
		err := json.Unmarshal(coinsData, &coins)
		if err != nil {
			return 0, err
		}
		gasBefore := ctx.GasMeter().GasConsumed()
		err = transferCoins(ctx, keeper, contractAddress, caller, coins)
		gasAfter := ctx.GasMeter().GasConsumed()
		return keeper.gasRegister.ToWasmVMGas(gasAfter - gasBefore), err
	}
}

func transferCoins(ctx sdk.Context, keeper Keeper, contractAddress, caller string, coins sdk.Coins) error {
	if !coins.IsZero() {
		contractAddr, err := sdk.WasmAddressFromBech32(contractAddress)
		if err != nil {
			return err
		}
		callerAddr, err := sdk.WasmAddressFromBech32(caller)
		if err != nil {
			return err
		}
		if err := keeper.bank.TransferCoins(ctx, callerAddr, contractAddr, coins); err != nil {
			return err
		}
	}
	return nil
}
