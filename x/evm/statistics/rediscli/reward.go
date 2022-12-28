package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *redisCli) InsertReward(reward *XenClaimReward) {
	exists, err := redis.Int(r.client.Do("EXISTS", reward.UserAddr))
	if err != nil || exists == 0 {
		panic(fmt.Sprintf("error %v or no exists %v", err, reward.UserAddr))
	}

	if del, err := redis.Int(r.client.Do("DEL", reward.UserAddr)); err != nil || del == 0 {
		panic(fmt.Sprintf("del %v error %v %v", reward, err, del))
	}
	if dup, err := redis.Int(r.client.Do("SADD", fmt.Sprintf("reward-%d", reward.Height), reward.UserAddr)); err != nil || dup == 0 {
		panic(fmt.Sprintf("sadd %v error %v or dup add %v", reward, err, dup))
	}
	r.client.Do("DEL", fmt.Sprintf("reward-%d", reward.Height-3))
}
