package baseapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	db "github.com/okex/exchain/libs/tm-db"
	"path/filepath"
	"sync"
)

const (
	FlagDBBackend    = "db_backend"

	HistoryGasUsedDbDir  = "data"
	HistoryGasUsedDBName = "hgu"

	FlagGasUsedFactor = "gu_factor"
)

var (
	once sync.Once
	guDB db.DB
	GasUsedFactor = 0.4
)

func InstanceOfHistoryGasUsedRecordDB() db.DB {
	once.Do(func() {
		guDB = initDb()
	})
	return guDB
}

func initDb() db.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, HistoryGasUsedDbDir)
	backend := viper.GetString(FlagDBBackend)
	if backend == "" {
		backend = string(db.GoLevelDBBackend)
	}

	return db.NewDB(HistoryGasUsedDBName, db.BackendType(backend), dbPath)
}
