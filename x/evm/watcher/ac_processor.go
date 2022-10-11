package watcher

import (
	"github.com/spf13/viper"
	"sync"
)

const (
	FlagWatchdbEnableAsyncCommit = "watchdb-enable-async-commit"
)

var (
	gACProcessor      *ACProcessor // global variable
	enableAsyncCommit              = false
	commitGapHeight   int64        = 100
	getFlagOnce       sync.Once
)

func init() {
	gACProcessor = &ACProcessor{
		commitList:  newCommitCache(), // for support to querier
		curMsgCache: newMessageCache(),
	}
}

func SetEnableAsyncCommit(enable bool) {
	enableAsyncCommit = enable
}

func GetEnableAsyncCommit() bool {
	GetACFlag()
	return enableAsyncCommit
}

func SetCommitGapHeight(height int64) {
	commitGapHeight = height
}

func GetCommitGapHeight() int64 {
	GetACFlag()
	return commitGapHeight
}

func GetACFlag() {
	getFlagOnce.Do(func() {
		SetEnableAsyncCommit(viper.GetBool(FlagWatchdbEnableAsyncCommit))
	})
}

// ACProcessor async commit processor
type ACProcessor struct {
	mtx         sync.RWMutex
	curMsgCache *MessageCache
	commitList  *commitCache
}

func (ap *ACProcessor) BatchDel(keys [][]byte) {
	ap.curMsgCache.BatchDel(keys)
}

func (ap *ACProcessor) BatchSet(wsgs []WatchMessage) {
	ap.curMsgCache.BatchSet(wsgs)
}

func (ap *ACProcessor) BatchSetEx(batchs []*Batch) {
	ap.curMsgCache.BatchSetEx(batchs)
}

//  Get  the return value of interface{} should pay attend to delta data type is []byte or Message
func (ap *ACProcessor) Get(key []byte) (WatchMessage, bool) {
	// from current message cache get key
	ap.mtx.RLock()
	if ap.curMsgCache != nil {
		if v, ok := ap.curMsgCache.Get(key); ok {
			ap.mtx.RUnlock()
			return v, ok
		}
	}
	ap.mtx.RUnlock()

	// from commitlist get key
	if v, ok := ap.commitList.getElementFromCache(key); ok {
		return v, ok
	}
	return nil, false
}

func (ap *ACProcessor) MoveToCommitList(version int64) {
	cur := ap.curMsgCache
	if cur != nil {
		ap.commitList.pushBack(version, &MessageCacheEvent{cur, version})
	}

	ap.mtx.Lock()
	ap.curMsgCache = newMessageCache()
	ap.mtx.Unlock()
}

func (ap *ACProcessor) PersistHandler(commitFn func(epochCache *MessageCache)) {
	// commit to db
	for {
		committer, ok := ap.commitList.getTop()
		if !ok {
			break
		}
		commitFn(committer.MessageCache)
		ap.commitList.remove(committer.version)
		committer.Clear()
	}
}

// PersistHandler after close channel should call this function
func (ap *ACProcessor) Close(version int64, commitFn func(epochCache *MessageCache)) {
	ap.MoveToCommitList(version)
	ap.PersistHandler(commitFn)
}
