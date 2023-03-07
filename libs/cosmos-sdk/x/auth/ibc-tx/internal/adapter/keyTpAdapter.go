package adapter

import (
	"errors"
	ethsecp256k12 "github.com/okx/okbchain/app/crypto/ethsecp256k1"
	"github.com/okx/okbchain/libs/cosmos-sdk/crypto/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/ethsecp256k1"
	secp256k1 "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/secp256k1"
	LagacyPubKey "github.com/okx/okbchain/libs/tendermint/crypto"
	secp256k12 "github.com/okx/okbchain/libs/tendermint/crypto/secp256k1"
)

func LagacyPubkey2ProtoBuffPubkey(pubKey LagacyPubKey.PubKey) types.PubKey {
	var pubKeyPB types.PubKey

	switch v := pubKey.(type) {
	case ethsecp256k12.PubKey:
		ethsecp256k1Pk := &ethsecp256k1.PubKey{
			Key: v,
		}
		pubKeyPB = ethsecp256k1Pk
	case secp256k12.PubKeySecp256k1:
		secp256k1Pk := &secp256k1.PubKey{
			Key: v[:],
		}
		pubKeyPB = secp256k1Pk
	default:
		panic("not supported key algo")
	}
	return pubKeyPB
}

func ProtoBufPubkey2LagacyPubkey(pkData types.PubKey) (LagacyPubKey.PubKey, error) {
	var pubkey LagacyPubKey.PubKey
	switch v := pkData.(type) {
	case *ethsecp256k1.PubKey:
		pubkey = ethsecp256k12.PubKey(v.Bytes())
	case *secp256k1.PubKey:
		secpPk := &secp256k12.PubKeySecp256k1{}
		copy(secpPk[:], v.Bytes())
		pubkey = (LagacyPubKey.PubKey)(secpPk)
	default:
		return pubkey, errors.New("not support pub key type")
	}
	return pubkey, nil
}
