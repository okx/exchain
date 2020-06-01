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

	TradePairKeyPrefix         = []byte{0x01}
	SavingKeyPrefix            = []byte{0x02}
	AccountKeyPrefix           = []byte{0x03}
	BorrowInfoKeyPrefix        = []byte{0x04}
	CalculateInterestKeyPrefix = []byte{0x05}
	DexWithdrawKeyPrefix       = []byte{0x06}
	DexWithdrawTimeKeyPrefix   = []byte{0x07}
)

func GetTradePairKey(product string) []byte {
	return append(TradePairKeyPrefix, []byte(product)...)
}

func GetAccountAddressKey(address sdk.AccAddress) []byte {
	return append(AccountKeyPrefix, address.Bytes()...)
}

func GetAccountAddressProductKey(address sdk.AccAddress, product string) []byte {
	return append(GetAccountAddressKey(address), []byte(product)...)
}

// GetDexWithdrawKey returns key of withdraw
func GetDexWithdrawKey(addr sdk.AccAddress) []byte {
	return append(DexWithdrawKeyPrefix, addr.Bytes()...)
}

// GetWithdrawTimeKey returns key of withdraw time
func GetWithdrawTimeKey(completeTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(completeTime)
	return append(DexWithdrawTimeKeyPrefix, bz...)
}

// GetDexWithdrawTimeAddressKey returns withdraw time address key
func GetDexWithdrawTimeAddressKey(completeTime time.Time, addr sdk.AccAddress) []byte {
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

func GetBorrowInfoAddressKey(address sdk.AccAddress) []byte {
	return append(BorrowInfoKeyPrefix, address.Bytes()...)
}

func GetBorrowInfoProductKey(address sdk.AccAddress, product string) []byte {
	return append(GetBorrowInfoAddressKey(address), []byte(product)...)
}

func GetBorrowInfoKey(address sdk.AccAddress, product string, blockHeight uint64) []byte {
	return append(GetBorrowInfoProductKey(address, product), sdk.Uint64ToBigEndian(blockHeight)...)
}

func GetCalculateInterestTimeKey(calculateTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(calculateTime)
	return append(CalculateInterestKeyPrefix, bz...)
}

func GetCalculateInterestKey(calculateTime time.Time, BorrowInfoKey []byte) []byte {
	return append(GetCalculateInterestTimeKey(calculateTime), BorrowInfoKey...)
}

func SplitCalculateInterestTimeKey(key []byte) (time.Time, []byte) {
	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}
	return endTime, key[1+lenTime:]
}
