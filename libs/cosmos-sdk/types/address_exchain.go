package types

import (
	"crypto/sha256"

	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/crypto"
)

// MustBech32ifyAccPub returns the result of Bech32ifyAccPub panicing on failure.
func MustBech32ifyAccPub(pub crypto.PubKey) string {
	return MustBech32ifyPubKey(Bech32PubKeyTypeAccPub, pub)
}

// Derive derives a new address from the main `address` and a derivation `key`.
func Derive(address []byte, key []byte) []byte {
	return Hash(UnsafeBytesToStr(address), key)
}

// Hash creates a new address from address type and key
func Hash(typ string, key []byte) []byte {
	hasher := sha256.New()
	_, err := hasher.Write(UnsafeStrToBytes(typ))
	// the error always nil, it's here only to satisfy the io.Writer interface
	errors.AssertNil(err)
	th := hasher.Sum(nil)

	hasher.Reset()
	_, err = hasher.Write(th)
	errors.AssertNil(err)
	_, err = hasher.Write(key)
	errors.AssertNil(err)
	return hasher.Sum(nil)
}
