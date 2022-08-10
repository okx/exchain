package watcher

import (
	"encoding/json"
	"log"
	"path/filepath"
	"sync"

	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	"github.com/okex/exchain/libs/cosmos-sdk/store/gaskv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	stypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

const (
	watchDBName = "wasm-watcher"
)

var (
	checkOnce     sync.Once
	checked       bool
	enableWatcher bool
	db            dbm.DB
	// used for parallel deliver txs mode
	txCacheMtx      sync.Mutex
	txStateCache    []*WatchMessage
	blockStateCache = make(map[string]*WatchMessage)

	accountKeyPrefix = []byte("wasm-account-")
)

func CheckEnable() bool {
	checkOnce.Do(func() {
		checked = true
		if viper.GetBool(watcher.FlagFastQuery) {
			enableWatcher = true
		}
	})
	return enableWatcher
}

func Enable() bool {
	if !checked {
		panic("fast query should be checked at init")
	}
	return enableWatcher
}

func InitDB() {
	if !Enable() {
		return
	}
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, watcher.WatchDbDir)
	backend := viper.GetString(watcher.FlagDBBackend)
	if backend == "" {
		backend = string(dbm.GoLevelDBBackend)
	}
	db = dbm.NewDB(watchDBName, dbm.BackendType(backend), dbPath)
	go taskRoutine()
}

func AccountKey(addr []byte) []byte {
	return append(accountKeyPrefix, addr...)
}
func GetAccount(addr sdk.AccAddress) (*types.EthAccount, error) {
	b, err := db.Get(AccountKey(addr.Bytes()))
	if err != nil {
		return nil, err
	}

	var acc types.EthAccount
	err = json.Unmarshal(b, &acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil

}

func SetAccount(acc *types.EthAccount) error {
	b, err := json.Marshal(acc)
	if err != nil {
		return err
	}
	return db.Set(AccountKey(acc.Address.Bytes()), b)
}

func DeleteAccount(addr sdk.AccAddress) {
	if err := db.Delete(AccountKey(addr.Bytes())); err != nil {
		log.Println("wasm watchDB delete account error", addr.String())
	}
}

func NewReadStore(pre []byte) sdk.KVStore {
	rs := &readStore{
		Store: dbadapter.Store{DB: db},
	}
	if len(pre) != 0 {
		return prefix.NewStore(rs, pre)
	}
	return rs
}

type Adapter struct{}

func (a Adapter) NewStore(gasMeter sdk.GasMeter, _ sdk.KVStore, pre []byte) sdk.KVStore {
	store := NewReadStore(pre)
	return gaskv.NewStore(store, gasMeter, stypes.KVGasConfig())
}

type readStore struct {
	dbadapter.Store
}

func (r *readStore) Set(key, value []byte) {}
func (r *readStore) Delete(key []byte)     {}
