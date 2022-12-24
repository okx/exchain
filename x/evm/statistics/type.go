package statistics

import "time"

type XenMint struct {
	Height     int64
	BlockTime  time.Time
	TxHash     string
	TxSender   string
	UserAddr   string
	Term       uint64
	Rank       uint64
	ExpireTime time.Time // BlockTime+Term+buffer(1d)
}

type XenClaimReward struct {
	TxHash       string
	TxSender     string
	UserAddr     string
	RewardAmount uint64
}
