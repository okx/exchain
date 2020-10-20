package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "farm"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// MintFarmingAccount as module account to be used for saving all mint farming tokens
	MintFarmingAccount = "mint_farming_account"

	// YieldFarmingAccount as module account to be used for saving all yield farming tokens
	YieldFarmingAccount = "yield_farming_account"

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	// Byte length of period occupied
	PeriodByteArrayLength = 8
)

var (
	FarmPoolPrefix              = []byte{0x01}
	Pool2AddressPrefix          = []byte{0x02}
	Address2PoolPrefix          = []byte{0x03}
	PoolsYieldNativeTokenPrefix = []byte{0x04}
	PoolHistoricalRewardsPrefix = []byte{0x05}
	PoolCurrentRewardsPrefix    = []byte{0x06}
)

const (
	poolNameFromLockInfoKeyIndex = sdk.AddrLen + 1
)

func GetFarmPoolKey(poolName string) []byte {
	return append(FarmPoolPrefix, []byte(poolName)...)
}

func GetAddressInFarmPoolKey(poolName string, addr sdk.AccAddress) []byte {
	return append(Pool2AddressPrefix, append([]byte(poolName), addr.Bytes()...)...)
}

func GetLockInfoKey(addr sdk.AccAddress, poolName string) []byte {
	return append(Address2PoolPrefix, append(addr.Bytes(), []byte(poolName)...)...)
}

func SplitPoolsYieldNativeTokenKey(keyWithPrefix []byte) (poolName string) {
	return string(keyWithPrefix[1:])
}

// GetWhitelistMemberKey builds the key for a available pool name
func GetWhitelistMemberKey(poolName string) []byte {
	return append(PoolsYieldNativeTokenPrefix, []byte(poolName)...)
}

// SplitPoolNameFromLockInfoKey splits the pool name out from a LockInfoKey
func SplitPoolNameFromLockInfoKey(lockInfoKey []byte) string {
	return string(lockInfoKey[poolNameFromLockInfoKeyIndex:])
}

// GetPoolHistoricalRewardsKey gets the key for a pool's historical reward
func GetPoolHistoricalRewardsKey(poolName string, period uint64) []byte {
	b := make([]byte, PeriodByteArrayLength)
	binary.LittleEndian.PutUint64(b, period)
	return append(PoolHistoricalRewardsPrefix, append([]byte(poolName), b...)...)
}

// GetPoolHistoricalRewardsPrefix gets the prefix key with pool name for a pool's historical rewards
func GetPoolHistoricalRewardsPrefix(poolName string) []byte {
	return append(PoolHistoricalRewardsPrefix, []byte(poolName)...)
}

// GetPoolCurrentRewardsKey gets the key for a pool's current period reward
func GetPoolCurrentRewardsKey(poolName string) []byte {
	return append(PoolCurrentRewardsPrefix, []byte(poolName)...)
}
