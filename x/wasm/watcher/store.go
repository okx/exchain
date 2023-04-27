package watcher

import (
	"encoding/json"
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"io"
	"log"
	"path/filepath"
	"sync"

	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
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
	txCacheMtx         sync.Mutex
	txStateCache       []*WatchMessage
	blockStateCache    = make(map[string]*WatchMessage)
	watchdbForSimulate = dbadapter.Store{}
	accountKeyPrefix   = []byte("wasm-account-")
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
	watchdbForSimulate = dbadapter.Store{DB: db}
	go taskRoutine()
}

func AccountKey(addr []byte) []byte {
	return append(accountKeyPrefix, addr...)
}
func GetAccount(addr sdk.WasmAddress) (*types.EthAccount, error) {
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

func DeleteAccount(addr sdk.WasmAddress) {
	if !Enable() {
		return
	}
	if err := db.Delete(AccountKey(addr.Bytes())); err != nil {
		log.Println("wasm watchDB delete account error", addr.String())
	}
}

func NewReadStore(pre []byte, store sdk.KVStore) sdk.KVStore {
	rs := &readStore{
		mp: make(map[string][]byte, 0),
		kv: store,
	}
	if len(pre) != 0 {
		return prefix.NewStore(rs, pre)
	}
	return rs
}

type Adapter struct{}

func (a Adapter) NewStore(ctx sdk.Context, storeKey sdk.StoreKey, pre []byte) sdk.KVStore {
	if ctx.WasmKvStoreForSimulate() != nil {
		return ctx.WasmKvStoreForSimulate()
	}
	s := NewReadStore(pre, ctx.KVStore(storeKey))
	ctx.SetWasmKvStoreForSimulate(s)
	return s
}

type readStore struct {
	mp map[string][]byte
	kv sdk.KVStore
}

func (r *readStore) GetStoreType() cosmost.StoreType {
	return r.kv.GetStoreType()
}

func (r *readStore) CacheWrap() cosmost.CacheWrap {
	return r.kv.CacheWrap()
}

func (r *readStore) CacheWrapWithTrace(w io.Writer, tc cosmost.TraceContext) cosmost.CacheWrap {
	return r.kv.CacheWrapWithTrace(w, tc)
}

func (r *readStore) Get(key []byte) []byte {
	if value, ok := r.mp[string(key)]; ok {
		return value
	}
	if value := watchdbForSimulate.Get(key); len(value) != 0 {
		return value
	}
	return r.kv.Get(key)
}

func (r *readStore) Has(key []byte) bool {
	if _, ok := r.mp[string(key)]; ok {
		return ok
	}
	return r.kv.Has(key)
}

func (r *readStore) Set(key, value []byte) {
	r.mp[string(key)] = value
}

func (r readStore) Delete(key []byte) {
	delete(r.mp, string(key))
}

func (r readStore) Iterator(start, end []byte) cosmost.Iterator {
	return r.kv.Iterator(start, end)
}

func (r readStore) ReverseIterator(start, end []byte) cosmost.Iterator {
	return r.kv.ReverseIterator(start, end)
}
