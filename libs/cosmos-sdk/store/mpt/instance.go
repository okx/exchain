package mpt

import (
	"encoding/binary"
	"path/filepath"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"
)

const (
	mptDataDir = "data"
	mptSpace   = "mpt"
)

var (
	gMptDatabase ethstate.Database = nil
	initMptOnce  sync.Once
	gStatic      = NewRuntimeState()
	gAsyncDB     *AsyncKeyValueStore
)

func InstanceOfMptStore() ethstate.Database {
	initMptOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		path := filepath.Join(homeDir, mptDataDir)

		backend := viper.GetString(sdk.FlagDBBackend)
		if backend == "" {
			backend = string(types.GoLevelDBBackend)
		}

		kvstore, e := types.CreateKvDB(mptSpace, types.BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}
		nkvstore := NewStatKeyValueStore(kvstore, gStatic)
		if EnableAsyncCommit && TrieAsyncDB {
			gAsyncDB = NewAsyncKeyValueStore(nkvstore, false)
			nkvstore = gAsyncDB
		}

		db := rawdb.NewDatabase(nkvstore)
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

func (ms *MptStore) HasVersion(height int64) bool {
	return ms.GetMptRootHash(uint64(height)) != ethcmn.Hash{}
}

func HasVersionByDiskDB(height int64) bool {
	hhash := sdk.Uint64ToBigEndian(uint64(height))
	rst, err := InstanceOfMptStore().TrieDB().DiskDB().Get(append(KeyPrefixAccRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return false
	}
	return true
}
