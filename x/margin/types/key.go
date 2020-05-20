package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "margin"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	lenTime = len(sdk.FormatTimeBytes(time.Now()))

	TradePairKeyPrefix = []byte{0x01}
	SavingKeyPrefix    = []byte{0x02}

	MagrinAssetKey = []byte{0x03}

	WithdrawKeyPrefix     = []byte{0x05}
	WithdrawTimeKeyPrefix = []byte{0x06}
)

func GetTradePairKey(product string) []byte {
	return append(TradePairKeyPrefix, []byte(product)...)
}

func GetMarginAllAssetKey(address string) []byte {
	return append(MagrinAssetKey, []byte(address)...)
}

func GetMarginProductAssetKey(address, product string) []byte {
	return append(GetMarginAllAssetKey(address), []byte(product)...)
}

// GetWithdrawKey returns key of withdraw
func GetWithdrawKey(addr sdk.AccAddress) []byte {
	return append(WithdrawKeyPrefix, addr.Bytes()...)
}

// GetWithdrawTimeKey returns key of withdraw time
func GetWithdrawTimeKey(completeTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(completeTime)
	return append(WithdrawTimeKeyPrefix, bz...)
}

// GetWithdrawTimeAddressKey returns withdraw time address key
func GetWithdrawTimeAddressKey(completeTime time.Time, addr sdk.AccAddress) []byte {
	return append(GetWithdrawTimeKey(completeTime), addr.Bytes()...)
}

//SplitWithdrawTimeKey splits the key and returns the complete time and address
func SplitWithdrawTimeKey(key []byte) (time.Time, sdk.AccAddress) {
	if len(key[1:]) != lenTime+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+sdk.AddrLen))
	}
	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}
	delAddr := sdk.AccAddress(key[1+lenTime:])
	return endTime, delAddr
}

func GetSavingKey(product string) []byte {
	return append(SavingKeyPrefix, []byte(product)...)
}
