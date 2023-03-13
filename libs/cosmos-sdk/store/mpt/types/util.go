package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
)

var (
	keccakStatePool = &sync.Pool{
		New: func() interface{} {
			return crypto.NewKeccakState()
		},
	}
)

func Keccak256HashWithSyncPool(data ...[]byte) (h ethcmn.Hash) {
	d := keccakStatePool.Get().(crypto.KeccakState)
	defer keccakStatePool.Put(d)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}
