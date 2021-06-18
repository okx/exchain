package watcher

import (
	"strings"
)

type RedisAndHbaseDB struct {
	rdb OperateDB
	hdb OperateDB
}

func initRedisAndHbaseDB(dbUrl string, dbPassword string) *RedisAndHbaseDB {
	dbUrls := strings.Split(dbUrl, ";")
	if len(dbUrls) != 2 {
		panic("wrong input , eg : redisUrl;hbaseUrl")
	}
	dbPasswords := strings.Split(dbPassword, ";")

	rdb := initRedisDB(dbUrls[0], dbPasswords[0], "172800") // two days
	hdb := initHbaseDB(dbUrls[1])
	return &RedisAndHbaseDB{rdb: rdb, hdb: hdb}
}

func (db *RedisAndHbaseDB) Set(key []byte, value []byte) {
	db.rdb.Set(key, value)
	db.hdb.Set(key, value)
}

func (db *RedisAndHbaseDB) Get(key []byte) ([]byte, error) {
	if result, err := db.rdb.Get(key); err == nil {
		return result, err
	}
	return db.hdb.Get(key)
}

func (db *RedisAndHbaseDB) Delete(key []byte) {
	db.rdb.Delete(key)
	db.hdb.Delete(key)
}

func (db *RedisAndHbaseDB) Has(key []byte) bool {
	return db.rdb.Has(key) || db.hdb.Has(key)
}
