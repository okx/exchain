package types

import (
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
)

type BlockPartsCache struct {
	parts map[int]*types.Part
}

func newBlockPartsCache() *BlockPartsCache {
	return &BlockPartsCache{
		parts: make(map[int]*types.Part),
	}
}

func (bpc *BlockPartsCache)addBlockPart(part *types.Part) bool {
	if bpc.parts[part.Index] == nil {
		bpc.parts[part.Index] = part
		return true
	}
	return false
}

type BlockPartsCacheMap struct {
	mtx sync.Mutex
	cache map[int64]map[int]*BlockPartsCache
}

func NewBlockPartsCacheMap() *BlockPartsCacheMap {
	cm := &BlockPartsCacheMap{}
	cm.cache = make(map[int64]map[int]*BlockPartsCache)

	return cm
}

func (cm *BlockPartsCacheMap) AddBlockPart(height int64, round int, part *types.Part) bool {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()

	if cm.cache[height] == nil {
		cm.cache[height] = make(map[int]*BlockPartsCache)
	}

	heightCache := cm.cache[height]
	if heightCache[round] == nil {
		heightCache[round] = newBlockPartsCache()
	}

	roundCache := heightCache[round]

	return roundCache.addBlockPart(part)
}

func (cm *BlockPartsCacheMap) GetBlockParts(height int64, round int)  {
	
}