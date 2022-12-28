package rediscli

import (
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
