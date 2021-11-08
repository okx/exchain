package baseapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	db "github.com/tendermint/tm-db"
	"path/filepath"
	"sync"
)

const (
	FlagDBBackend    = "db_backend"

	WatchDbDir  = "data"
	WatchDBName = "gas"
)

var once sync.Once
var guDB db.DB

func InstanceOfGasUsedRecordDB() db.DB {
	once.Do(func() {
		guDB = initDb()
	})
	return guDB
}

func initDb() db.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, WatchDbDir)
	backend := viper.GetString(FlagDBBackend)
	if backend == "" {
		backend = string(db.GoLevelDBBackend)
	}

	return db.NewDB(WatchDBName, db.BackendType(backend), dbPath)
}
