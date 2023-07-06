package keeper

import "C"

import (
	"fmt"
	wasmvm "github.com/CosmWasm/wasmvm"
	"github.com/CosmWasm/wasmvm/api"
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

	wasmCache api.Cache
)

func SetWasmKeeper(k *Keeper) {
	wasmKeeper = *k
}

func SetWasmCache(cache api.Cache) {
	wasmCache = cache
}

func GetCacheInfo() (wasmvm.GoAPI, api.Cache) {
	return cosmwasmAPI, wasmCache
}

func GenerateCallerInfo(q unsafe.Pointer, contractAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter) {
	goQuerier := *(*wasmvm.Querier)(q)
	qq := goQuerier.(QueryHandler)
	code, store, querier, gasMeter := generateCallerInfo(qq.Ctx, contractAddress)
	return code, store, querier, gasMeter
}

func generateCallerInfo(ctx sdk.Context, addr string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter) {
	contractAddress, err := sdk.WasmAddressFromBech32(addr)
	if err != nil {
		panic(fmt.Sprintln("WasmAddressFromBech32 err", err))
	}
	_, codeInfo, prefixStore, err := wasmKeeper.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, nil, nil, nil
	}
	queryHandler := wasmKeeper.newQueryHandler(ctx, contractAddress)
	return codeInfo.CodeHash, prefixStore, queryHandler, wasmKeeper.gasMeter(ctx)
}
