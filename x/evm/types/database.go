package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

var (
	gEvmMptDatabase ethstate.Database = nil

	initOnce sync.Once
)

const (
	EvmDataDir = "data"
	EvmSpace   = "evm"
	//FreezerSpace = "freezer"
)

type Watcher interface {
	SaveAccount(account auth.Account, isDirectly bool)
	SaveState(addr ethcmn.Address, key, value []byte)
	Enabled() bool
	SaveContractBlockedListItem(addr sdk.AccAddress)
	SaveContractDeploymentWhitelistItem(addr sdk.AccAddress)
	DeleteContractBlockedList(addr sdk.AccAddress)
	DeleteContractDeploymentWhitelist(addr sdk.AccAddress)
}

type DefaultPrefixDb struct {
}

func (d DefaultPrefixDb) NewStore(parent types.KVStore, Prefix []byte) StoreProxy {
	return prefix.NewStore(parent, Prefix)
}

type StoreProxy interface {
	Set(key, value []byte)
	Get(key []byte) []byte
	Delete(key []byte)
	Has(key []byte) bool
}

type DbAdapter interface {
	NewStore(parent types.KVStore, prefix []byte) StoreProxy
}

func InstanceOfEvmStore() ethstate.Database {
	initOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		file := filepath.Join(homeDir, EvmDataDir, EvmSpace+".db")
		//freezerPath := filepath.Join(homeDir, EvmDataDir, FreezerSpace)

		kvdb, err := leveldb.New(file, 128, 1024, EvmSpace, false)
		if err != nil {
			panic(fmt.Sprintf("fail to open level database: %v", err))
		}

		db := rawdb.NewDatabase(kvdb)
		//frdb, err := rawdb.NewDatabaseWithFreezer(kvdb, freezerPath, EvmSpace, false)
		//if err != nil {
		//	kvdb.Close()
		//	panic(fmt.Sprintf("fail to init evm mpt database: %v", err))
		//}

		gEvmMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     256,
			Journal:   "",
			Preimages: true,
		})
	})

	return gEvmMptDatabase
}
