package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
	"log"
	"path/filepath"
	"sync"
)

var gFlatStore *FlatStore = nil
var flatOnce sync.Once

const (
	FlagDBBackend = "db_backend"
	FlatDBDir     = "data"
	FlatDBName    = "flat"
)

type FlatStore struct {
	db dbm.DB
}

func InstanceOfFlatStore() *FlatStore {
	flatOnce.Do(func() {
		gFlatStore = &FlatStore{db: initFlatDb()}
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

	//backend = string(dbm.MemDBBackend)

	return dbm.NewDB(FlatDBName, dbm.BackendType(backend), dbPath)
}

func (w FlatStore) Set(key []byte, value []byte) {
	err := w.db.Set(key, value)
	if err != nil {
		log.Println("flatdb error: ", err.Error())
	}
}

func (w FlatStore) Get(key []byte) ([]byte, error) {
	return w.db.Get(key)
}

func (w FlatStore) Delete(key []byte) {
	err := w.db.Delete(key)
	if err != nil {
		log.Printf("flatdb error: " + err.Error())
	}
}

func (w FlatStore) Has(key []byte) bool {
	res, err := w.db.Has(key)
	if err != nil {
		log.Println("flatdb error: " + err.Error())
		return false
	}
	return res
}

func (w FlatStore) Iterator(start, end []byte) dbm.Iterator {
	it, err := w.db.Iterator(start, end)
	if err != nil {
		log.Println("flatdb error: " + err.Error())
		return nil
	}
	return it
}
