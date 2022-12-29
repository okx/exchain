package rediscli

import (
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func Test_redisCli_InsertClaim(t *testing.T) {
	GetInstance().Init()
	GetInstance().InsertClaim(&XenMint{
		Height:    1,
		BlockTime: time.Now(),
		TxHash:    "txhash",
		TxSender:  "txsender",
		UserAddr:  "useraddr",
		Term:      1,
		Rank:      "111",
	})

	GetInstance().InsertReward(&XenClaimReward{
		Height:       2,
		BlockTime:    time.Now(),
		TxHash:       "txhash",
		TxSender:     "txsender",
		UserAddr:     "useraddr",
		RewardAmount: "1",
	})
}

func TestHSet(t *testing.T) {
	GetInstance().Init()
	claim := &XenMint{
		Height:    15429182,
		BlockTime: time.Now(),
		TxHash:    "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		Term:      30,
		Rank:      "11600774",
	}
	db := GetInstance().GetClientPool().Get()
	defer db.Close()
	_, err := redis.Int(db.Do("HSET", claim.UserAddr, "height", claim.Height,
		"btime", claim.BlockTime.Unix(), "txhash", claim.TxHash, "term", claim.Term, "rank", claim.Rank, "reward", 0))
	if err != nil {
		panic(err)
	}
}
