package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
)

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}

// UpgradeKeeper expected upgrade keeper
type UpgradeKeeper interface {
	ClearIBCState(ctx sdk.Context, lastHeight int64)
	GetUpgradePlan(ctx sdk.Context) (plan upgrade.Plan, havePlan bool)
	GetUpgradedClient(ctx sdk.Context, height int64) ([]byte, bool)
	SetUpgradedClient(ctx sdk.Context, planHeight int64, bz []byte) error
	GetUpgradedConsensusState(ctx sdk.Context, lastHeight int64) ([]byte, bool)
	SetUpgradedConsensusState(ctx sdk.Context, planHeight int64, bz []byte) error
	ScheduleUpgrade(ctx sdk.Context, plan upgrade.Plan) error
}
