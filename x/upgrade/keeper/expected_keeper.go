package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common/proto"
	"github.com/okex/okexchain/x/staking/exported"

	//"github.com/okex/okexchain/x/staking/types"
	"github.com/okex/okexchain/x/staking/types"
)

// BankKeeper shows the expected action of bank keeper in this module
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// StakingKeeper shows the expected action of staking keeper in this module
type StakingKeeper interface {
	IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool))
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator types.Validator, found bool)
}

// ProtocolKeeper shows the expected action of proto keeper in this module
type ProtocolKeeper interface {
	IsValidVersion(ctx sdk.Context, version uint64) bool
	GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig proto.AppUpgradeConfig, found bool)
	GetUpgradeConfig(ctx sdk.Context) (upgradeConfig proto.AppUpgradeConfig, found bool)
	SetUpgradeConfig(ctx sdk.Context, upgradeConfig proto.AppUpgradeConfig)
	GetCurrentVersion(ctx sdk.Context) uint64
	SetLastFailedVersion(ctx sdk.Context, lastFailedVersion uint64)
	SetCurrentVersion(ctx sdk.Context, currentVersion uint64)
	ClearUpgradeConfig(ctx sdk.Context)
	GetLastFailedVersion(ctx sdk.Context) uint64
}
