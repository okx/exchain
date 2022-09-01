package watcher

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
)

const (
	FlagWatchdbEnableAsyncCommit = "watchdb-enable-async-commit"
	FlagWatchdbCommitGapHeight   = "watchdb-commit-gap-height"
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

func SetAnableAsyncCommit(enable bool) {
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
		SetAnableAsyncCommit(viper.GetBool(FlagWatchdbEnableAsyncCommit))
		gap := viper.GetInt64(FlagWatchdbCommitGapHeight)
		if gap != 0 {
			SetCommitGapHeight(gap)
		}
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

func (ap *ACProcessor) PersistHander(commitFn func(epochCache *MessageCache)) {
	// commit to db
	for {
		cmmiter, ok := ap.commitList.getTop()
		if !ok {
			break
		}
		s := len(cmmiter.MessageCache.mp)
		st := time.Now()
		commitFn(cmmiter.MessageCache)
		ed := time.Now()
		ap.commitList.remove(cmmiter.version)
		ed1 := time.Now()
		cmmiter.Clear()
		ed2 := time.Now()
		fmt.Printf("****** lyh ACProcessor cur commiter size %d, repeat count %d, commitlist len %d;;; cost time commit %v, remove %v, clear %v \n",
			s, cmmiter.count-s, ap.commitList.size(), ed.Sub(st), ed1.Sub(ed), ed2.Sub(ed1))
	}
}

// PersistHander after close channel should call this function
func (ap *ACProcessor) Close(version int64, commitFn func(epochCache *MessageCache)) {
	ap.MoveToCommitList(version)
	ap.PersistHander(commitFn)
}
