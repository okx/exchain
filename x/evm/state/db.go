package state

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/trie"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/ethereum/go-ethereum/core/state"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
)

type stateStore struct {
	db state.Database
}

var gStateStore *stateStore = nil

func InstanceOfStateStore() *stateStore {
	if gStateStore == nil {
		homeDir := viper.GetString(flags.FlagHome)
		dbPath := filepath.Join(homeDir, "data/storage.db")
		//set cache and handle value as a test number
		db, e := rawdb.NewLevelDBDatabase(dbPath, 1024, 102400, "evmState")
		if e == nil {
			gStateStore = &stateStore{db: state.NewDatabase(db)}
		}

	}
	return gStateStore
}

func (s stateStore) GetDb() state.Database {
	return s.db
}

func (s stateStore) PruningTrie(root ethcmn.Hash) error {
	t, e := s.db.OpenTrie(root)
	if e != nil {
		return e
	}
	disk := s.db.TrieDB().DiskDB()
	it := trie.NewIterator(t.NodeIterator(nil))
	for it.Next() {
		e := disk.Delete(it.Key)
		if e != nil {
			return e
		}
	}
	e = disk.Delete(root.Bytes())
	return e
}
