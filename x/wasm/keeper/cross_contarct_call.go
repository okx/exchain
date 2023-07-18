package keeper

import "C"

import (
	"encoding/json"
	"errors"
	wasmvm "github.com/CosmWasm/wasmvm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"unsafe"
)

var (
	wasmKeeper Keeper

	// wasmvm cache param
	filePath            string
	supportedFeatures   string
	contractMemoryLimit uint32 = ContractMemoryLimit
	contractDebugMode   bool
	memoryCacheSize     uint32

	wasmCache wasmvm.Cache
)

func SetWasmKeeper(k *Keeper) {
	wasmKeeper = *k
}

func SetWasmCache(cache wasmvm.Cache) {
	wasmCache = cache
}

func GetWasmCacheInfo() (wasmvm.GoAPI, wasmvm.Cache) {
	return cosmwasmAPI, wasmCache
}

func GetWasmCallInfo(q unsafe.Pointer, contractAddress, storeAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, error) {
	goQuerier := *(*wasmvm.Querier)(q)
	qq, ok := goQuerier.(QueryHandler)
	if !ok {
		return nil, nil, nil, nil, errors.New("can not switch the pointer to the QueryHandler")
	}
	return getCallerInfo(qq.Ctx, contractAddress, storeAddress)
}

func getCallerInfo(ctx sdk.Context, contractAddress, storeAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, error) {
	cAddr, err := sdk.WasmAddressFromBech32(contractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// 1. get wasm code from contractAddress
	_, codeInfo, prefixStore, err := wasmKeeper.contractInstance(ctx, cAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// 2. contractAddress == storeAddress and direct return
	if contractAddress == storeAddress {
		queryHandler := wasmKeeper.newQueryHandler(ctx, cAddr)
		return codeInfo.CodeHash, prefixStore, queryHandler, wasmKeeper.gasMeter(ctx), nil
	}
	// 3. get store from storeaddress
	sAddr, err := sdk.WasmAddressFromBech32(storeAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	_, _, prefixStore, err = wasmKeeper.contractInstance(ctx, sAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queryHandler := wasmKeeper.newQueryHandler(ctx, sAddr)
	return codeInfo.CodeHash, prefixStore, queryHandler, wasmKeeper.gasMeter(ctx), nil
}

func TransferCoins(q unsafe.Pointer, contractAddress, caller string, coinsData []byte) error {
	goQuerier := *(*wasmvm.Querier)(q)
	qq, ok := goQuerier.(QueryHandler)
	if !ok {
		return errors.New("can not switch the pointer to the QueryHandler")
	}
	var coins sdk.Coins
	err := json.Unmarshal(coinsData, &coins)
	if err != nil {
		return err
	}
	return transferCoins(qq.Ctx, contractAddress, caller, coins)
}
func transferCoins(ctx sdk.Context, contractAddress, caller string, coins sdk.Coins) error {
	if !coins.IsZero() {
		contractAddr, err := sdk.WasmAddressFromBech32(contractAddress)
		if err != nil {
			return err
		}
		callerAddr, err := sdk.WasmAddressFromBech32(caller)
		if err != nil {
			return err
		}
		if err := wasmKeeper.bank.TransferCoins(ctx, callerAddr, contractAddr, coins); err != nil {
			return err
		}
	}
	return nil
}
