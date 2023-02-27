package snapshot

import "github.com/ethereum/go-ethereum/ethdb"

var gConfigure configure

type configure struct {
	diskDB ethdb.Database
}

func SetDiskDB(db ethdb.Database) {
	gConfigure.diskDB = db
}

func GetDiskDB() ethdb.Database {
	return gConfigure.diskDB
}
