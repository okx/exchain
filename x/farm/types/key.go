package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "farm"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	FarmPoolPrefix              = []byte{0x01}
	pool2AddressPrefix          = []byte{0x02}
	address2PoolPrefix          = []byte{0x03}
	PoolsYieldNativeTokenPrefix = []byte{0x04}
)

func GetFarmPoolKey(poolName string) []byte {
	return append(FarmPoolPrefix, []byte(poolName)...)
}

func GetLockInfoKey(addr sdk.AccAddress, poolName string) []byte {
	return append(address2PoolPrefix, append(addr.Bytes(), []byte(poolName)...)...)
}

func SplitPoolsYieldNativeTokenKey(keyWithPrefix []byte) (poolName string) {
	return string(keyWithPrefix[1:])
}

// GetWhitelistMemberKey builds the key for a available pool name
func GetWhitelistMemberKey(poolName string) []byte {
	return append(PoolsYieldNativeTokenPrefix, []byte(poolName)...)
}
