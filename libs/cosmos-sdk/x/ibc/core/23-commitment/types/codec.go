package types

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)

// RegisterInterfaces registers the commitment interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"oec.ibc.core.commitment.v1.Root",
		(*exported.Root)(nil),
	)
	registry.RegisterInterface(
		"oec.ibc.core.commitment.v1.Prefix",
		(*exported.Prefix)(nil),
	)
	registry.RegisterInterface(
		"oec.ibc.core.commitment.v1.Path",
		(*exported.Path)(nil),
	)
	registry.RegisterInterface(
		"oec.ibc.core.commitment.v1.Proof",
		(*exported.Proof)(nil),
	)

	registry.RegisterImplementations(
		(*exported.Root)(nil),
		&MerkleRoot{},
	)
	registry.RegisterImplementations(
		(*exported.Prefix)(nil),
		&MerklePrefix{},
	)
	registry.RegisterImplementations(
		(*exported.Path)(nil),
		&MerklePath{},
	)
	registry.RegisterImplementations(
		(*exported.Proof)(nil),
		&MerkleProof{},
	)
}
