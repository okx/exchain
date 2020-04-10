package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the dex module
	ModuleName        = "dex"
	DefaultParamspace = ModuleName
	DefaultCodespace  = ModuleName

	// QuerierRoute is the querier route for the dex module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the dex module
	RouterKey = ModuleName

	// StoreKey is the string store representation
	StoreKey = ModuleName

	TokenPairStoreKey      = "token_pair"
	QueryProductsDelisting = "products_delisting"

	QueryProducts   = "products"
	QueryDeposits   = "deposits"
	QueryMatchOrder = "match-order"
	QueryParameters = "params"
)

var (
	lenTime = len(sdk.FormatTimeBytes(time.Now()))

	TokenPairKey             = []byte{0x01} // the address prefix of the token pair's symbol
	TokenPairNumberKey       = []byte{0x02} // key for token pair number address
	TokenPairLockKeyPrefix   = []byte{0x03}
	PrefixWithdrawAddressKey = []byte{0x53}
	PrefixWithdrawTimeKey    = []byte{0x54}
	PrefixUserTokenPairKey   = []byte{0x06}
)

func GetUserTokenPairAddressPrefix(Owner sdk.AccAddress) []byte {
	return append(PrefixUserTokenPairKey, Owner.Bytes()...)
}

func GetUserTokenPairAddress(Owner sdk.AccAddress, assertPair string) []byte {
	return append(GetUserTokenPairAddressPrefix(Owner), []byte(assertPair)...)
}

// GetTokenPairAddress returns store key of token pair
func GetTokenPairAddress(key string) []byte {
	return append(TokenPairKey, []byte(key)...)
}

// GetWithdrawAddressKey returns key of withdraw address
func GetWithdrawAddressKey(addr sdk.AccAddress) []byte {
	return append(PrefixWithdrawAddressKey, addr.Bytes()...)
}

// GetWithdrawTimeKey returns key of withdraw time
func GetWithdrawTimeKey(completeTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(completeTime)
	return append(PrefixWithdrawTimeKey, bz...)
}

// GetWithdrawTimeAddressKey returns
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

// GetProductKey returns key of token pair
func GetLockProductKey(product string) []byte {
	return append(TokenPairLockKeyPrefix, []byte(product)...)
}

// GetKey returns keys between index 1 to the end
func GetKey(it sdk.Iterator) string {
	return string(it.Key()[1:])
}
