package redis_cgi

import (
	"fmt"
	redisgo "github.com/garyburd/redigo/redis"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
	"time"
)

const (
	// unit: Millisecond
	lockerExpire = 4000
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

	return &RedisClient{pool, ttl, l, lockerID}, nil
}

func (r *RedisClient) GetLocker() bool {
	res, err := r.FetchDistLock(deltaLockerKey, r.lockerID, lockerExpire)
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
func (r *RedisClient) ResetLatestHeightAfterUpload(height int64, getBytes func() ([]byte, error)) bool {
	var res bool
	conn := r.pool.Get()
	defer conn.Close()
	// get latestHeight
	h, err := redisgo.Int64(conn.Do("GET", latestHeightKey))
	if err != nil && err != redisgo.ErrNil {
		return res
	}

	// get upload bytes
	bytes, err := getBytes()
	if h < height && err == nil {
		// upload and set latestHeight

		deltaKey := calcDeltaKey(height)
		reply, err := uploadAndResetHeight.Do(conn, deltaKey, bytes, latestHeightKey, height)
		if err == nil && reply == "OK" {
			r.logger.Info(fmt.Sprintf("uploadAndResetHeight: set key(%s) with valueLen(%d), and resetLasteHeight %d to %d",
				deltaKey, len(bytes), h, height))
			res = true
		} else {
			r.logger.Error("Failed to reset LatestHeightKey", "err", err, "reply", reply)
		}
	}

	return res
}

func (r *RedisClient) GetDeltas(height int64) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := redisgo.Bytes(conn.Do("GET", calcDeltaKey(height)))
	if err == redisgo.ErrNil {
		return nil, fmt.Errorf("get empty delta")
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func calcDeltaKey(height int64) string {
	return fmt.Sprintf("DH-%d:%d", types.DeltaVersion, height)
}
