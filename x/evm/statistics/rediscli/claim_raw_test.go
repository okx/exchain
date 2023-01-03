package rediscli

import (
	"testing"
	"time"
)

func Test_redisCli_insertRawMint(t *testing.T) {
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
	GetInstance().InsertClaim(claim)

	claim = &XenMint{
		Height:    15429183,
		BlockTime: time.Unix(time.Now().Unix()-100*24*60*60, 0),
		TxHash:    "0x811991657398dda93c4b2db124c9ddcd85899e886b9aee2041bae976a5595a6c",
		TxSender:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		UserAddr:  "0x1826080876d1dfbb06aa4f722876fec7b243b59c",
		Term:      30,
		Rank:      "11600774",
	}
	GetInstance().InsertClaim(claim)
}
