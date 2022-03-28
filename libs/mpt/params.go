package mpt

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	StoreTypeMPT = types.StoreTypeMPT

	TriesInMemory = 100

	// StoreKey is string representation of the store key for mpt
	StoreKey = "mpt"

	FlagAccStoreCache = "account-store-cache"
)

var (
	KeyPrefixRootMptHash             = []byte{0x11}
	KeyPrefixLatestStoredHeight      = []byte{0x12}
	AccStoreCache               uint = 2048 // MB

	GAccToPrefetchChannel    = make(chan [][]byte, 2000)
	GAccTryUpdateTrieChannel = make(chan struct{})
	GAccTrieUpdatedChannel   = make(chan struct{})
)

var (
	NilHash = ethcmn.Hash{}

	// EmptyCodeHash is the known hash of an empty code.
	EmptyCodeHash      = crypto.Keccak256Hash(nil)
	EmptyCodeHashBytes = crypto.Keccak256(nil)

	// EmptyRootHash is the known root hash of an empty trie.
	EmptyRootHash      = ethtypes.EmptyRootHash
	EmptyRootHashBytes = EmptyRootHash.Bytes()
)
