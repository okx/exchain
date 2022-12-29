package rediscli

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	// ExpireDate date --date='2023-02-01 00:00:00' +%s
	// TermThreshold = 1675209600

	// 2022/12/29 02:00:43 lrp
	TermThreshold = 1672279243
)

func (r *redisCli) InsertClaim(claim *XenMint) {

	if claim.Height < int64(r.height) {
		return
	}
	if (claim.BlockTime).Add(time.Duration(claim.Term)*time.Duration(24)*time.Hour).Unix() > TermThreshold {
		return
	}
	db := r.client.Get()
	defer db.Close()
	//	exists, err := redis.Int(r.client.Do("EXISTS", claim.UserAddr))
	//	if (err != nil || exists == 1) && claim.Height != r.height {
	//		panic(fmt.Sprintf("error %v or exists %v", err, claim.UserAddr))
	//	}

	_, err := redis.Int(db.Do("HSET", claim.UserAddr, "height", claim.Height,
		"btime", claim.BlockTime.Unix(), "txhash", claim.TxHash, "term", claim.Term, "rank", claim.Rank, "reward", 0))
	if err != nil && claim.Height != r.height {
		panic(fmt.Sprintf("hset %v error %v", claim, err))
	}

	if claim.Height > r.height {
		db.Do("SET", "claim-height", claim.Height)
	}
	db.Do("INCR", "mint-counter")
}