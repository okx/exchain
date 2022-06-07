package watcher

import (
	"path/filepath"
	"sync"

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
	enableWatcher   bool
	db              dbm.DB
	once            sync.Once
	txStateCache    []*watcherMessage
	blockStateCache = make(map[string]*watcherMessage)
)

type watcherMessage struct {
	key      []byte
	value    []byte
	isDelete bool
}

func Enable() bool {
	once.Do(initDB)
	return enableWatcher
}

func NewReadStore(pre []byte) sdk.KVStore {
	once.Do(initDB)
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

func initDB() {
	v := viper.Get(watcher.FlagFastQuery)
	if enable, ok := v.(bool); !ok || !enable {
		return
	}
	enableWatcher = true
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, watcher.WatchDbDir)
	backend := viper.GetString(watcher.FlagDBBackend)
	if backend == "" {
		backend = string(dbm.GoLevelDBBackend)
	}
	db = dbm.NewDB(watchDBName, dbm.BackendType(backend), dbPath)
	go taskRoutine()
}

var tasks = make(chan func(), 10)

func taskRoutine() {
	for task := range tasks {
		task()
	}
}
