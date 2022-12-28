package rediscli

import (
	"github.com/gomodule/redigo/redis"
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
	client redis.Conn
	height int64
}

func (r *redisCli) Init() {
	var err error
	r.client, err = redis.Dial("tcp", redisAddr)
	if err != nil {
		panic(err)
	}
}

func (r *redisCli) Close() {
	r.client.Close()
}

func (r *redisCli) GetRawClient() redis.Conn {
	return r.client
}

func (r *redisCli) UpdateHeight() {
	claimHeight, _ := redis.Int64(r.client.Do("GET", "claim-height"))
	rewardHeight, _ := redis.Int64(r.client.Do("GET", "reward-height"))
	if claimHeight > rewardHeight {
		r.height = claimHeight
		return
	}
	r.height = rewardHeight
}
