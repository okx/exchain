package watcher

import (
	"encoding/json"
	"log"
	"path/filepath"
	"sync"

	"github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/dbadapter"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/gaskv"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/prefix"
	stypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	dbm "github.com/okx/okbchain/libs/tm-db"
	"github.com/okx/okbchain/x/evm/watcher"
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

func Enable() bool {
	checkOnce.Do(func() {
		checked = true
		if viper.GetBool(watcher.FlagFastQuery) {
			enableWatcher = true
			InitDB()
		}
	})
	return enableWatcher
}

func ensureChecked() {
	if !checked {
		panic("fast query should be checked at init")
	}
}

func InitDB() {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, watcher.WatchDbDir)

	var err error
	db, err = sdk.NewDB(watchDBName, dbPath)
	if err != nil {
		panic(err)
	}
	go taskRoutine()
}

func AccountKey(addr []byte) []byte {
	return append(accountKeyPrefix, addr...)
}
func GetAccount(addr sdk.AccAddress) (*types.EthAccount, error) {
	if !Enable() {
		return nil, nil
	}
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
	if !Enable() {
		return nil
	}
	b, err := json.Marshal(acc)
	if err != nil {
		return err
	}
	return db.Set(AccountKey(acc.Address.Bytes()), b)
}

func DeleteAccount(addr sdk.AccAddress) {
	if !Enable() {
		return
	}
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
