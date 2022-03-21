package mpt

import (
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
	KeyPrefixRootMptHash             = []byte{0x01}
	KeyPrefixLatestStoredHeight      = []byte{0x02}
	AccStoreCache               uint = 2048 // MB

	GAccToPrefetchChannel    = make(chan [][]byte, 2000)
	GAccTryUpdateTrieChannel = make(chan struct{})
	GAccTrieUpdatedChannel = make(chan struct{})
)
