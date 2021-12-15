package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"strconv"
	"time"
)

const TTL = 0

type RedisClient struct {
	rdb *redis.Client
	compressBroker compress.CompressBroker
	logger log.Logger
}

func NewRedisClient(url string, l log.Logger) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{rdb,
		// todo can config different compress algorithm
		&compress.Flate{}, l}
}

func (r *RedisClient) SetBlock(block *tmtypes.Block) error {
	t0 := time.Now()
	if block == nil || block.Size() <= 0 {
		return fmt.Errorf("block is empty")
	}
	blockBytes, err := block.Marshal()
	if err != nil {
		return err
	}

	t1 := time.Now()
	compressBytes := r.compressBroker.DefaultCompress(blockBytes)

	t2 := time.Now()
	req := r.rdb.SetNX(context.Background(), setBlockKey(block.Height), compressBytes, TTL)

	t3 := time.Now()
	r.logger.Info("SetDeltas", "marshal", t1.Sub(t0), "compress", t2.Sub(t1), "setRedis", t3.Sub(t2),
		"blockBytes", len(blockBytes), "compressBytes", len(compressBytes))

	return req.Err()
}

func (r *RedisClient) SetDeltas(deltas *tmtypes.Deltas) error {
	t0 := time.Now()
	if deltas == nil || deltas.Size() == 0 {
		return fmt.Errorf("delta is empty")
	}
	deltaBytes, err := deltas.Marshal()
	if err != nil {
		return err
	}

	t1 := time.Now()
	compressBytes := r.compressBroker.DefaultCompress(deltaBytes)

	t2 := time.Now()
	req := r.rdb.SetNX(context.Background(), setDeltaKey(deltas.Height), compressBytes, TTL)

	t3 := time.Now()
	r.logger.Info("SetDeltas", "marshal", t1.Sub(t0), "compress", t2.Sub(t1), "setRedis", t3.Sub(t2),
		"deltaBytes", len(deltaBytes), "compressBytes", len(compressBytes))

	return req.Err()
}

func (r *RedisClient) GetBlock(height int64) (*tmtypes.Block, error) {
	t0 := time.Now()
	compressBytes, err := r.rdb.Get(context.Background(), setBlockKey(height)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	blockBytes := r.compressBroker.UnCompress(compressBytes)

	t2 := time.Now()
	block := &tmtypes.Block{}
	if err = block.Unmarshal(blockBytes); err != nil {
		return nil, err
	}

	t3 := time.Now()
	r.logger.Info("GetBlock", "getRedis", t1.Sub(t0), "uncompress", t2.Sub(t1), "unmarshal", t3.Sub(t2),
		"blockBytes", len(blockBytes), "compressBytes", len(compressBytes))
	return block, nil
}

func (r *RedisClient) GetDeltas(height int64) (*tmtypes.Deltas, error) {
	t0 := time.Now()
	compressBytes, err := r.rdb.Get(context.Background(), setDeltaKey(height)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	deltaBytes := r.compressBroker.UnCompress(compressBytes)

	t2 := time.Now()
	deltas := &tmtypes.Deltas{}
	if err = deltas.Unmarshal(deltaBytes); err != nil {
		return nil, err
	}

	t3 := time.Now()
	r.logger.Info("GetDeltas", "getRedis", t1.Sub(t0), "uncompress", t2.Sub(t1), "unmarshal", t3.Sub(t2),
		"deltaBytes", len(deltaBytes), "compressBytes", len(compressBytes))
	return deltas, nil
}

func setBlockKey(height int64) string {
	return "BH:" + strconv.Itoa(int(height))
}

func setDeltaKey(height int64) string {
	return "DH:" + strconv.Itoa(int(height))
}
