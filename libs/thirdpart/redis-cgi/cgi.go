package redis_cgi

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

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
	_, err = r.rdb.SetNX(context.Background(), setBlockKey(block.Height), blockBytes, time.Second)
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