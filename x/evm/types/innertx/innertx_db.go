package innertx

import (
	"errors"
	"fmt"

	ethvm "github.com/ethereum/go-ethereum/core/vm"
	dbm "github.com/okex/exchain/libs/tm-db"
)

func InitDB(innerTxPath, dbBackendStr string) error {
	var creator ethvm.DBCreator
	switch dbm.BackendType(dbBackendStr) {
	case dbm.RocksDBBackend:
		creator = newRocksDBCreator()
	case dbm.GoLevelDBBackend:
		creator = ethvm.DefaultCreator()
	default:

		return errors.New(fmt.Sprintf("Unknown db_backend %s", dbBackendStr))
	}
	return ethvm.InitDB(innerTxPath, creator)
}

func CloseDB() []error {
	return ethvm.CloseDB()
}
