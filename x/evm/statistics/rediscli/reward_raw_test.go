package rediscli

import (
	"testing"
	"time"
)

func Test_redisCli_insertRawReward(t *testing.T) {
	GetInstance().Init()
	reward := &XenClaimReward{
		Height:       15429182,
		BlockTime:    time.Now(),
		TxHash:       "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		RewardAmount: "122",
	}
	GetInstance().InsertReward(reward)

	reward = &XenClaimReward{
		Height:       15429183,
		BlockTime:    time.Now(),
		TxHash:       "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:     "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		RewardAmount: "122",
	}
	GetInstance().InsertReward(reward)

}
