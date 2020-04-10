package protocol

import (
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/tendermint/tendermint/crypto"
)

func newPubKey(pubKey string) (res crypto.PubKey) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	var pubKeyEd25519 ed25519.PubKeyEd25519
	copy(pubKeyEd25519[:], pubKeyBytes[:])
	return pubKeyEd25519
}
