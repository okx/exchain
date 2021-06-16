package watcher

import (
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/okex/exchain/x/stream/common"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"time"
)

type RedisDB struct {
	db     *redis.Pool
	logger log.Logger
}

func initRedisDB(dbUrl string, dbPassword string) *RedisDB {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	pool, err := common.NewPool(dbUrl, dbPassword, logger)
	if err != nil {
		panic(err)
	}
	return &RedisDB{db: pool, logger: logger}
}

func (db *RedisDB) Set(key []byte, value []byte) {
	conn := db.db.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", hex.EncodeToString(key), hex.EncodeToString(value)))
	if nil != err {
		db.logger.Error(fmt.Sprintf("redis: trying to set key(%s) with value(%s), err(%+v)", hex.EncodeToString(key), hex.EncodeToString(value), err))
	}
}

func (db *RedisDB) Get(key []byte) ([]byte, error) {
	conn := db.db.Get()
	defer conn.Close()

	//todo del
	fromTime := time.Now()
	result, err := redis.String(conn.Do("GET", hex.EncodeToString(key)))
	fmt.Println("RedisDB get spend time ", time.Since(fromTime))
	if nil != err {
		db.logger.Error(fmt.Sprintf("redis: trying to get key(%s) , err(%+v)", hex.EncodeToString(key), err))
		return nil, err
	}
	return hex.DecodeString(result)
}

func (db *RedisDB) Delete(key []byte) {
	conn := db.db.Get()
	defer conn.Close()

	_, err := redis.Bool(conn.Do("DEL", hex.EncodeToString(key)))
	if nil != err {
		db.logger.Error(fmt.Sprintf("redis: trying to del key(%s) , err(%+v)", hex.EncodeToString(key), err))
	}
}

func (db *RedisDB) Has(key []byte) bool {
	conn := db.db.Get()
	defer conn.Close()

	result, err := redis.Bool(conn.Do("EXISTS", hex.EncodeToString(key)))
	if nil != err {
		db.logger.Error(fmt.Sprintf("redis: trying to exits key(%s) , err(%+v)", hex.EncodeToString(key), err))
	}
	return result
}
