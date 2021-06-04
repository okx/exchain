package watcher

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"path/filepath"
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
	result, err := db.db.Get(key, nil)
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
