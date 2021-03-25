package state

import (
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const FlagFastQuery = "fast-query"

type stateStore struct {
	db *leveldb.DB
}

var gStateStore *stateStore = nil

func InstanceOfStateStore() *stateStore {
	if gStateStore == nil {
		db, e := initDb()
		if e == nil {
			gStateStore = &stateStore{db: db}
		} else {
			panic(e)
		}
	}
	return gStateStore
}

func initDb() (*leveldb.DB, error) {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, "data/storage.db")
	return leveldb.OpenFile(dbPath, nil)
}

func (s stateStore) Set(key []byte, value []byte) {
	s.db.Put(key, value, nil)
}

func (s stateStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}
