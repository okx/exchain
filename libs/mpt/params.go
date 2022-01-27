package mpt

import "github.com/okex/exchain/libs/cosmos-sdk/types"

const (
	StoreTypeMPT = types.StoreTypeMPT

	TriesInMemory = 100

	// StoreKey is string representation of the store key for mpt
	StoreKey = "mpt"

	FlagDBBackend = "db_backend"
)
