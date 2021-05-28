package watcher

import (
	"encoding/hex"
	"github.com/go-redis/redis"
)

type RedisDB struct {
	db *redis.Client
}

func initRedisDB(dbUrl string, dbPassword string) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     dbUrl,
		Password: dbPassword,
	})

	return &RedisDB{db: client}
}

func (db *RedisDB) Set(key []byte, value []byte) {
	db.db.Set(hex.EncodeToString(key), hex.EncodeToString(value), 0).Result()
}

func (db *RedisDB) Get(key []byte) ([]byte, error) {
	result, err := db.db.Get(hex.EncodeToString(key)).Result()
	if err != nil {
		return nil, err
	}
	bz, err := hex.DecodeString(result)
	return bz, err
}

func (db *RedisDB) Delete(key []byte) {
	db.db.Del(hex.EncodeToString(key)).Result()
}

func (db *RedisDB) Has(key []byte) bool {
	result, _ := db.db.Exists(hex.EncodeToString(key)).Result()
	return result > 0
}
