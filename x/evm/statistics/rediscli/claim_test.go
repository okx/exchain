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

func TestHSet(t *testing.T) {
	GetInstance().Init()
	claim := &XenMint{
		Height:    15429182,
		BlockTime: time.Unix(time.Now().Unix()-100*24*60*60, 0),
		TxHash:    "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		Term:      30,
		Rank:      "11600774",
	}
	db := GetInstance().GetClientPool().Get()
	defer db.Close()
	GetInstance().InsertClaim(claim)

	reward := &XenClaimReward{
		Height:       15429182,
		BlockTime:    time.Now(),
		TxHash:       "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		RewardAmount: "122",
	}
	GetInstance().insertRewardSingleEx(reward)
	claim.UserAddr = "useraddr"
	GetInstance().InsertClaim(claim)
}

func Test_redisCli_parseXenMint(t *testing.T) {
	GetInstance().Init()
	mint := GetInstance().parseXenMint("0x1826080876d1dfbb06aa4f722876fec7b243b59c")
	t.Log(mint)
}

func Test_redisCli_parseXenReward(t *testing.T) {
	GetInstance().Init()
	mint := GetInstance().parseXenReward("0x1826080876d1dfbb06aa4f722876fec7b243b59c")
	t.Log(mint)
}
