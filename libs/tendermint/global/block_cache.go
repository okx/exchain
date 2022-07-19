package global

import (
	"github.com/VictoriaMetrics/fastcache"
	"math/big"
	"unsafe"
)

var blockEvmTxGasCache = fastcache.New(1024 * 1024 * 32)

func GetBlockEvmTxGasLimit(height int64) *big.Int {
	k := (*[8]byte)(unsafe.Pointer(&height))
	gasBytes, ok := blockEvmTxGasCache.HasGet(nil, k[:])
	if !ok {
		return nil
	}
	gas := new(big.Int)
	gas.SetBytes(gasBytes)
	return gas
}

func SetBlockEvmTxGasLimit(height int64, gas *big.Int) {
	k := (*[8]byte)(unsafe.Pointer(&height))
	blockEvmTxGasCache.Set(k[:], gas.Bytes())
}
