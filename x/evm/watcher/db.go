package watcher

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	dbm "github.com/okex/exchain/libs/tm-db"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
)

const (
	FlagFastQuery    = "fast-query"
	FlagFastQueryLru = "fast-lru"
	FlagDBBackend    = "db_backend"
	FlagCheckWd      = "check_watchdb"

	WatchDbDir  = "data"
	WatchDBName = "watch"
)

type WatchStore struct {
	db          dbm.DB
	params      evmtypes.Params
	paramsMutex sync.RWMutex
}

var gWatchStore *WatchStore = nil
var once sync.Once

func InstanceOfWatchStore() *WatchStore {
	once.Do(func() {
		if IsWatcherEnabled() {
			gWatchStore = &WatchStore{db: initDb(), params: evmtypes.DefaultParams()}
		}
	})
	return gWatchStore
}

func initDb() dbm.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, WatchDbDir)
	backend := viper.GetString(FlagDBBackend)
	if backend == "" {
		backend = string(dbm.GoLevelDBBackend)
	}

	versionPath := filepath.Join(dbPath, WatchDBName+".db", "VERSION")
	if !checkVersion(versionPath) {
		os.RemoveAll(filepath.Join(dbPath, WatchDBName+".db"))
	}

	db := dbm.NewDB(WatchDBName, dbm.BackendType(backend), dbPath)
	writeVersion(versionPath)
	return db
}

func checkVersion(versionPath string) bool {
	_, err := os.Stat(versionPath)
	if os.IsNotExist(err) {
		writeVersion(versionPath)
		return true
	}
	content, err := ioutil.ReadFile(versionPath)
	if err != nil {
		panic(err)
	}
	if string(content) == version {
		return true
	}
	return false
}

func writeVersion(versionPath string) {
	ioutil.WriteFile(versionPath, []byte(version), 0666)
}

func (w WatchStore) Set(key []byte, value []byte) {
	err := w.db.Set(key, value)
	if err != nil {
		log.Println("watchdb error: ", err.Error())
	}
}

func (w WatchStore) Get(key []byte) ([]byte, error) {
	return w.db.Get(key)
}

func (w WatchStore) GetUnsafe(key []byte, processor dbm.UnsafeValueProcessor) (interface{}, error) {
	return w.db.GetUnsafeValue(key, processor)
}

func (w WatchStore) Delete(key []byte) {
	err := w.db.Delete(key)
	if err != nil {
		log.Printf("watchdb error: " + err.Error())
	}
}

func (w WatchStore) Has(key []byte) bool {
	res, err := w.db.Has(key)
	if err != nil {
		log.Println("watchdb error: " + err.Error())
		return false
	}
	return res
}

func (w WatchStore) Iterator(start, end []byte) dbm.Iterator {
	it, err := w.db.Iterator(start, end)
	if err != nil {
		log.Println("watchdb error: " + err.Error())
		return nil
	}
	return it
}

func (w WatchStore) GetEvmParams() evmtypes.Params {
	w.paramsMutex.RLock()
	defer w.paramsMutex.RUnlock()
	return w.params
}

func (w *WatchStore) SetEvmParams(params evmtypes.Params) {
	w.paramsMutex.Lock()
	defer w.paramsMutex.Unlock()
	w.params = params
}
