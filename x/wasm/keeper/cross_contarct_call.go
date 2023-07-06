package keeper

import "C"

import (
	"fmt"
	wasmvm "github.com/CosmWasm/wasmvm"
	"github.com/CosmWasm/wasmvm/api"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/types"
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

func SetWasmCacheParam(dataDir string,
	supportedFeatures string,
	memoryLimit uint32,
	printDebug bool,
	cacheSize uint32) {
	filePath = dataDir
	supportedFeatures = supportedFeatures
	contractMemoryLimit = memoryLimit
	contractDebugMode = printDebug
	memoryCacheSize = cacheSize
}

func GetCacheInfo() (wasmvm.GoAPI, api.Cache) {
	return cosmwasmAPI, wasmCache
}

func GenerateCallerInfo(q unsafe.Pointer, contractAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter) {
	goQuerier := *(*wasmvm.Querier)(q)
	qq := goQuerier.(QueryHandler)
	code, _, store, querier, gasMeter := generateCallerInfo(qq.Ctx, contractAddress)
	return code, store, querier, gasMeter
}

func generateCallerInfo(ctx sdk.Context, addr string) ([]byte, wasmvmtypes.Env, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter) {
	contractAddress, err := sdk.WasmAddressFromBech32(addr)
	if err != nil {
		panic(fmt.Sprintln("WasmAddressFromBech32 err", err))
	}
	env := types.NewEnv(ctx, contractAddress)
	_, codeInfo, prefixStore, err := wasmKeeper.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, env, nil, nil, nil
	}
	queryHandler := wasmKeeper.newQueryHandler(ctx, contractAddress)
	return codeInfo.CodeHash, env, prefixStore, queryHandler, wasmKeeper.gasMeter(ctx)
}
