package mpt

import (
	"encoding/binary"
	"github.com/okex/exchain/libs/types"
	"path/filepath"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"
)

const (
	AccMptDataDir = "data"
	AccMptSpace   = "acc"
)

var (
	gAccMptDatabase ethstate.Database = nil
	initAccOnce     sync.Once
)

func InstanceOfMptStore() ethstate.Database {
	initAccOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		path := filepath.Join(homeDir, AccMptDataDir)

		backend := viper.GetString(FlagDBBackend)
		if backend == "" {
			backend = string(types.GoLevelDBBackend)
		}

		kvstore, e := types.CreateKvDB(AccMptSpace, types.BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}
		db := rawdb.NewDatabase(kvstore)
		gAccMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     int(types.TrieCacheSize),
			Journal:   "",
			Preimages: true,
		})
	})

	return gAccMptDatabase
}

var (
	KeyPrefixLatestStoredHeight = []byte{0x01}
	KeyPrefixRootMptHash        = []byte{0x02}
)

// GetLatestStoredBlockHeight get latest mpt storage height
func (ms *MptStore) GetLatestStoredBlockHeight() uint64 {
	rst, err := ms.db.TrieDB().DiskDB().Get(KeyPrefixLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (ms *MptStore) SetLatestStoredBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	ms.db.TrieDB().DiskDB().Put(KeyPrefixLatestStoredHeight, hhash)
}

// GetMptRootHash gets root mpt hash from block height
func (ms *MptStore) GetMptRootHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := ms.db.TrieDB().DiskDB().Get(append(KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func (ms *MptStore) SetMptRootHash(height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	ms.db.TrieDB().DiskDB().Put(append(KeyPrefixRootMptHash, hhash...), hash.Bytes())
}
