package watcher

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"sync"
)

const (
	FlagFastQuery         = "fast-query"
	FlagWatcherDisLockUrl = "watcher-dis-lock-url"
	FlagWatcherDBType     = "watcher-db-type"
	FlagHbaseDBUrl        = "hbase-db-url"

	DBTypeLevel = "levelDB"
	DBTypeHbase = "hbaseDB"
)

type WatchStore struct {
	db OperateDB
}

type OperateDB interface {
	Set(key []byte, value []byte)
	Get(key []byte) ([]byte, error)
	Delete(key []byte)
	Has(key []byte) bool
}

var gWatchStore *WatchStore = nil
var once sync.Once
var db OperateDB

func InstanceOfWatchStore() *WatchStore {
	once.Do(func() {
		if IsWatcherEnabled() {
			// set db by FlagWatcherDB
			dbType := viper.GetString(FlagWatcherDBType)
			if dbType == DBTypeLevel {
				db = initLevelDB(viper.GetString(flags.FlagHome))
			} else if dbType == DBTypeHbase {
				db = initHbaseDB(viper.GetString(FlagHbaseDBUrl))
			}
			if db == nil {
				panic("db has not been initialized")
			}
			gWatchStore = &WatchStore{db: db}
		}
	})
	return gWatchStore
}

func (w *WatchStore) Set(key []byte, value []byte) {
	w.db.Set(key, value)
}

func (w *WatchStore) Get(key []byte) ([]byte, error) {
	return w.db.Get(key)
}

func (w *WatchStore) Delete(key []byte) {
	w.db.Delete(key)
}

func (w *WatchStore) Has(key []byte) bool {
	return w.db.Has(key)
}
