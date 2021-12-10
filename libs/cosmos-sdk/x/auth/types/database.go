package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
	"path/filepath"
	"strings"
	"sync"
)

var (
	gEvmMptDatabase ethstate.Database = nil

	initOnce sync.Once

	TrieDirtyDisabled = false
	TrieCacheSize uint = 1554 // MB
)

func init() {
	dbCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewLevelDB(name, dir)
	}
	registerDBCreator(GoLevelDBBackend, dbCreator, false)
}

func NewLevelDB(name string, dir string) (ethdb.KeyValueStore, error) {
	file := filepath.Join(dir, name+".db")
	return leveldb.New(file, 128, 1024, EvmSpace, false)
}

const (
	EvmDataDir = "data"
	EvmSpace   = "evm"

	FlagDBBackend    = "db_backend"

	FlagTrieDirtyDisabled = "trie-dirty-disabled"
	FlagTrieCacheSize = "trie-cache-size"
)

func InstanceOfEvmStore() ethstate.Database {
	initOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		path := filepath.Join(homeDir, EvmDataDir)

		backend := viper.GetString(FlagDBBackend)
		if backend == "" {
			backend = string(dbm.GoLevelDBBackend)
		}

		kvstore, e := CreateKvDB(EvmSpace, BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}
		db := rawdb.NewDatabase(kvstore)

		//freezerPath := filepath.Join(homeDir, EvmDataDir, FreezerSpace)
		//frdb, err := rawdb.NewDatabaseWithFreezer(kvdb, freezerPath, EvmSpace, false)
		//if err != nil {
		//	kvdb.Close()
		//	panic(fmt.Sprintf("fail to init evm mpt database: %v", err))
		//}

		gEvmMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     int(TrieCacheSize),
			Journal:   "",
			Preimages: true,
		})
	})

	return gEvmMptDatabase
}


type BackendType string

// These are valid backend types.
const (
	// GoLevelDBBackend represents goleveldb (github.com/syndtr/goleveldb - most
	// popular implementation)
	//   - pure go
	//   - stable
	GoLevelDBBackend BackendType = "goleveldb"

	// RocksDBBackend represents rocksdb (uses github.com/tecbot/gorocksdb)
	//   - EXPERIMENTAL
	//   - requires gcc
	//   - use rocksdb build tag (go build -tags rocksdb)
	RocksDBBackend BackendType = "rocksdb"
)
type dbCreator func(name string, dir string) (ethdb.KeyValueStore, error)

var backends = map[BackendType]dbCreator{}

func registerDBCreator(backend BackendType, creator dbCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

func CreateKvDB(name string, backend BackendType, dir string) (ethdb.KeyValueStore, error) {
	dbCreator, ok := backends[backend]
	if !ok {
		keys := make([]string, len(backends))
		i := 0
		for k := range backends {
			keys[i] = string(k)
			i++
		}
		panic(fmt.Sprintf("Unknown db_backend %s, expected either %s", backend, strings.Join(keys, " or ")))
	}

	return dbCreator(name, dir)
}
