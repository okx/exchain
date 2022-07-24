package adapter

import (
	"errors"
	ethsecp256k12 "github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/ed25519"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/ethsecp256k1"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/multisig"
	secp256k1 "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/secp256k1"
	LegacyPubKey "github.com/okex/exchain/libs/tendermint/crypto"
	ed255192 "github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	multisig2 "github.com/okex/exchain/libs/tendermint/crypto/multisig"
	secp256k12 "github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
)

func LegacyPubkey2ProtoBuffPubkey(pubKey LegacyPubKey.PubKey) types.PubKey {
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

func ProtoBufPubkey2LegacyPubkey(pkData types.PubKey) (LegacyPubKey.PubKey, error) {
	var pubkey LegacyPubKey.PubKey
	switch v := pkData.(type) {
	case *ethsecp256k1.PubKey:
		pubkey = ethsecp256k12.PubKey(v.Bytes())
	case *secp256k1.PubKey:
		secpPk := &secp256k12.PubKeySecp256k1{}
		copy(secpPk[:], v.Bytes())
		pubkey = (LegacyPubKey.PubKey)(secpPk)
	case *multisig.LegacyAminoPubKey:
		var pubkeys []LegacyPubKey.PubKey
		for _, pubkeyItem := range v.PubKeys {
			item, ok := pubkeyItem.GetCachedValue().(types.PubKey)
			if !ok {
				return pubkey, errors.New("multisig not support pub key interface")
			}
			switch vv := item.(type) {
			case *ethsecp256k1.PubKey:
				pubkeys = append(pubkeys, ethsecp256k12.PubKey(vv.Bytes()))
			case *secp256k1.PubKey:
				secpPk := &secp256k12.PubKeySecp256k1{}
				copy(secpPk[:], vv.Bytes())
				pubkeys = append(pubkeys, (LegacyPubKey.PubKey)(secpPk))
			case *ed25519.PubKey:
				ed25519Pubkey := &ed255192.PubKeyEd25519{}
				copy(ed25519Pubkey[:], vv.Bytes())
				pubkeys = append(pubkeys, (LegacyPubKey.PubKey)(ed25519Pubkey))
			default:
				return pubkey, errors.New("multisig not support pub key type")
			}
		}
		pubkey = multisig2.NewPubKeyMultisigThreshold(int(v.Threshold), pubkeys)
	default:
		return pubkey, errors.New("not support pub key type")
	}
	return pubkey, nil
}
