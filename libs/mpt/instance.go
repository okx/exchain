package mpt

import (
	"encoding/binary"
	"path/filepath"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/mpt/types"
	"github.com/spf13/viper"
)

const (
	mptDataDir = "data"
	mptSpace   = "mpt"
)

var (
	gMptDatabase ethstate.Database = nil
	initMptOnce  sync.Once
)

func InstanceOfMptStore() ethstate.Database {
	initMptOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		path := filepath.Join(homeDir, mptDataDir)

		backend := sdk.DBBackend
		if backend == "" {
			backend = string(types.GoLevelDBBackend)
		}

		kvstore, e := types.CreateKvDB(mptSpace, types.BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}
		db := rawdb.NewDatabase(kvstore)
		gMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     int(TrieCacheSize),
			Journal:   "",
			Preimages: true,
		})
	})

	return gMptDatabase
}

// GetLatestStoredBlockHeight get latest mpt storage height
func (ms *MptStore) GetLatestStoredBlockHeight() uint64 {
	rst, err := ms.db.TrieDB().DiskDB().Get(KeyPrefixAccLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (ms *MptStore) SetLatestStoredBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	ms.db.TrieDB().DiskDB().Put(KeyPrefixAccLatestStoredHeight, hhash)
}

// GetMptRootHash gets root mpt hash from block height
func (ms *MptStore) GetMptRootHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := ms.db.TrieDB().DiskDB().Get(append(KeyPrefixAccRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func (ms *MptStore) SetMptRootHash(height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	ms.db.TrieDB().DiskDB().Put(append(KeyPrefixAccRootMptHash, hhash...), hash.Bytes())
}
