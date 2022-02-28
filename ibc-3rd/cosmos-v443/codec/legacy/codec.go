package legacy

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	cryptocodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
	cryptotypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/types"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec.LegacyAmino

func init() {
	Cdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	codec.RegisterEvidences(Cdc)
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey cryptotypes.PrivKey, err error) {
	err = Cdc.Unmarshal(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey cryptotypes.PubKey, err error) {
	err = Cdc.Unmarshal(pubKeyBytes, &pubKey)
	return
}
