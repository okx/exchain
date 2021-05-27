package watcher

import (
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
	db.db.Set(string(key), string(value), 0).Result()
}

func (db *RedisDB) Get(key []byte) ([]byte, error) {
	result, err := db.db.Get(string(key)).Result()
	return []byte(result), err
}

func (db *RedisDB) Delete(key []byte) {
	db.db.Del(string(key)).Result()
}

func (db *RedisDB) Has(key []byte) bool {
	result, _ := db.db.Exists(string(key)).Result()
	return result > 0
}
