package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the dex module
	ModuleName = "dex"
	// DefaultParamspace defines default param space
	DefaultParamspace = ModuleName
	// DefaultCodespace defines default code space
	DefaultCodespace = ModuleName
	// QuerierRoute is the querier route for the dex module
	QuerierRoute = ModuleName
	// RouterKey is the msg router key for the dex module
	RouterKey = ModuleName
	// StoreKey is the string store representation
	StoreKey = ModuleName
	// TokenPairStoreKey is the token pair store key
	TokenPairStoreKey = "token_pair"

	// QueryProductsDelisting defines delisting query route path
	QueryProductsDelisting = "products_delisting"
	// QueryProducts defines products query route path
	QueryProducts = "products"
	// QueryDeposits defines deposits query route path
	QueryDeposits = "deposits"
	// QueryMatchOrder defines match-order query route path
	QueryMatchOrder = "match-order"
	// QueryParameters defines 	QueryParameters = "params" query route path
	QueryParameters = "params"
	// QueryOperator defines operator query route path
	QueryOperator = "operator"
	// QueryOperators defines operators query route path
	QueryOperators = "operators"
)

var (
	lenTime = len(sdk.FormatTimeBytes(time.Now()))

	// TokenPairKey is the store key for token pair
	TokenPairKey = []byte{0x01}
	// TokenPairNumberKey is the store key for token pair num
	TokenPairNumberKey = []byte{0x02}
	// TokenPairLockKeyPrefix is the store key  prefix for token pair
	TokenPairLockKeyPrefix = []byte{0x03}

	// DEXOperatorKeyPrefix is the store key prefix for DEXOperator
	DEXOperatorKeyPrefix = []byte{0x04}

	// DEXOperatorTokenPairCntPrefix is the store key prefix for DEXOperator's tokenpair count
	DEXOperatorTokenPairCntPrefix = []byte{0x05}

	// PrefixWithdrawAddressKey is the store key for withdraw address
	PrefixWithdrawAddressKey = []byte{0x53}
	// PrefixWithdrawTimeKey is the store key for withdraw time
	PrefixWithdrawTimeKey = []byte{0x54}
	// PrefixUserTokenPairKey is the store key for user token pair num
	PrefixUserTokenPairKey = []byte{0x06}
)

// GetUserTokenPairAddressPrefix returns token pair address prefix key
func GetUserTokenPairAddressPrefix(owner sdk.AccAddress) []byte {
	return append(PrefixUserTokenPairKey, owner.Bytes()...)
}

// GetUserTokenPairAddress returns token pair address key
func GetUserTokenPairAddress(owner sdk.AccAddress, assertPair string) []byte {
	return append(GetUserTokenPairAddressPrefix(owner), []byte(assertPair)...)
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

// GetLockProductKey returns key of token pair
func GetLockProductKey(product string) []byte {
	return append(TokenPairLockKeyPrefix, []byte(product)...)
}

// GetKey returns keys between index 1 to the end
func GetKey(it sdk.Iterator) string {
	return string(it.Key()[1:])
}

// GetOperatorAddressKey returns key of operator address
func GetOperatorAddressKey(addr sdk.AccAddress) []byte {
	return append(DEXOperatorKeyPrefix, addr.Bytes()...)
}
