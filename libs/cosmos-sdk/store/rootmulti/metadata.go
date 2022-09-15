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
		rwLock   sync.RWMutex
		mtCache  map[int64]MetadataItem
		db       dbm.DB
		task     chan MetadataItem
		quit     chan struct{}
		acFlush  chan struct{}
		mut      sync.Mutex
		isClosed bool
	}
)

func NewMetadata(db dbm.DB) *Metadata {
	mt := &Metadata{
		mtCache:  make(map[int64]MetadataItem),
		db:       db,
		task:     make(chan MetadataItem, iavl.CommitGapHeight/2),
		quit:     make(chan struct{}),
		acFlush:  make(chan struct{}),
		isClosed: false,
	}
	if iavl.EnableAsyncCommit {
		go func() {
			for {
				select {
				case mc, ok := <-mt.task:
					if ok {
						flushMetadata(mt.db, mc.version, mc.cInfo, mc.pruneHeights, mc.versions)
						mt.unCacheMt(mc.version)
						if iavl.ShouldPersist(mc.version) {
							//must wait flush over
							mt.acFlush <- struct{}{}
						}
					} else {
						mt.quit <- struct{}{}
					}
				}
			}
		}()
	}
	return mt
}
func (mt *Metadata) cacheMt(version int64, mtItem MetadataItem) {
	mt.rwLock.Lock()
	defer mt.rwLock.Unlock()
	mt.mtCache[version] = mtItem
}
func (mt *Metadata) unCacheMt(version int64) {
	mt.rwLock.Lock()
	defer mt.rwLock.Unlock()
	delete(mt.mtCache, version)
}

func (mt *Metadata) GetCommitInfoFromCache(version int64) (commitInfo, error) {
	mt.rwLock.RLock()
	defer mt.rwLock.RUnlock()
	mtItem, ok := mt.mtCache[version]
	if !ok {
		return commitInfo{}, fmt.Errorf("no commitInfo from cache")
	}
	return mtItem.cInfo, nil
}

func (mt *Metadata) notifyFlushMetadata(version int64, cInfo commitInfo, pruneHeights []int64, versions []int64) {
	mt.cacheMt(version, MetadataItem{
		version, cInfo, pruneHeights, versions,
	})

	mt.mut.Lock()
	defer mt.mut.Unlock()

	if !mt.isClosed {
		mt.task <- MetadataItem{
			version, cInfo, pruneHeights, versions,
		}
	}
	if iavl.ShouldPersist(version) {
		<-mt.acFlush
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
