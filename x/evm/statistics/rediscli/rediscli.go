package rediscli

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	redisAddr = ":6379"
)

var cli *redisCli

func init() {
	cli = &redisCli{}
}

func GetInstance() *redisCli {
	return cli
}

type redisCli struct {
	client *redis.Pool
	height int64
}

func (r *redisCli) Init() {
	var err error
	r.client = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}
	if err != nil {
		panic(err)
	}
}

func (r *redisCli) Close() {
	r.client.Close()
}

func (r *redisCli) GetClientPool() *redis.Pool {
	return r.client
}

func (r *redisCli) UpdateHeight() {
	db := r.client.Get()
	defer db.Close()
	claimHeight, _ := redis.Int64(db.Do("GET", "claim-height"))
	rewardHeight, _ := redis.Int64(db.Do("GET", "reward-height"))
	if claimHeight > rewardHeight {
		r.height = claimHeight
		return
	}
	r.height = rewardHeight
}
