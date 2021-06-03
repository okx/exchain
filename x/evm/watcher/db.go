package watcher

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"sync"
)

const (
	FlagFastQuery                 = "fast-query"
	FlagFastQueryLru              = "fast-lru"
	FlagWatcherDisLockUrl         = "watcher-dis-lock-url"
	FlagWatcherDisLockUrlPassword = "watcher-dis-lock-password"
	FlagWatcherDBType             = "watcher-db-type"
	FlagWatcherDBUrl              = "watcher-db-url"
	FlagWatcherDBPassword         = "watcher-db-password"

	DBTypeLevel = "levelDB"
	DBTypeHbase = "hbaseDB"
	DBTypeRedis = "redisDB"
)

type WatchStore struct {
	OperateDB
}

type OperateDB interface {
	Set(key []byte, value []byte)
	Get(key []byte) ([]byte, error)
	Delete(key []byte)
	Has(key []byte) bool
}

var gWatchStore *WatchStore = nil
var once sync.Once

func InstanceOfWatchStore() *WatchStore {
	once.Do(func() {
		if IsWatcherEnabled() {
			var db OperateDB
			// set db by FlagWatcherDB
			dbType := viper.GetString(FlagWatcherDBType)
			if dbType == DBTypeLevel {
				db = initLevelDB(viper.GetString(flags.FlagHome))
			} else if dbType == DBTypeHbase {
				db = initHbaseDB(viper.GetString(FlagWatcherDBUrl))
			} else if dbType == DBTypeRedis {
				db = initRedisDB(viper.GetString(FlagWatcherDBUrl), viper.GetString(FlagWatcherDBPassword))
			}
			if db == nil {
				panic("db has not been initialized")
			}
			gWatchStore = &WatchStore{db}
		}
	})
	return gWatchStore
}
