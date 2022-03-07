package codec

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"

	ibckey "github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/ibc-key"

	//"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/ed25519"
	//"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/multisig"
	//"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/secp256k1"
	//"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterInterface("cosmos.crypto.PubKey", pk)
	//registry.RegisterImplementations(pk, &ed25519.PubKey{})
	registry.RegisterImplementations(pk, &ibckey.PubKey{})
	//registry.RegisterImplementations(pk, &multisig.LegacyAminoPubKey{})
	//secp256r1.RegisterInterfaces(registry)
}
