package types

import (
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var CommitStateDBPool = sync.Pool{
	New: func() interface{} {
		return &CommitStateDB{
			stateObjects:        make(map[ethcmn.Address]*stateObject),
			stateObjectsPending: make(map[ethcmn.Address]struct{}),
			stateObjectsDirty:   make(map[ethcmn.Address]struct{}),
			preimages:           make(map[ethcmn.Hash][]byte),
			logs:                make(map[ethcmn.Hash][]*ethtypes.Log),
			codeCache:           make(map[ethcmn.Address]CacheCode, 0),
			updatedAccount:      make(map[ethcmn.Address]struct{}),
		}
	},
}
