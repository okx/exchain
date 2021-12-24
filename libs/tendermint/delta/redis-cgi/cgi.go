package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"time"
)

const (
	LatestHeightKey = "LatestHeight"
	DeltaLockerKey  = "DeltaLocker"
	LockerExpire    = 4 * time.Second
)

type RedisClient struct {
	rdb    *redis.Client
	ttl    time.Duration
	logger log.Logger
}

func NewRedisClient(url, auth string, ttl time.Duration, l log.Logger) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: auth, // no password set
		DB:       0,    // use default DB
	})
	return &RedisClient{rdb, ttl, l}
}

func (r *RedisClient) GetLocker() bool {
	res, err := r.rdb.SetNX(context.Background(), DeltaLockerKey, true, LockerExpire).Result()
	if err != nil {
		r.logger.Error("GetLocker err", err)
		return false
	}
	return res
}

func (r *RedisClient) ReleaseLocker() {
	_, err := r.rdb.Del(context.Background(), DeltaLockerKey).Result()
	if err != nil {
		r.logger.Error("ReleaseLocker err", err)
	}
}

// return bool: if change the value of latest_height, need to upload
func (r *RedisClient) SetLatestHeight(height int64) bool {
	h, err := r.rdb.Get(context.Background(), LatestHeightKey).Int64()
	if err != nil && err != redis.Nil {
		return false
	}
	// h is not exist(h==0) or h < height
	// set h and need to upload
	if h < height {
		if r.rdb.Set(context.Background(), LatestHeightKey, height, 0).Err() != nil {
			return false
		}
		return true
	}
	// h is exist and h > height, no need to upload
	return false
}

func (r *RedisClient) SetBlock(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("block is empty")
	}
	req := r.rdb.SetNX(context.Background(), setBlockKey(height), bytes, r.ttl)
	return req.Err()
}

func (r *RedisClient) SetDeltas(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("delta is empty")
	}
	req := r.rdb.SetNX(context.Background(), setDeltaKey(height), bytes, r.ttl)
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
	return fmt.Sprintf("BH:%d", height)
}

func setDeltaKey(height int64) string {
	return fmt.Sprintf("DH-%d:%d", types.DeltaVersion, height)
}
