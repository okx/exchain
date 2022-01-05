package flatkv

import (
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const FlagEnable = "enable-flat-kv"

// Store wraps app_flat_kv.db for read performance.
type Store struct {
	db         dbm.DB
	cache      map[string][]byte
	readTime   int64
	writeTime  int64
	readCount  int64
	writeCount int64
	enable     bool
}

func NewStore(db dbm.DB) *Store {
	return &Store{
		db:         db,
		cache:      make(map[string][]byte),
		readTime:   0,
		writeTime:  0,
		readCount:  0,
		writeCount: 0,
		enable:     viper.GetBool(FlagEnable),
	}
}

func (st *Store) Get(key []byte) []byte {
	if !st.enable {
		return nil
	}
	if cacheVal, ok := st.getCache(key); ok {
		return cacheVal
	}
	ts := time.Now()
	value, err := st.db.Get(key)
	st.addDBReadTime(time.Now().Sub(ts).Nanoseconds())
	st.addDBReadCount()
	if err == nil && len(value) != 0 {
		return value
	}
	return nil
}

func (st *Store) Set(key, value []byte) {
	if !st.enable {
		return
	}
	st.addCache(key, value)
}

func (st *Store) Has(key []byte) bool {
	if !st.enable {
		return false
	}
	if _, ok := st.getCache(key); ok {
		return true
	}
	st.addDBReadCount()
	if ok, err := st.db.Has(key); err == nil && ok {
		return true
	}
	return false
}

func (st *Store) Delete(key []byte) {
	if !st.enable {
		return
	}
	ts := time.Now()
	st.db.Delete(key)
	st.addDBWriteTime(time.Now().Sub(ts).Nanoseconds())
	st.addDBWriteCount()
	st.deleteCache(key)
}

func (st *Store) Commit() {
	if !st.enable {
		return
	}
	ts := time.Now()
	// commit to flat kv db
	batch := st.db.NewBatch()
	defer batch.Close()
	for key, value := range st.cache {
		batch.Set([]byte(key), value)
	}
	batch.Write()
	st.addDBWriteTime(time.Now().Sub(ts).Nanoseconds())
	st.addDBWriteCount()
	// clear cache
	st.cache = make(map[string][]byte)
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

func (st *Store) getCache(key []byte) (value []byte, ok bool) {
	strKey := string(key)
	value, ok = st.cache[strKey]
	return
}

func (st *Store) addCache(key, value []byte) {
	strKey := string(key)
	st.cache[strKey] = value
}

func (st *Store) deleteCache(key []byte) {
	strKey := string(key)
	delete(st.cache, strKey)
}
