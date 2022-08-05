package watcher

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	"github.com/okex/exchain/libs/cosmos-sdk/store/gaskv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	stypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

const (
	watchDBName = "wasm-watcher"
)

var (
	checked       bool
	enableWatcher bool
	db            dbm.DB
	// used for parallel deliver txs mode
	txCacheMtx      sync.Mutex
	txStateCache    []*watcherMessage
	blockStateCache = make(map[string]*watcherMessage)
)

type watcherMessage struct {
	key      []byte
	value    []byte
	isDelete bool
}

var checkOnce sync.Once

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

var tasks = make(chan func(), 5*3)

func taskRoutine() {
	for task := range tasks {
		task()
	}
}
