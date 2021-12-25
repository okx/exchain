package redis_cgi

import (
	"context"
	"fmt"
	redisgo "github.com/garyburd/redigo/redis"
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

var uploadAndResetHeight = redisgo.NewScript(1, `
	if redis.call("set", KEYS[1], ARGV[1])
	then
		return redis.call("set", ARGV[2], ARGV[3])
	else
		return -1
	end
`)

var once sync.Once

func init() {
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
	rdb      *redis.Client
	pool     *redisgo.Pool
	ttl      time.Duration
	logger   log.Logger
	lockerID string // unique identifier of locker
}

func NewRedisClient(url, auth, lockerID string, ttl time.Duration, l log.Logger) (*RedisClient, error) {
	pool, err := NewPool("redis://" + url, auth, l)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: auth, // no password set
		DB:       0,    // use default DB
	})

	return &RedisClient{rdb, pool, ttl, l, lockerID}, nil
}

func (r *RedisClient) GetLocker() bool {
	res, err := r.rdb.SetNX(context.Background(), deltaLockerKey, r.lockerID, lockerExpire).Result()
	if err != nil {
		r.logger.Error("GetLocker err", err)
		return false
	}
	return res
}

func (r *RedisClient) ReleaseLocker() {
	r.ReleaseDistLock(deltaLockerKey, r.lockerID)
}

// return bool: if change the value of latest_height, need to upload
func (r *RedisClient) ResetLatestHeightAfterUpload(height int64, uploadBytes []byte) bool {
	var res bool
	h, err := r.rdb.Get(context.Background(), latestHeightKey).Int64()
	if err != nil && err != redis.Nil {
		return res
	}

	conn := r.pool.Get()
	defer conn.Close()

	if h < height {
		deltaKey := setDeltaKey(height)
		reply, err := uploadAndResetHeight.Do(conn, deltaKey, uploadBytes, latestHeightKey, height)
		r.logger.Debug(fmt.Sprintf("uploadAndResetHeight: trying to set key(%s) with value(%s), and resetLasteHeight %d to %d. reply(%T, %+v)",
			deltaKey, uploadBytes, h, height, reply, reply))
		if err == nil && reply == "OK" {
			r.logger.Info(fmt.Sprintf("uploadAndResetHeight: set key(%s) with valueLen(%d), and resetLasteHeight %d to %d",
				deltaKey, len(uploadBytes), h, height))
			res = true
		} else {
			r.logger.Error("Failed to reset LatestHeightKey", "err", err, "reply", reply)
		}
	} else {
		r.logger.Info("uploadAndResetHeight: latestHeight is bigger, no need to upload")
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
