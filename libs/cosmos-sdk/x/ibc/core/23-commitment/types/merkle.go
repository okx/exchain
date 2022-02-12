package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
	"net/url"
)

var _ exported.Root = (*MerkleRoot)(nil)

// String implements fmt.Stringer.
// This represents the path in the same way the tendermint KeyPath will
// represent a key path. The backslashes partition the key path into
// the respective stores they belong to.
func (mp MerklePath) String() string {
	pathStr := ""
	for _, k := range mp.KeyPath {
		pathStr += "/" + url.PathEscape(k)
	}
	return pathStr
}


// NewMerkleRoot constructs a new MerkleRoot
func NewMerkleRoot(hash []byte) MerkleRoot {
	return MerkleRoot{
		Hash: hash,
	}
}


// GetHash implements RootI interface
func (mr MerkleRoot) GetHash() []byte {
	return mr.Hash
}

// Empty returns true if the root is empty
func (mr MerkleRoot) Empty() bool {
	return len(mr.GetHash()) == 0
}

var _ exported.Prefix = (*MerklePrefix)(nil)

// NewMerklePrefix constructs new MerklePrefix instance
func NewMerklePrefix(keyPrefix []byte) MerklePrefix {
	return MerklePrefix{
		KeyPrefix: keyPrefix,
	}
}

// Bytes returns the key prefix bytes
func (mp MerklePrefix) Bytes() []byte {
	return mp.KeyPrefix
}

// Empty returns true if the prefix is empty
func (mp MerklePrefix) Empty() bool {
	return len(mp.Bytes()) == 0
}
