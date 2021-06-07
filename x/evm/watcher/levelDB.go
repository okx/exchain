package watcher

import (
	"encoding/hex"
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
	start := time.Now()
	result, err := db.db.Get(key, nil)
	//todo del
	log.Println(fmt.Sprintf("levelDB get key(%s) , value(%s), time(%s)", hex.EncodeToString(key), hex.EncodeToString(result), time.Since(start)))
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
