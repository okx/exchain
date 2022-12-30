package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *redisCli) InsertReward(reward *XenClaimReward) {
	if reward.Height < int64(r.height) {
		return
	}
	db := r.client.Get()
	defer db.Close()
	//	exists, err := redis.Int(db.Do("EXISTS", reward.UserAddr))
	//	if (err != nil || exists == 0) && reward.Height != r.height {
	//		panic(fmt.Sprintf("error %v or no exists %v", err, reward.UserAddr))
	//	}

	if _, err := redis.Int(db.Do("DEL", reward.UserAddr)); (err != nil) && reward.Height != r.height {
		panic(fmt.Sprintf("del %v error %v ", reward, err))
	}
	if reward.Height > r.height {
		db.Do("SET", "reward-height", reward.Height)
	}
	db.Do("INCR", "reward-counter")
}
