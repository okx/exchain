package types

import (
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
)

/*
Keeps track of blockParts cache from round 0 to round 'round'.
*/
type HeightBPCache struct {
	mtx    sync.Mutex
	height int64
	cache  map[int]*types.Part
	count  int
}

func NewBPCache(height int64) *HeightBPCache {
	if height <= 0 {
		panic("Cannot make BPCache for height <= 0, doesn't make sense.")
	}

	hbc := &HeightBPCache{
		height: height,
		cache:  make(map[int]*types.Part),
	}

	return hbc
}

// Duplicate votes return added=false, err=nil.
// By convention, peerID is "" if origin is self.
func (hbc *HeightBPCache) AddBlockPart(height int64, part *types.Part) {
	hbc.mtx.Lock()
	defer hbc.mtx.Unlock()
	if hbc.height != height {
		return
	}

	if hbc.cache[part.Index] == nil {
		hbc.count++
	}

	hbc.cache[part.Index] = part

	return
}

func (hbc *HeightBPCache) Count() int {
	hbc.mtx.Lock()
	defer hbc.mtx.Unlock()
	return hbc.count
}

func (hbc *HeightBPCache) Height() int64 {
	hbc.mtx.Lock()
	defer hbc.mtx.Unlock()
	return hbc.height
}
func (hbc *HeightBPCache) Cache() map[int]*types.Part {
	hbc.mtx.Lock()
	defer hbc.mtx.Unlock()
	return hbc.cache
}

func (hbc *HeightBPCache) GetPart(index int) *types.Part {
	hbc.mtx.Lock()
	defer hbc.mtx.Unlock()
	return hbc.cache[index]
}
