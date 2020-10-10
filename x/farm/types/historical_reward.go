package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// PoolHistoricalRewards records the reward ratio of one user in a pool
type PoolHistoricalRewards struct {
	CumulativeRewardRatio sdk.DecCoins
	ReferenceCount        uint16
}

// NewPoolHistoricalRewards creates a new instance of PoolHistoricalRewards
func NewPoolHistoricalRewards(cumulativeRewardRatio sdk.DecCoins, referenceCount uint16) PoolHistoricalRewards {
	return PoolHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

// PoolCurrentPeriod records the current period
type PoolCurrentPeriod struct {
	StartBlockHeight             int64
	Period                       int64
	LastAmountYieldedNativeToken sdk.DecCoin
}

// NewPoolCurrentPeriod creates a new instance of PoolCurrentPeriod
func NewPoolCurrentPeriod(startBlockHeight, period int64, token sdk.DecCoin) PoolCurrentPeriod {
	return PoolCurrentPeriod{
		StartBlockHeight:             startBlockHeight,
		Period:                       period,
		LastAmountYieldedNativeToken: token,
	}
}
