package mpt

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okx/okbchain/libs/cosmos-sdk/types"
)

const (
	StoreTypeMPT = types.StoreTypeMPT

	TriesInMemory = 100

	// StoreKey is string representation of the store key for mpt
	StoreKey = "mpt"
)

const (
	FlagTrieDirtyDisabled = "trie.dirty-disabled"
	FlagTrieCacheSize     = "trie.cache-size"
	FlagTrieNodesLimit    = "trie.nodes-limit"
	FlagTrieImgsLimit     = "trie.imgs-limit"
)

var (
	TrieDirtyDisabled       = false
	TrieCacheSize     uint  = 2048 // MB
	TrieNodesLimit    uint  = 256  // MB
	TrieImgsLimit     uint  = 4    // MB
	TrieCommitGap     int64 = 100
)

var (
	KeyPrefixAccRootMptHash        = []byte{0x11}
	KeyPrefixAccLatestStoredHeight = []byte{0x12}

	GAccToPrefetchChannel    = make(chan [][]byte, 2000)
	GAccTryUpdateTrieChannel = make(chan struct{})
	GAccTrieUpdatedChannel   = make(chan struct{})
)

var (
	NilHash = ethcmn.Hash{}

	// EmptyCodeHash is the known hash of an empty code.
	EmptyCodeHash      = crypto.Keccak256Hash(nil)
	EmptyCodeHashBytes = crypto.Keccak256(nil)
)
