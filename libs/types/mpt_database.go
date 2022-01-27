package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"path/filepath"
	"strings"
)

const (
	FlagTrieDirtyDisabled = "trie-dirty-disabled"
	FlagTrieCacheSize     = "trie-cache-size"
	FlagEnableDoubleWrite = "enable-double-write"
)

var (
	TrieDirtyDisabled      = false
	TrieCacheSize     uint = 2048 // MB
	EnableDoubleWrite      = false
)

//------------------------------------------
type (
	BackendType string

	dbCreator func(name string, dir string) (ethdb.KeyValueStore, error)
)

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

	// MemDBBackend represents in-memory key value store, which is mostly used
	// for testing.
	MemDBBackend BackendType = "memdb"
)

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

//------------------------------------------
//	Register go-ethereum memdb and leveldb
//------------------------------------------
func init() {
	levelDBCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewMptLevelDB(name, dir)
	}

	memDBCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewMptMemDB(name, dir)
	}

	registerDBCreator(GoLevelDBBackend, levelDBCreator, false)
	registerDBCreator(MemDBBackend, memDBCreator, false)
}

func NewMptLevelDB(name string, dir string) (ethdb.KeyValueStore, error) {
	file := filepath.Join(dir, name+".db")
	return leveldb.New(file, 128, 1024, name, false)
}

func NewMptMemDB(name string, dir string) (ethdb.KeyValueStore, error) {
	return memorydb.New(), nil
}