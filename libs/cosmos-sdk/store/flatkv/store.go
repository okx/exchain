package flatkv

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/iavl"

	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	latestVersionKey = "s/latest"

	FlagEnable = "enable-flat-kv"
)

// Store wraps app_flat_kv.db for read performance.
type Store struct {
	db          dbm.DB
	cache       *Cache
	readTime    int64
	writeTime   int64
	readCount   int64
	writeCount  int64
	enable      bool
	asyncCommit bool
	tree        Tree
	preloadCh   chan []byte
}

func NewStore(db dbm.DB, tree Tree) *Store {
	st := &Store{
		db:          db,
		cache:       newCache(),
		readTime:    0,
		writeTime:   0,
		readCount:   0,
		writeCount:  0,
		enable:      viper.GetBool(FlagEnable),
		asyncCommit: iavl.EnableAsyncCommit,
		tree:        tree,
		preloadCh:   make(chan []byte, 0),
	}
	if st.enable {
		st.preloadSchedule()
	}
	return st
}

func (st *Store) preloadSchedule() {
	num := runtime.NumCPU()
	for i := 0; i < num; i++ {
		go func() {
			for key := range st.preloadCh {
				go st.tree.Get(key)
			}
		}()
	}
}

func (st *Store) Enable() bool {
	return st.enable
}
func (st *Store) Get(key []byte) []byte {
	if !st.enable {
		return nil
	}
	if cacheVal, ok := st.cache.get(key); ok {
		return cacheVal
	}
	ts := time.Now()
	value, err := st.db.Get(key)
	st.addDBReadTime(time.Now().Sub(ts).Nanoseconds())
	st.addDBReadCount()
	if err == nil && len(value) != 0 {
		st.cache.add(key, value, false, false)
		st.preloadCh <- key
		return value
	}
	return nil
}

func (st *Store) Set(key, value []byte) {
	if !st.enable {
		return
	}
	st.cache.add(key, value, false, true)
}

func (st *Store) Has(key []byte) bool {
	if !st.enable {
		return false
	}
	value := st.Get(key)
	return value != nil
}

func (st *Store) Delete(key []byte) {
	if !st.enable {
		return
	}
	st.cache.add(key, nil, true, true)
}

func (st *Store) Commit(version int64) {
	if !st.enable {
		return
	}
	if !st.asyncCommit {
		st.write(version)
		return
	}

	if st.tree.ShouldPersist(version) {
		go st.write(version)
	}
}

// writes cache to db when app is stopped and async commit is enabled
func (st *Store) Stop() {
	if !st.enable || !st.asyncCommit {
		return
	}
	st.write(0)
}

func (st *Store) write(version int64) {
	ts := time.Now()
	// commit to flat kv db
	batch := st.db.NewBatch()
	defer batch.Close()
	st.cache.write(batch)
	st.setLatestVersion(batch, version)
	batch.Write()
	st.addDBWriteTime(time.Now().Sub(ts).Nanoseconds())
	st.addDBWriteCount()
}

func (st *Store) ResetCount() {
	if !st.enable {
		return
	}
	st.resetDBReadTime()
	st.resetDBWriteTime()
	st.resetDBReadCount()
	st.resetDBWriteCount()
}

func (st *Store) GetDBReadTime() int {
	if !st.enable {
		return 0
	}
	return int(atomic.LoadInt64(&st.readTime))
}

func (st *Store) addDBReadTime(ts int64) {
	atomic.AddInt64(&st.readTime, ts)
}

func (st *Store) resetDBReadTime() {
	atomic.StoreInt64(&st.readTime, 0)
}

func (st *Store) GetDBWriteTime() int {
	if !st.enable {
		return 0
	}
	return int(atomic.LoadInt64(&st.writeTime))
}

func (st *Store) addDBWriteTime(ts int64) {
	atomic.AddInt64(&st.writeTime, ts)
}

func (st *Store) resetDBWriteTime() {
	atomic.StoreInt64(&st.writeTime, 0)
}

func (st *Store) GetDBReadCount() int {
	if !st.enable {
		return 0
	}
	return int(atomic.LoadInt64(&st.readCount))
}

func (st *Store) addDBReadCount() {
	atomic.AddInt64(&st.readCount, 1)
}

func (st *Store) resetDBReadCount() {
	atomic.StoreInt64(&st.readCount, 0)
}

func (st *Store) GetDBWriteCount() int {
	if !st.enable {
		return 0
	}
	return int(atomic.LoadInt64(&st.writeCount))
}

func (st *Store) addDBWriteCount() {
	atomic.AddInt64(&st.writeCount, 1)
}

func (st *Store) resetDBWriteCount() {
	atomic.StoreInt64(&st.writeCount, 0)
}

func (st *Store) GetLatestVersion() int64 {
	if !st.enable {
		return 0
	}
	return getLatestVersion(st.db)
}

func getLatestVersion(db dbm.DB) int64 {
	var latest int64
	latestBytes, err := db.Get([]byte(latestVersionKey))
	if err != nil {
		panic(err)
	} else if latestBytes == nil {
		return 0
	}

	err = cdc.UnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}

	return latest
}

func (st *Store) setLatestVersion(batch dbm.Batch, version int64) {
	latestBytes := cdc.MustMarshalBinaryLengthPrefixed(version)
	batch.Set([]byte(latestVersionKey), latestBytes)
}
