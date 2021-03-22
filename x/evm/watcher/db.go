package watcher

import (
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const FlagFastQuery = "fast-query"

type WatchStore struct {
	db *leveldb.DB
}

var gWatchStore *WatchStore = nil

func InstanceOfWatchStore() *WatchStore {
	if gWatchStore == nil {
		db, e := initDb()
		if e == nil {

			gWatchStore = &WatchStore{db: db}
		}
	}
	return gWatchStore
}

func initDb() (*leveldb.DB, error) {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, "data/watch.db")
	return leveldb.OpenFile(dbPath, nil)
}

func (w WatchStore) Set(key []byte, value []byte) {
	w.db.Put(key, value, nil)
}

func (w WatchStore) Get(key []byte) ([]byte, error) {
	return w.db.Get(key, nil)
}
