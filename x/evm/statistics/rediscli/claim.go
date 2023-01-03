package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

const (
	// ExpireDate date --date='2023-02-01 00:00:00' +%s
	// TermThreshold = 1675209600

	// 2022/12/29 02:00:43 lrp
	TermThreshold = 1672279243
)

func (r *redisCli) InsertClaim(claim *XenMint) {
	r.insertRawMint(claim)
	//	if global.RedisPassword == "" {
	//		r.insertSingle(claim, true)
	//	} else {
	//		r.insertMulti(claim)
	//	}
}

func (r *redisCli) parseXenMint(uaddr string) *XenMint {
	userAddr := "c" + uaddr
	db := r.client.Get()
	defer db.Close()
	mintValues, _ := redis.StringMap(db.Do("HGETALL", userAddr))
	var ret XenMint
	ret.UserAddr = userAddr
	for key, value := range mintValues {
		parseClaim(&ret, key, value)
	}

	return &ret
}

func (r *redisCli) insertMulti(claim *XenMint) {
	if (claim.BlockTime).Add(time.Duration(claim.Term+8)*time.Duration(24)*time.Hour).Unix() > TermThreshold {
		return
	}
	db := r.client.Get()
	defer db.Close()

	exists, _ := redis.Int(db.Do("EXISTS", claim.UserAddr))
	if exists == 1 {
		mint := r.parseXenMint(claim.UserAddr)
		if mint.Height < claim.Height {
			r.insertSingle(claim, false)
		}
	} else {
		r.insertSingle(claim, false)
	}
}

func (r *redisCli) insertSingle(claim *XenMint, updateMeta bool) {
	if claim.Height < int64(r.height) {
		return
	}
	if (claim.BlockTime).Add(time.Duration(claim.Term+8)*time.Duration(24)*time.Hour).Unix() > TermThreshold {
		return
	}
	db := r.client.Get()
	defer db.Close()
	//	exists, err := redis.Int(r.client.Do("EXISTS", claim.UserAddr))
	//	if (err != nil || exists == 1) && claim.Height != r.height {
	//		panic(fmt.Sprintf("error %v or exists %v", err, claim.UserAddr))
	//	}

	userAddrKey := claim.UserAddr
	if !updateMeta {
		userAddrKey = "c" + claim.UserAddr
	}
	_, err := redis.Int(db.Do("HSET", userAddrKey, "height", claim.Height,
		"btime", claim.BlockTime.Unix(), "txhash", claim.TxHash, "term", claim.Term, "rank", claim.Rank, "reward", 0))
	if err != nil && claim.Height != r.height {
		panic(fmt.Sprintf("hset %v error %v", claim, err))
	}

	if updateMeta {
		if claim.Height > r.height {
			db.Do("SET", "claim-height", claim.Height)
		}
		db.Do("INCR", "mint-counter")
	}
}

func parseClaim(claim *XenMint, key, value string) {
	switch key {
	case "height":
		height, _ := strconv.Atoi(value)
		claim.Height = int64(height)
	case "txhash":
		claim.TxHash = value
	case "term":
		term, _ := strconv.Atoi(value)
		claim.Term = int64(term)
	case "btime":
		utc, _ := strconv.Atoi(value)
		tim := time.Unix(int64(utc), 0)
		claim.BlockTime = tim
	}
}
