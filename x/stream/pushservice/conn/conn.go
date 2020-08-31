package conn

import (
	"fmt"

	// "github.com/apache/pulsar/pulsar-client-go/pulsar"

	"github.com/go-redis/redis"
	"github.com/tendermint/tendermint/libs/log"
)

type Client struct {
	redisCli *redis.Client
	log      log.Logger
}

func NewClient(redisUrl, redisPassword string, db int, log log.Logger) (client *Client, err error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword, // no password set
		DB:       db,            // use default DB
	})

	client = &Client{redisCli: rc, log: log}
	return client, err
}

//PrivatePub push data to private topic
func (c Client) PrivatePub(key, val string) (err error) {
	logger := c.log.With("module", "redis-Client")
	err = c.redisCli.Publish(key, val).Err()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("Redis Publish error")
		} else {
			return err
		}
	}
	return
}

//PublicPub push data to public topic
func (c Client) PublicPub(key, val string) (err error) {
	logger := c.log.With("module", "redis-Client")
	err = c.redisCli.Publish(key, val).Err()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("Redis Publish error")
		} else {
			return err
		}
	}
	return
}

//DepthPub push data to depth topic
func (c Client) DepthPub(key, val string) (err error) {
	logger := c.log.With("module", "redis-Client")
	err = c.redisCli.Publish(key, val).Err()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("Redis Publish error")
		} else {
			return err
		}
	}
	return
}

//Set push data to redis, only used by public channel
func (c Client) Set(key, val string) (err error) {
	logger := c.log.With("module", "redis-Client")
	err = c.redisCli.Set(key, val, 0).Err()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("Redis set", "key not found, set a new k-v")
		} else {
			return err
		}
	}
	return nil
}

//Get get data from redis
func (c Client) Get(key string) (val string, err error) {
	val, err = c.redisCli.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get val for key=%s, err=%s", key, err.Error())
	}
	return val, nil
}

//MGet get data from redis, many keys once
func (c Client) MGet(keys []string) (vals []interface{}, err error) {
	vals, err = c.redisCli.MGet(keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get vals for keys=%v, err=%s", keys, err.Error())
	}
	return vals, nil
}

//Close close redis connect to server
func (c Client) Close() error {
	return c.redisCli.Close()
}

//HGetAll
func (c Client) HGetAll(key string) (vals map[string]string, err error) {
	vals, err = c.redisCli.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to HGETALL vals for key=%v, err=%s", key, err.Error())
	}
	return vals, nil
}
