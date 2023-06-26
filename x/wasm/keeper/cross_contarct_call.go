package keeper

import "C"

import (
	"fmt"
	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/types"
	"unsafe"
)

var (
	wasmKeeper Keeper
)

func SetWasmKeeper(k *Keeper) {
	wasmKeeper = *k
}

func GenerateCallerInfo(q unsafe.Pointer, contractAddress string) ([]byte, wasmvm.KVStore, wasmvm.Querier) {
	goQuerier := *(*wasmvm.Querier)(q)
	qq := goQuerier.(QueryHandler)
	code, _, store, querier := generateCallerInfo(qq.Ctx, contractAddress)
	return code, store, querier
}

func generateCallerInfo(ctx sdk.Context, addr string) ([]byte, wasmvmtypes.Env, wasmvm.KVStore, wasmvm.Querier) {
	contractAddress, err := sdk.WasmAddressFromBech32(addr)
	if err != nil {
		panic(fmt.Sprintln("WasmAddressFromBech32 err", err))
	}
	env := types.NewEnv(ctx, contractAddress)
	_, codeInfo, prefixStore, err := wasmKeeper.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, env, nil, nil
	}
	queryHandler := wasmKeeper.newQueryHandler(ctx, contractAddress)
	return codeInfo.CodeHash, env, prefixStore, &queryHandler
}
