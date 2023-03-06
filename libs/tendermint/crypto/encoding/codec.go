package encoding

import (
	"errors"
	"fmt"

	"github.com/okx/okbchain/libs/tendermint/crypto"
	"github.com/okx/okbchain/libs/tendermint/crypto/ed25519"
	"github.com/okx/okbchain/libs/tendermint/crypto/secp256k1"
	pc "github.com/okx/okbchain/libs/tendermint/proto/crypto/keys"
)

type PubKeyType uint8

const (
	Unknown PubKeyType = iota
	Ed25519
	Secp256k1
)

// PubKeyToProto takes crypto.PubKey and transforms it to a protobuf Pubkey
func PubKeyToProto(k crypto.PubKey) (pc.PublicKey, error) {
	if k == nil {
		return pc.PublicKey{}, errors.New("nil PublicKey")
	}
	var kp pc.PublicKey
	switch k := k.(type) {
	case ed25519.PubKeyEd25519:
		kp = pc.PublicKey{
			Sum: &pc.PublicKey_Ed25519{
				Ed25519: k[:],
			},
		}
	case secp256k1.PubKeySecp256k1:
		kp = pc.PublicKey{
			Sum: &pc.PublicKey_Secp256K1{
				Secp256K1: k[:],
			},
		}
	default:
		return kp, fmt.Errorf("toproto: key type %v is not supported", k)
	}
	return kp, nil
}

// PubKeyFromProto takes a protobuf Pubkey and transforms it to a crypto.Pubkey
// Return one more parameter to prevent of slowing down the whole procedure
func PubKeyFromProto(k *pc.PublicKey) (crypto.PubKey, PubKeyType, error) {
	if k == nil {
		return nil, Unknown, errors.New("nil PublicKey")
	}
	switch k := k.Sum.(type) {
	case *pc.PublicKey_Ed25519:
		if len(k.Ed25519) != ed25519.PubKeyEd25519Size {
			return nil, Unknown, fmt.Errorf("invalid size for PubKeyEd25519. Got %d, expected %d",
				len(k.Ed25519), ed25519.PubKeyEd25519Size)
		}
		var pk ed25519.PubKeyEd25519
		copy(pk[:], k.Ed25519)
		return pk, Ed25519, nil
	case *pc.PublicKey_Secp256K1:
		if len(k.Secp256K1) != secp256k1.PubKeySecp256k1Size {
			return nil, Unknown, fmt.Errorf("invalid size for PubKeySecp256k1. Got %d, expected %d",
				len(k.Secp256K1), secp256k1.PubKeySecp256k1Size)
		}
		var pk secp256k1.PubKeySecp256k1
		copy(pk[:], k.Secp256K1)
		return pk, Secp256k1, nil
	default:
		return nil, Unknown, fmt.Errorf("fromproto: key type %v is not supported", k)
	}
}
