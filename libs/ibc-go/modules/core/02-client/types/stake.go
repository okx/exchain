package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	"time"
)

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}
