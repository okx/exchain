package proto

import (
	"fmt"
	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// ProtocolDefinition is the struct of app-upgrade detail info
type ProtocolDefinition struct {
	Version   uint64  `json:"version"`
	Software  string  `json:"software"`
	Height    uint64  `json:"height"`
	Threshold sdk.Dec `json:"threshold"`
}

// NewProtocolDefinition creates a new instance of ProtocolDefinition
func NewProtocolDefinition(version uint64, software string, height uint64, threshold sdk.Dec) ProtocolDefinition {
	return ProtocolDefinition{
		version,
		software,
		height,
		threshold,
	}
}

// AppUpgradeConfig is the struct of app-upgrade-specific params
type AppUpgradeConfig struct {
	ProposalID  uint64             `json:"proposal_id"`
	ProtocolDef ProtocolDefinition `json:"protocol_def"`
}

// NewAppUpgradeConfig creates a new instance of AppUpgradeConfig
func NewAppUpgradeConfig(proposalID uint64, protocolDef ProtocolDefinition) AppUpgradeConfig {
	return AppUpgradeConfig{
		proposalID,
		protocolDef,
	}
}

// DefaultUpgradeConfig returns a default AppUpgradeConfig object
func DefaultUpgradeConfig(software string) AppUpgradeConfig {
	return AppUpgradeConfig{
		ProposalID:  uint64(0),
		ProtocolDef: NewProtocolDefinition(uint64(0), software, uint64(1), sdk.NewDecWithPrec(9, 1)),
	}
}

// VersionKeeper shows the expected behaviour of a version keeper
type VersionKeeper interface {
	GetCurrentVersionByStore(store sdk.KVStore) uint64
	GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig AppUpgradeConfig, found bool)
}

// ProtocolKeeper is designed for a protocol controller
type ProtocolKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewProtocolKeeper creates a new instance of ProtocolKeeper
func NewProtocolKeeper(key sdk.StoreKey) ProtocolKeeper {
	return ProtocolKeeper{key, cdc}
}

// GetCurrentVersionByStore gets the current version of protocol from store
func (pk ProtocolKeeper) GetCurrentVersionByStore(store sdk.KVStore) uint64 {
	bz := store.Get(currentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

// GetCurrentVersion gets the current version from context
func (pk ProtocolKeeper) GetCurrentVersion(ctx sdk.Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	return pk.GetCurrentVersionByStore(store)
}

// GetUpgradeConfigByStore gets the upgrade config from store
func (pk ProtocolKeeper) GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig AppUpgradeConfig, found bool) {
	bz := store.Get(upgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

// SetCurrentVersion sets current version
func (pk ProtocolKeeper) SetCurrentVersion(ctx sdk.Context, currentVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(currentVersion)
	store.Set(currentVersionKey, bz)
}

// GetLastFailedVersion gets last failed version
func (pk ProtocolKeeper) GetLastFailedVersion(ctx sdk.Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(lastFailedVersionKey)
	if bz == nil {
		return 0 // default value
	}
	var lastFailedVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastFailedVersion)
	return lastFailedVersion
}

// SetLastFailedVersion sets last failed version
func (pk ProtocolKeeper) SetLastFailedVersion(ctx sdk.Context, lastFailedVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(lastFailedVersion)
	store.Set(lastFailedVersionKey, bz)
}

// GetUpgradeConfig gets upgrade config
func (pk ProtocolKeeper) GetUpgradeConfig(ctx sdk.Context) (upgradeConfig AppUpgradeConfig, found bool) {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(upgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

// SetUpgradeConfig sets upgrade config
func (pk ProtocolKeeper) SetUpgradeConfig(ctx sdk.Context, upgradeConfig AppUpgradeConfig) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(upgradeConfigKey, bz)
}

// ClearUpgradeConfig removes the upgrade config in the store
func (pk ProtocolKeeper) ClearUpgradeConfig(ctx sdk.Context) {
	store := ctx.KVStore(pk.storeKey)
	store.Delete(upgradeConfigKey)
}

// IsValidVersion checks whether the version is available
func (pk ProtocolKeeper) IsValidVersion(ctx sdk.Context, version uint64) bool {
	currentVersion := pk.GetCurrentVersion(ctx)
	lastFailedVersion := pk.GetLastFailedVersion(ctx)
	return isValidVersion(currentVersion, lastFailedVersion, version)
}

// rule: new version should be currentVersion+1 or lastFailedVersion or lastFailedVersion+1
func isValidVersion(currentVersion uint64, lastFailedVersion uint64, version uint64) bool {
	if currentVersion >= lastFailedVersion {
		return currentVersion+1 == version
	}
	return lastFailedVersion == version || lastFailedVersion+1 == version

}

// String returns a human readable string representation of AppUpgradeConfig
func (auc AppUpgradeConfig) String() string {
	return fmt.Sprintf(`AppUpgradeConfig:
	ProposalID:			 %d
	ProtocolDefinition:  %v
`,
		auc.ProposalID, auc.ProtocolDef)
}
