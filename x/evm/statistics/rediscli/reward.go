package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *redisCli) InsertReward(reward *XenClaimReward) {
	if reward.Height < int64(r.height) {
		return
	}
	exists, err := redis.Int(r.client.Do("EXISTS", reward.UserAddr))
	if (err != nil || exists == 0) && reward.Height != r.height {
		panic(fmt.Sprintf("error %v or no exists %v", err, reward.UserAddr))
	}

	if del, err := redis.Int(r.client.Do("DEL", reward.UserAddr)); (err != nil || del == 0) && reward.Height != r.height {
		panic(fmt.Sprintf("del %v error %v %v", reward, err, del))
	}
	if reward.Height > r.height {
		r.client.Do("SET", "reward-height", reward.Height)
	}
}
