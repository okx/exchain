package state

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/ethereum/go-ethereum/core/state"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
)

const FlagFastQuery = "fast-query"

type stateStore struct {
	db state.Database
}

var gStateStore *stateStore = nil

func InstanceOfStateStore() *stateStore {
	if gStateStore == nil {
		homeDir := viper.GetString(flags.FlagHome)
		dbPath := filepath.Join(homeDir, "data/storage.db")
		db, e := rawdb.NewLevelDBDatabase(dbPath, 1024, 102400, "state")
		if e == nil {
			gStateStore = &stateStore{db: state.NewDatabase(db)}
		}

	}
	return gStateStore
}

func (s stateStore) GetDb() state.Database {
	return s.db
}

func (s stateStore) Set(key []byte, value []byte) {

}

func (s stateStore) Get(key []byte) ([]byte, error) {
	return nil, nil
}
