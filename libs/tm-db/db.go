package db

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BackendType string

// These are valid backend types.
const (
	// GoLevelDBBackend represents goleveldb (github.com/syndtr/goleveldb - most
	// popular implementation)
	//   - pure go
	//   - stable
	GoLevelDBBackend BackendType = "goleveldb"
	// CLevelDBBackend represents cleveldb (uses levigo wrapper)
	//   - fast
	//   - requires gcc
	//   - use cleveldb build tag (go build -tags cleveldb)
	CLevelDBBackend BackendType = "cleveldb"
	// MemDBBackend represents in-memory key value store, which is mostly used
	// for testing.
	MemDBBackend BackendType = "memdb"
	// BoltDBBackend represents bolt (uses etcd's fork of bolt -
	// github.com/etcd-io/bbolt)
	//   - EXPERIMENTAL
	//   - may be faster is some use-cases (random reads - indexer)
	//   - use boltdb build tag (go build -tags boltdb)
	BoltDBBackend BackendType = "boltdb"
	// RocksDBBackend represents rocksdb (uses github.com/tecbot/gorocksdb)
	//   - EXPERIMENTAL
	//   - requires gcc
	//   - use rocksdb build tag (go build -tags rocksdb)
	RocksDBBackend BackendType = "rocksdb"
	// UnknownDBBackend unknown db type
	UnknownDBBackend BackendType = "unknown"

	FlagRocksdbOpts = "rocksdb.opts"
)

type dbCreator func(name string, dir string) (DB, error)

var backends = map[BackendType]dbCreator{}

func registerDBCreator(backend BackendType, creator dbCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

// NewDB creates a new database of type backend with the given name.
// NOTE: function panics if:
//   - backend is unknown (not registered)
//   - creator function, provided during registration, returns error
func NewDB(name string, backend BackendType, dir string) DB {
	dataType := checkDBType(name, dir)
	if dataType != UnknownDBBackend && dataType != backend {
		panic(fmt.Sprintf("Invalid db_backend for <%s> ; expected %s, got %s",
			filepath.Join(dir, name+".db"),
			dataType,
			backend))
	}

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

	db, err := dbCreator(name, dir)
	if err != nil {
		panic(fmt.Sprintf("Error initializing DB: %v", err))
	}
	return db
}

// checkDBType check whether the db file is goleveldb or rocksdb,
// only goleveldb and rocksdb are supported, otherwise it returns unknown.
// Ignore artificial changes to db files
func checkDBType(name string, dir string) BackendType {
	logPath := filepath.Join(dir, name+".db", "LOG")
	file, err := os.Open(logPath)
	if err != nil {
		return UnknownDBBackend
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var firstLine, secondLine string
	line := 0
	for scanner.Scan() {
		line++
		if line == 1 {
			firstLine = scanner.Text()
		}
		if line == 2 {
			secondLine = scanner.Text()
			break
		}
	}

	if strings.Contains(firstLine, "RocksDB") {
		return RocksDBBackend
	}
	if strings.Contains(secondLine, "Level") {
		return GoLevelDBBackend
	}

	return UnknownDBBackend
}
