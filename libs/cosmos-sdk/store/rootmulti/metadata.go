package rootmulti

import (
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"sync"
)

type (
	MetadataItem struct {
		version      int64
		cInfo        commitInfo
		pruneHeights []int64
		versions     []int64
	}

	Metadata struct {
		cache    MetadataCache
		db       dbm.DB
		task     chan MetadataItem
		quit     chan struct{}
		mut      sync.Mutex
		isClosed bool
	}

	MetadataCache struct {
		version int64
		//cInfo        commitInfo
		mtx          sync.RWMutex
		cInfos       map[int64]commitInfo
		deletions    map[int64]commitInfo
		pruneHeights []int64
		versions     []int64
	}
)

func NewMetadata(db dbm.DB) *Metadata {
	mt := &Metadata{
		cache: MetadataCache{
			version:      0,
			cInfos:       make(map[int64]commitInfo),
			deletions:    make(map[int64]commitInfo),
			pruneHeights: []int64{},
			versions:     []int64{},
		},
		db:       db,
		task:     make(chan MetadataItem, iavl.CommitGapHeight),
		quit:     make(chan struct{}),
		isClosed: false,
	}
	if iavl.EnableAsyncCommit {
		go func() {
			for {
				select {
				case mc, ok := <-mt.task:
					if ok {
						flushMetadata(mt.db, mc.version, mc.cInfo, mc.pruneHeights, mc.versions)
					} else {
						mt.quit <- struct{}{}
					}
				}
			}
		}()
	}
	return mt
}

func (mt *Metadata) CacheMetadata(version int64, cInfo commitInfo, pruneHeights []int64, versions []int64) {
	mt.cache.mtx.Lock()
	mt.cache.cInfos[version] = cInfo
	mt.cache.mtx.Unlock()
	mt.cache.version = version
	mt.cache.pruneHeights = pruneHeights
	mt.cache.versions = versions
}

func (mt *Metadata) AddCommitInfo(version int64, cInfo commitInfo) {
	mt.cache.mtx.Lock()
	defer mt.cache.mtx.Lock()
	mt.cache.cInfos[version] = cInfo
}

func (mt *Metadata) GetCommitInfo(version int64) (commitInfo, error) {
	mt.cache.mtx.RLock()
	defer mt.cache.mtx.RUnlock()
	cInfo, ok := mt.cache.cInfos[version]
	if !ok {
		cInfo, ok := mt.cache.deletions[version]
		if !ok {
			return cInfo, fmt.Errorf("metadata failed to GetCommitInfo")
		}
		return cInfo, nil
	}

	return cInfo, nil
}

func (mt *Metadata) GetLatestVersion() int64 {
	//todo
	return mt.cache.version
}

func (mt *Metadata) notifyFlushMetadata(version int64, cInfo commitInfo, pruneHeights []int64, versions []int64) {
	mt.mut.Lock()
	defer mt.mut.Unlock()

	if !mt.isClosed {
		mt.task <- MetadataItem{
			version, cInfo, pruneHeights, versions,
		}
	}
}

func (mt *Metadata) GracefulExit() {
	if iavl.EnableAsyncCommit {
		mt.mut.Lock()
		defer mt.mut.Unlock()
		if !mt.isClosed {
			mt.isClosed = true
			close(mt.task)
			<-mt.quit
		}
	}
}
