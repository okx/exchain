package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"time"
)

func (r *redisCli) InsertReward(reward *XenClaimReward) {
	r.insertRawReward(reward)
	//	if global.RedisPassword == "" {
	//		r.insertRewardSingle(reward)
	//	} else {
	//		r.insertRewardMulti(reward)
	//	}
}

func (r *redisCli) insertRewardSingleEx(reward *XenClaimReward) {
	db := r.client.Get()
	defer func() {
		db.Do("SELECT", 0)
		db.Close()
	}()
	_, err := db.Do("SELECT", 1)
	if err != nil {
		log.Printf("claim err %v %v %v \n", err, reward.TxHash, reward.UserAddr)
		return
	}
	_, err = redis.Int(db.Do("HSET", "r"+reward.UserAddr, "height", reward.Height,
		"btime", reward.BlockTime.Unix(), "txsender", reward.TxSender, "txhash", reward.TxHash, "amount", reward.RewardAmount, "reward", 0))

	if err != nil {
		log.Printf("claim err %v %v %v \n", err, reward.TxHash, reward.UserAddr)
	}
}

func (r *redisCli) insertRewardSingle(reward *XenClaimReward) {
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

func (r *redisCli) insertRewardMulti(reward *XenClaimReward) {
	db := r.client.Get()
	defer db.Close()

	exists, _ := redis.Int(db.Do("EXISTS", reward.UserAddr))
	if exists == 1 {
		rewardDB := r.parseXenReward(reward.UserAddr)
		if rewardDB.Height < reward.Height {
			r.insertRewardSingleEx(reward)
		}
	} else {
		r.insertRewardSingleEx(reward)
	}
}

func (r *redisCli) parseXenReward(uaddr string) *XenClaimReward {
	userAddr := "r" + uaddr
	db := r.client.Get()
	defer func() {
		db.Do("SELECT", 0)
		db.Close()
	}()
	_, err := db.Do("SELECT", 1)
	if err != nil {
		log.Printf("parse xen select error %v %v\n", err, userAddr)
	}
	mintValues, _ := redis.StringMap(db.Do("HGETALL", userAddr))
	var ret XenClaimReward
	ret.UserAddr = userAddr
	for key, value := range mintValues {
		parseReward(&ret, key, value)
	}

	return &ret
}

func parseReward(reward *XenClaimReward, key, value string) {
	switch key {
	case "height":
		height, _ := strconv.Atoi(value)
		reward.Height = int64(height)
	case "txhash":
		reward.TxHash = value
	case "btime":
		utc, _ := strconv.Atoi(value)
		tim := time.Unix(int64(utc), 0)
		reward.BlockTime = tim
	}
}
