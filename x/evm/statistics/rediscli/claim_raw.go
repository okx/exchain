package rediscli

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

func (r *redisCli) insertRawMint(mint *XenMint) {
	db := r.client.Get()
	defer db.Close()

	key := "c" + mint.UserAddr + "_" + strconv.Itoa(int(mint.Height))
	_, err := redis.Int(db.Do("HSET", key, "height", mint.Height,
		"btime", mint.BlockTime.Unix(), "txhash", mint.TxHash, "term", mint.Term, "rank", mint.Rank))
	if err != nil {
		panic(err)
	}
	r.insertUserAddr(mint, key, int(mint.Height))
}

func (r *redisCli) insertUserAddr(mint *XenMint, key string, score int) {
	db := r.client.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Do("ZADD", mint.UserAddr, score, key)
	if err != nil {
		panic(err)
	}
}
