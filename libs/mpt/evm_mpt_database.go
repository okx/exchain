package mpt

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/types"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

var (
	gEvmMptDatabase ethstate.Database = nil
	initEvmOnce     sync.Once
)

const (
	EvmDataDir = "data"
	EvmSpace   = "evm"
)

func InstanceOfEvmStore(homeDir string) ethstate.Database {
	initEvmOnce.Do(func() {
		path := filepath.Join(homeDir, EvmDataDir)

		backend := viper.GetString(FlagDBBackend)
		if backend == "" {
			backend = string(types.GoLevelDBBackend)
		}

		kvstore, e := types.CreateKvDB(EvmSpace, types.BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}

		db := rawdb.NewDatabase(kvstore)
		gEvmMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     int(types.TrieCacheSize),
			Journal:   "",
			Preimages: true,
		})
	})

	return gEvmMptDatabase
}
