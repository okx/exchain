package watcher

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"path/filepath"
	"time"
)

type LevelDB struct {
	db *leveldb.DB
}

func initLevelDB(homeDir string) *LevelDB {
	dbPath := filepath.Join(homeDir, "data/watch.db")
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDB{db: db}
}

func (db *LevelDB) Set(key []byte, value []byte) {
	err := db.db.Put(key, value, nil)
	if err != nil {
		log.Println("LevelDB error: ", err.Error())
	}
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	//todo del
	fromTime := time.Now()
	result, err := db.db.Get(key, nil)
	fmt.Println("LevelDB get spend time ", time.Since(fromTime))
	return result, err
}

func (db *LevelDB) Delete(key []byte) {
	err := db.db.Delete(key, nil)
	if err != nil {
		log.Printf("LevelDB error: " + err.Error())
	}
}

func (db *LevelDB) Has(key []byte) bool {
	res, err := db.db.Has(key, nil)
	if err != nil {
		log.Println("LevelDB error: " + err.Error())
		return false
	}
	return res
}
