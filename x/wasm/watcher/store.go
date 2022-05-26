package watcher

import (
	"path/filepath"
	"sync"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

const (
	watchDBName = "wasm-watcher"
)

var (
	enableWatcher bool
	db            dbm.DB
	once          sync.Once
)

var (
	QueryTypeKey     = "query-type"
	QueryWatchDBOnly = "watchDB-only"

	notFoundValue  = []byte("Not Found,.;'[-%&^${") // some random bytes
	wasmStateCache = make(map[string][]byte)
)

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
}
