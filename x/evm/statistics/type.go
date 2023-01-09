package statistics

import "time"

type XenMint struct {
	Height    int64
	BlockTime time.Time
	TxHash    string
	TxSender  string
	To        string
	UserAddr  string
	Term      int64
	Rank      string
}

type XenClaimReward struct {
	Height       int64
	BlockTime    time.Time
	TxHash       string
	TxSender     string
	To           string
	UserAddr     string
	RewardAmount string
}
