package etherhash

import "golang.org/x/crypto/sha3"

// Sum returns the non-standard Keccak256 of the bz.
func Sum(bz []byte) []byte {
	// NewLegacyKeccak256 uses non-standard padding
	// and is incompatible with sha3.Sum256
	sha := sha3.NewLegacyKeccak256()
	sha.Write(bz)
	return sha.Sum(nil)
}
