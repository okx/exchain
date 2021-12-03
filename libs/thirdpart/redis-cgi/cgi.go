package redis_cgi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
)

const TTL = 1500 * time.Second

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(tmtypes.DataCenterUrl),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{rdb}
}

func (r *RedisClient) SetBlock(block *tmtypes.Block) error {
	if block == nil || block.Size() <= 0 {
		return fmt.Errorf("block is empty")
	}
	blockBytes, err := block.Marshal()
	if err != nil {
		return err
	}
	return r.rdb.SetNX(context.Background(), setBlockKey(block.Height), blockBytes, TTL).Err()
}

func (r *RedisClient) SetDelta(deltas *tmtypes.Deltas) error {
	if deltas.Size() <= 0 {
		return fmt.Errorf("delta is empty")
	}
	deltaBytes, err := deltas.Marshal()
	if err != nil {
		return err
	}
	return r.rdb.SetNX(context.Background(), setDeltaKey(deltas.Height), deltaBytes, TTL).Err()
}

func (r *RedisClient) SetWatch(watch *tmtypes.WatchData) error {
	if watch.Size() <= 0 {
		return fmt.Errorf("watch is empty")
	}
	return r.rdb.SetNX(context.Background(), setWatchKey(watch.Height), watch.WatchDataByte, TTL).Err()
}

func (r *RedisClient) GetBlock(height int64) (*tmtypes.Block, error) {
	blockBytes, err := r.rdb.Get(context.Background(), setBlockKey(height)).Bytes()
	if err != nil {
		return nil, err
	}
	block := &tmtypes.Block{}
	if err = block.Unmarshal(blockBytes); err != nil {
		return nil, err
	}
	return block, nil
}

func (r *RedisClient) GetDeltas(height int64) (*tmtypes.Deltas, error) {
	deltaBytes, err := r.rdb.Get(context.Background(), setDeltaKey(height)).Bytes()
	if err != nil {
		return nil, err
	}
	deltas := &tmtypes.Deltas{}
	if err = deltas.Unmarshal(deltaBytes); err != nil {
		return nil, err
	}
	return deltas, nil
}

func (r *RedisClient) GetWatch(height int64) (*tmtypes.WatchData, error) {
	watchBytes, err := r.rdb.Get(context.Background(), setWatchKey(height)).Bytes()
	if err != nil {
		return nil, err
	}
	return &tmtypes.WatchData{
		WatchDataByte: watchBytes,
		Height:        height,
	}, nil
}

func setBlockKey(height int64) string {
	return "BH:" + strconv.Itoa(int(height))
}

func setDeltaKey(height int64) string {
	return "DH:" + strconv.Itoa(int(height))
}

func setWatchKey(height int64) string {
	return "WH:" + strconv.Itoa(int(height))
}
