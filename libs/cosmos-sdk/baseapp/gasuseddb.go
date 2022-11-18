package baseapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const (
	HistoryGasUsedDbDir  = "data"
	HistoryGasUsedDBName = "hgu"

	FlagGasUsedFactor = "gu_factor"
)

var (
	//once          sync.Once
	guDB          db.DB
	GasUsedFactor = 0.4

	hguPath string
)

func InstanceOfHistoryGasUsedRecordDB() db.DB {
	return guDB
}

func CreateHguDB() {
	guDB = initDb()
}

func DeleteHguDB() {
	err := guDB.Close()
	if err != nil {
		log.Println("Close guDB error:", err)
	}
	if hguPath != "" {
		if err = os.RemoveAll(hguPath); err != nil {
			log.Println("Remove guDB error:", err)
		}
	}
	guDB = nil
}

func RecreateHguDB() {
	DeleteHguDB()
	guDB = initDb()
}

func initDb() db.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, HistoryGasUsedDbDir)

	hguPath = filepath.Join(dbPath, HistoryGasUsedDBName+".db")
	log.Println("hguPath:", hguPath)
	db, err := sdk.NewDB(HistoryGasUsedDBName, dbPath)
	if err != nil {
		panic(err)
	}
	return db
}
