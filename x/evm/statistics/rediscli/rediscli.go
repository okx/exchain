package rediscli

import "github.com/gomodule/redigo/redis"

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
