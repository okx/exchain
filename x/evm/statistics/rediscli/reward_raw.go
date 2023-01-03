package rediscli

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

func (r *redisCli) insertRawReward(reward *XenClaimReward) {
	db := r.client.Get()
	defer db.Close()
	db.Do("SELECT", 1)

	key := "r" + reward.UserAddr + "_" + strconv.Itoa(int(reward.Height))
	_, err := redis.Int(db.Do("HSET", key, "height", reward.Height,
		"btime", reward.BlockTime.Unix(), "txhash", reward.TxHash))
	if err != nil {
		panic(err)
	}
	r.insertUserAddrReward(reward, key, int(reward.Height))
}

func (r *redisCli) insertUserAddrReward(reward *XenClaimReward, key string, score int) {
	db := r.client.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 3)
	if err != nil {
		panic(err)
	}
	_, err = db.Do("ZADD", reward.UserAddr, score, key)
	if err != nil {
		panic(err)
	}
}
