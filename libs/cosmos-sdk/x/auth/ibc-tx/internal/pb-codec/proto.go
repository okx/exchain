package codec

import (
	codectypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/ethsecp256k1"
	ibckey "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/secp256k1"

	//"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys/ed25519"
	//"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys/multisig"
	//"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys/secp256k1"
	//"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/okx/okbchain/libs/cosmos-sdk/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterInterface("cosmos.crypto.PubKey", pk)
	registry.RegisterInterface("ethermint.crypto.v1.ethsecp256k1.PubKey", pk)
	registry.RegisterImplementations(pk, &ibckey.PubKey{})
	registry.RegisterImplementations(pk, &ethsecp256k1.PubKey{})
}
