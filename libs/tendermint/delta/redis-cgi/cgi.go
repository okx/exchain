package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"strconv"
)

const TTL = 0

type RedisClient struct {
	rdb *redis.Client
	logger log.Logger
}

func NewRedisClient(url, auth string, l log.Logger) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: auth, // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{rdb, l}
}

func (r *RedisClient) SetBlock(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("block is empty")
	}
	req := r.rdb.SetNX(context.Background(), setBlockKey(height), bytes, TTL)
	return req.Err()
}

func (r *RedisClient) SetDeltas(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("delta is empty")
	}
	req := r.rdb.SetNX(context.Background(), setDeltaKey(height), bytes, TTL)
	return req.Err()
}

func (r *RedisClient) GetBlock(height int64) ([]byte, error) {
	bytes, err := r.rdb.Get(context.Background(), setBlockKey(height)).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("get empty block")
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (r *RedisClient) GetDeltas(height int64) ([]byte, error) {
	bytes, err := r.rdb.Get(context.Background(), setDeltaKey(height)).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("get empty delta")
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func setBlockKey(height int64) string {
	return "BH:" + strconv.Itoa(int(height))
}

func setDeltaKey(height int64) string {
	return "DH:" + strconv.Itoa(int(height))
}
