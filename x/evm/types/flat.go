package types

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
	"log"
	"path/filepath"
	"sync"
)

var gFlatStore *FlatStore = nil
var flatOnce sync.Once
var EnableFlatDB = false

const (
	FlagDBBackend = "db_backend"
	FlatDBDir     = "data"
	FlatDBName    = "flat"

	FlatEnableFlagDB = "evm-use-flatdb"
)

type FlatStore struct {
	db      dbm.DB
	kvCache *fastcache.Cache
}

func InstanceOfFlatStore() *FlatStore {
	flatOnce.Do(func() {
		gFlatStore = &FlatStore{
			db:      initFlatDb(),
			kvCache: fastcache.New(2 * 1024 * 1024 * 1024),
		}
	})
	return gFlatStore
}

func initFlatDb() dbm.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, FlatDBDir)
	backend := viper.GetString(FlagDBBackend)
	if backend == "" {
		backend = string(dbm.GoLevelDBBackend)
	}

	return dbm.NewDB(FlatDBName, dbm.BackendType(backend), dbPath)
}

func (w *FlatStore) Set(key []byte, value []byte) {
	err := w.db.Set(key, value)
	if err != nil {
		log.Println("flatdb error: ", err.Error())
	}
	w.kvCache.Set(key, value)

}

func (w *FlatStore) Get(key []byte) ([]byte, error) {
	if w.kvCache.Has(key) {
		return w.kvCache.Get(nil, key), nil
	} else {
		return w.db.Get(key)
	}
}

func (w *FlatStore) Delete(key []byte) {
	err := w.db.Delete(key)
	if err != nil {
		log.Printf("flatdb error: " + err.Error())
	}
	w.kvCache.Del(key)
}

func (w *FlatStore) Has(key []byte) bool {
	res, err := w.db.Has(key)
	if err != nil {
		log.Println("flatdb error: " + err.Error())
		return false
	}
	return res
}

func (w *FlatStore) Iterator(start, end []byte) dbm.Iterator {
	it, err := w.db.Iterator(start, end)
	if err != nil {
		log.Println("flatdb error: " + err.Error())
		return nil
	}
	return it
}

type FlatBatch struct {
	db    *FlatStore
	batch dbm.Batch
}

func (w *FlatStore) NewBatch() *FlatBatch {
	return &FlatBatch{
		db:    w,
		batch: w.db.NewBatch(),
	}
}

// Write writes the batch, possibly without flushing to disk. Only Close() can be called after,
// other methods will panic.
func (w *FlatBatch) Write() error {
	return w.batch.Write()
}

// WriteSync writes the batch and flushes it to disk. Only Close() can be called after, other
// methods will panic.
func (w *FlatBatch) WriteSync() error {
	return w.batch.WriteSync()
}

// Close closes the batch. It is idempotent, but any other calls afterwards will panic.
func (w *FlatBatch) Close() {
	w.batch.Close()
}

func (w *FlatBatch) Set(key, value []byte) {
	//w.batch.Set(key, value)
	w.db.kvCache.Set(key, value)
}

// Delete deletes a key/value pair.
// CONTRACT: key readonly []byte
func (w *FlatBatch) Delete(key []byte) {
	//w.batch.Delete(key)
	w.db.kvCache.Del(key)
}
