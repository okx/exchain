package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// PoolHistoricalRewards records the reward ratio of one user in a pool
type PoolHistoricalRewards struct {
	CumulativeRewardRatio sdk.SysCoins
	ReferenceCount        uint16
}

// NewPoolHistoricalRewards creates a new instance of PoolHistoricalRewards
func NewPoolHistoricalRewards(cumulativeRewardRatio sdk.SysCoins, referenceCount uint16) PoolHistoricalRewards {
	return PoolHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

// PoolCurrentRewards records the current period
type PoolCurrentRewards struct {
	StartBlockHeight int64
	Period           uint64
	Rewards          sdk.SysCoins
}

// NewPoolCurrentRewards creates a new instance of PoolCurrentRewards
func NewPoolCurrentRewards(startBlockHeight int64, period uint64, token sdk.SysCoins) PoolCurrentRewards {
	return PoolCurrentRewards{
		StartBlockHeight: startBlockHeight,
		Period:           period,
		Rewards:          token,
	}
}
