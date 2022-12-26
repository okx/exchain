package statistics

import "time"

type XenMint struct {
	Height    int64
	BlockTime time.Time
	TxHash    string
	TxSender  string
	UserAddr  string
	Term      int64
	Rank      int64
}

type XenClaimReward struct {
	Height       int64
	BlockTime    time.Time
	TxHash       string
	TxSender     string
	UserAddr     string
	RewardAmount int64
}
