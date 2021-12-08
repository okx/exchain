package redis_cgi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

const TTL = 1500 * time.Second

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
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
	if deltas == nil || deltas.Size() == 0 {
		return fmt.Errorf("delta is empty")
	}
	deltaBytes, err := deltas.Marshal()
	if err != nil {
		return err
	}
	return r.rdb.SetNX(context.Background(), setDeltaKey(deltas.Height), deltaBytes, TTL).Err()
}

func (r *RedisClient) GetBlock(height int64) (*tmtypes.Block, error) {
	blockBytes, err := r.rdb.Get(context.Background(), setBlockKey(height)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
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
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	deltas := &tmtypes.Deltas{}
	if err = deltas.Unmarshal(deltaBytes); err != nil {
		return nil, err
	}
	return deltas, nil
}

func setBlockKey(height int64) string {
	return "BH:" + strconv.Itoa(int(height))
}

func setDeltaKey(height int64) string {
	return "DH:" + strconv.Itoa(int(height))
}
