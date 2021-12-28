package etherhash

import "golang.org/x/crypto/sha3"

// Sum returns the Keccak256 of the bz.
func Sum(bz []byte) []byte {
	h := sha3.Sum256(bz)
	return h[:]
}
