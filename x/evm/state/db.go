package state

import (
	"fmt"
	"path/filepath"

	"github.com/status-im/keycard-go/hexutils"

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

func (s stateStore) PruningTrie(root ethcmn.Hash, dirtyKeys map[string]int) error {
	t, e := s.db.OpenTrie(root)
	if e != nil {
		return e
	}
	disk := s.db.TrieDB().DiskDB()
	it := t.NodeIterator(nil)
	for it.Next(true) {
		if !it.Leaf() {
			continue
		}
		if len(dirtyKeys) > 0 {
			_, ok := dirtyKeys[hexutils.BytesToHex(it.LeafKey())]
			if ok {
				disk.Delete(it.Parent().Bytes())
			}
		}
	}
	fmt.Println("delRoot key : " + hexutils.BytesToHex(root.Bytes()))
	disk.Delete(root.Bytes())
	return nil
}

func (s stateStore) PruningTrie2(dirtyKeys [][]byte) error {
	disk := s.db.TrieDB().DiskDB()
	for _, dk := range dirtyKeys {
		_, err := disk.Get(dk)
		if err != nil {
			fmt.Println(err.Error())
		}
		e := disk.Delete(dk)
		if e != nil {
			return e
		}
	}
	return nil
}

func (s stateStore) getDirtyKeys() [][]byte {
	return nil
}

func (s stateStore) CommitPruningRoot(height uint64, addr, root []byte) {
	store := InstanceOfStateStore().GetDb().TrieDB().DiskDB()
	prefix := PruningRootPrefix(height)
	store.Put(append(prefix, addr...), root)
}
