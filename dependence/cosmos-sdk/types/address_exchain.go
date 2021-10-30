package types

import "github.com/okex/exchain/dependence/tendermint/crypto"

// MustBech32ifyAccPub returns the result of Bech32ifyAccPub panicing on failure.
func MustBech32ifyAccPub(pub crypto.PubKey) string {
	return MustBech32ifyPubKey(Bech32PubKeyTypeAccPub, pub)
}