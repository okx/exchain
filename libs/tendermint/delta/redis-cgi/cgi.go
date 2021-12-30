package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
	"time"
)

const (
	lockerExpire = 4 * time.Second
)

var (
	latestHeightKey string
	deltaLockerKey  string
)

var once sync.Once
func init()  {
	const (
		latestHeight = "LatestHeight"
		deltaLocker  = "DeltaLocker"
	)
	once.Do(func() {
		latestHeightKey = fmt.Sprintf("dds:%d:%s", types.DeltaVersion, latestHeight)
		deltaLockerKey = fmt.Sprintf("dds:%d:%s", types.DeltaVersion, deltaLocker)
	})
}

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
	res, err := r.rdb.SetNX(context.Background(), deltaLockerKey, true, lockerExpire).Result()
	if err != nil {
		r.logger.Error("GetLocker err", err)
		return false
	}
	return res
}

func (r *RedisClient) ReleaseLocker() {
	_, err := r.rdb.Del(context.Background(), deltaLockerKey).Result()
	if err != nil {
		r.logger.Error("Failed to Release Locker", "err", err)
	}
}

// return bool: if change the value of latest_height, need to upload
func (r *RedisClient) ResetLatestHeightAfterUpload(height int64, upload func() bool) bool {
	var res bool
	h, err := r.rdb.Get(context.Background(), latestHeightKey).Int64()
	if err != nil && err != redis.Nil {
		return res
	}

	if h < height && upload() {
		err = r.rdb.Set(context.Background(), latestHeightKey, height, 0).Err()
		if err == nil {
			r.logger.Info("Reset LatestHeightKey", "height", height)
			res = true
		} else {
			r.logger.Error("Failed to reset LatestHeightKey","err", err)
		}
	}
	return res
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
	req := r.rdb.SetNX(context.Background(), genDeltaKey(height), bytes, r.ttl)
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

func (r *RedisClient) GetDeltas(height int64) ([]byte, error, int64) {
	latestHeight := r.getLatestHeight()
	bytes, err := r.rdb.Get(context.Background(), genDeltaKey(height)).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("get empty delta"), latestHeight
	}
	return bytes, err, latestHeight
}

func (r *RedisClient) getLatestHeight() (latestHeight int64) {
	latestHeight = -1
	h, err := r.rdb.Get(context.Background(), latestHeightKey).Int64()
	if err == nil {
		latestHeight = h
	} else if err == redis.Nil {
		latestHeight = 0
	}
	return
}

func setBlockKey(height int64) string {
	return fmt.Sprintf("BH:%d", height)
}

func genDeltaKey(height int64) string {
	return fmt.Sprintf("DH-%d:%d", types.DeltaVersion, height)
}
