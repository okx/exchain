package global

import (
	"github.com/VictoriaMetrics/fastcache"
	"math/big"
	"unsafe"
)

const cacheBlockSize = 2000

var blockEvmTxGasCache = fastcache.New(cacheBlockSize * (8 + 8))

func GetBlockEvmTxGasUsed(height int64) *big.Int {
	k := (*[8]byte)(unsafe.Pointer(&height))
	gasBytes, ok := blockEvmTxGasCache.HasGet(nil, k[:])
	if !ok {
		return nil
	}
	gas := new(big.Int)
	gas.SetBytes(gasBytes)
	return gas
}

func SetBlockEvmTxGasUsed(height int64, gas *big.Int) {
	k := (*[8]byte)(unsafe.Pointer(&height))
	blockEvmTxGasCache.Set(k[:], gas.Bytes())
}
