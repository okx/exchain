package cryptoamino

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/libs/tendermint/crypto/sr25519"
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

// nameTable is used to map public key concrete types back
// to their registered amino names. This should eventually be handled
// by amino. Example usage:
// nameTable[reflect.TypeOf(ed25519.PubKeyEd25519{})] = ed25519.PubKeyAminoName
var nameTable = make(map[reflect.Type]string, 3)

func init() {
	// NOTE: It's important that there be no conflicts here,
	// as that would change the canonical representations,
	// and therefore change the address.
	// TODO: Remove above note when
	// https://github.com/tendermint/go-amino/issues/9
	// is resolved
	RegisterAmino(cdc)

	// TODO: Have amino provide a way to go from concrete struct to route directly.
	// Its currently a private API
	nameTable[reflect.TypeOf(ed25519.PubKeyEd25519{})] = ed25519.PubKeyAminoName
	nameTable[reflect.TypeOf(sr25519.PubKeySr25519{})] = sr25519.PubKeyAminoName
	nameTable[reflect.TypeOf(secp256k1.PubKeySecp256k1{})] = secp256k1.PubKeyAminoName
	nameTable[reflect.TypeOf(multisig.PubKeyMultisigThreshold{})] = multisig.PubKeyMultisigThresholdAminoRoute
}

// PubkeyAminoName returns the amino route of a pubkey
// cdc is currently passed in, as eventually this will not be using
// a package level codec.
func PubkeyAminoName(cdc *amino.Codec, key crypto.PubKey) (string, bool) {
	route, found := nameTable[reflect.TypeOf(key)]
	return route, found
}

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PubKeySr25519{},
		sr25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{},
		multisig.PubKeyMultisigThresholdAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{},
		ed25519.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PrivKeySr25519{},
		sr25519.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
		secp256k1.PrivKeyAminoName, nil)
}

// RegisterKeyType registers an external key type to allow decoding it from bytes
func RegisterKeyType(o interface{}, name string) {
	cdc.RegisterConcrete(o, name, nil)
	nameTable[reflect.TypeOf(o)] = name
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}

// hard code here for performance
var typePubKeySecp256k1Prefix = []byte{0xeb, 0x5a, 0xe9, 0x87}
var typePubKeyEd25519Prefix = []byte{0x16, 0x24, 0xde, 0x64}
var typePubKeySr25519Prefix = []byte{0x0d, 0xfb, 0x10, 0x05}

// UnmarshalPubKeyFromAminoWithTypePrefix decode pubkey from amino bytes,
// bytes should start with type prefix
func UnmarshalPubKeyFromAminoWithTypePrefix(data []byte) (crypto.PubKey, error) {
	const typePrefixAndSizeLen = 4 + 1

	if data[0] == 0x00 {
		return nil, errors.New("unmarshal pubkey with disamb do not implement")
	}
	if len(data) < typePrefixAndSizeLen {
		return nil, errors.New("pubkey raw data size error")
	}

	prefix := data[0:4]
	size := data[4]

	if size == 0 {
		return nil, nil
	}
	if len(data) == typePrefixAndSizeLen {
		return nil, errors.New("pubkey raw data size error")
	}
	if size&0x80 == 0x80 {
		return nil, errors.New("pubkey amino data size should use one byte")
	}

	data = data[typePrefixAndSizeLen:]

	if len(data) < int(size) {
		return nil, errors.New("pubkey raw data size error")
	}
	if bytes.Compare(typePubKeySecp256k1Prefix, prefix) == 0 {
		if size != secp256k1.PubKeySecp256k1Size {
			return nil, errors.New("pubkey secp256k1 size error")
		}
		pubKey := secp256k1.PubKeySecp256k1{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else if bytes.Compare(typePubKeyEd25519Prefix, prefix) == 0 {
		if size != ed25519.PubKeyEd25519Size {
			return nil, errors.New("pubkey ed25519 size error")
		}
		pubKey := ed25519.PubKeyEd25519{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else if bytes.Compare(typePubKeySr25519Prefix, prefix) == 0 {
		if size != sr25519.PubKeySr25519Size {
			return nil, errors.New("pubkey sr25519 size error")
		}
		pubKey := sr25519.PubKeySr25519{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else {
		return nil, errors.New("unknown pubkey type")
	}
}
