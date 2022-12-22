package rediscli

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

func (r *redisCli) insertRawMint(mint *XenMint) {
	db := r.client.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}

	key := "c" + mint.UserAddr + "_" + strconv.Itoa(int(mint.Height))
	_, err = redis.Int(db.Do("HSET", key, "height", mint.Height,
		"btime", mint.BlockTime.Unix(), "txhash", mint.TxHash, "term", mint.Term, "rank", mint.Rank))
	if err != nil {
		panic(err)
	}
	r.insertUserAddr(mint, int(mint.Height))
}

func (r *redisCli) insertUserAddr(mint *XenMint, height int) {
	db := r.client.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Do("ZADD", "c"+mint.UserAddr, height, height)
	if err != nil {
		panic(err)
	}
}
