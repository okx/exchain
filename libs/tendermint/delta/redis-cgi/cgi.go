package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"strconv"
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
	if block == nil || block.Size() <= 0 {
		return fmt.Errorf("block is empty")
	}
	blockBytes, err := block.Marshal()
	if err != nil {
		return err
	}
	compressBytes := r.compressBroker.DefaultCompress(blockBytes)
	r.logger.Info("compress block", "blockBytes", len(blockBytes), "compressBytes", len(compressBytes))
	return r.rdb.SetNX(context.Background(), setBlockKey(block.Height), compressBytes, TTL).Err()
}

func (r *RedisClient) SetDeltas(deltas *tmtypes.Deltas) error {
	if deltas == nil || deltas.Size() == 0 {
		return fmt.Errorf("delta is empty")
	}
	deltaBytes, err := deltas.Marshal()
	if err != nil {
		return err
	}
	compressBytes := r.compressBroker.DefaultCompress(deltaBytes)
	r.logger.Info("compress delta", "deltaBytes", len(deltaBytes), "compressBytes", len(compressBytes))
	return r.rdb.SetNX(context.Background(), setDeltaKey(deltas.Height), compressBytes, TTL).Err()
}

func (r *RedisClient) GetBlock(height int64) (*tmtypes.Block, error) {
	compressBytes, err := r.rdb.Get(context.Background(), setBlockKey(height)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	blockBytes := r.compressBroker.UnCompress(compressBytes)
	r.logger.Info("uncompress block", "blockBytes", len(blockBytes), "compressBytes", len(compressBytes))
	block := &tmtypes.Block{}
	if err = block.Unmarshal(blockBytes); err != nil {
		return nil, err
	}
	return block, nil
}

func (r *RedisClient) GetDeltas(height int64) (*tmtypes.Deltas, error) {
	compressBytes, err := r.rdb.Get(context.Background(), setDeltaKey(height)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	deltaBytes := r.compressBroker.UnCompress(compressBytes)
	r.logger.Info("uncompress delta", "deltaBytes", len(deltaBytes), "compressBytes", len(compressBytes))
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
