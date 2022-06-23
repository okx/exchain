package watcher

import (
	"log"
	"sync"

	dbm "github.com/okex/exchain/libs/tm-db"
)

const (
	FlagFastQuery = "fast-query"
)

type WatchStore struct {
	db *dbm.MemDB
}

var gWatchStore *WatchStore = nil
var once sync.Once

func InstanceOfWatchStore() *WatchStore {
	once.Do(func() {
		if IsWatcherEnabled() {
			gWatchStore = &WatchStore{db: dbm.NewMemDB()}
		}
	})
	return gWatchStore
}

func (w WatchStore) Set(key []byte, value []byte) {
	err := w.db.Set(key, value)
	if err != nil {
		log.Println("watchdb error: ", err.Error())
	}
}

func (w WatchStore) Get(key []byte) ([]byte, error) {
	return w.db.Get(key)
}

func (w WatchStore) GetUnsafe(key []byte, processor dbm.UnsafeValueProcessor) (interface{}, error) {
	return w.db.GetUnsafeValue(key, processor)
}

func (w WatchStore) Delete(key []byte) {
	err := w.db.Delete(key)
	if err != nil {
		log.Printf("watchdb error: " + err.Error())
	}
}

func (w WatchStore) Has(key []byte) bool {
	res, err := w.db.Has(key)
	if err != nil {
		log.Printf("watchdb error: " + err.Error())
		return false
	}
	return res
}

func (w WatchStore) Iterator(start, end []byte) dbm.Iterator {
	it, err := w.db.Iterator(start, end)
	if err != nil {
		log.Printf("watchdb error: " + err.Error())
		return nil
	}
	return it
}
