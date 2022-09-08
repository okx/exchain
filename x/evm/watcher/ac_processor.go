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
	gstat             = make(map[string]*Stat)
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
		// for test
		go StaticMemory()
	})
}

// ACProcessor async commit processor
type ACProcessor struct {
	mtx         sync.RWMutex
	curMsgCache *MessageCache
	commitList  *commitCache
	totalCommit int
	totalRepeat int
	total       int
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
		curstat := cmmiter.Clear()
		ed2 := time.Now()
		for k, v := range curstat {
			value, ok := gstat[k]
			if !ok {
				gstat[k] = &Stat{count: v.count, dbSize: v.dbSize, structSize: v.structSize}
			} else {
				value.count += v.count
				value.dbSize += v.dbSize
				value.structSize += v.structSize
				gstat[k] = value
			}
		}

		ap.totalCommit += s
		ap.totalRepeat += cmmiter.count - s
		ap.total += cmmiter.count
		fmt.Printf("****** lyh ACProcessor cur commiter size %d, repeat count %d;;; total commit %d, repeat %d, total %d;;; cost time commit %v, remove %v, clear %v \n",
			s, cmmiter.count-s,
			ap.totalCommit, ap.totalRepeat, ap.total,
			ed.Sub(st), ed1.Sub(ed), ed2.Sub(ed1))
	}

	var totaldbsize int
	var totalstructsize int
	for k, v := range gstat {
		dbsize := float64(v.dbSize) / float64(1024*1024)
		structsize := float64(v.structSize) / float64(1024*1024)
		totaldbsize += v.dbSize
		totalstructsize += v.structSize
		fmt.Printf("**** lyh ****** glbal static %s, count %d, dbSize %.3f, structSize %.3f \n", k, v.count, dbsize, structsize)
	}
	fmt.Printf("**** lyh ****** glbal static total dbSize %.3f, structSize %.3f \n", float64(totaldbsize)/float64(1024*1024), float64(totalstructsize)/float64(1024*1024))
}

// PersistHander after close channel should call this function
func (ap *ACProcessor) Close(version int64, commitFn func(epochCache *MessageCache)) {
	ap.MoveToCommitList(version)
	ap.PersistHander(commitFn)
}
