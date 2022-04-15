package etherhash

import (
	"hash"
	"sync"

	"golang.org/x/crypto/sha3"
)

var keccakPool = sync.Pool{
	// NewLegacyKeccak256 uses non-standard padding
	// and is incompatible with sha3.Sum256
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// Sum returns the non-standard Keccak256 of the bz.
func Sum(bz []byte) []byte {
	sha := keccakPool.Get().(keccakState)
	defer func() {
		// better to reset before putting it to the pool
		sha.Reset()
		keccakPool.Put(sha)
	}()
	sha.Write(bz)

	var hashData [32]byte
	sha.Read(hashData[:])
	return hashData[:]
}
