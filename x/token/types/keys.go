package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "token"

	DefaultParamspace = ModuleName
	DefaultCodespace  = ModuleName

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName

	KeyLock = "lock"
	KeyMint = "mint"

	// query endpoints supported by the governance Querier
	QueryInfo       = "info"
	QueryTokens     = "tokens"
	QueryParameters = "params"
	QueryCurrency   = "currency"
	QueryAccount    = "accounts"
	QueryKeysNum    = "store"

	QueryAccountV2 = "accountsV2"
	QueryTokensV2  = "tokensV2"
	QueryTokenV2   = "tokenV2"
)

var (
	TokenKey           = []byte{0x00} // the address prefix of the token's symbol
	TokenNumberKey     = []byte{0x01} // key for token number address
	LockKey            = []byte{0x02} // the address prefix of the locked coins
	LockedFeeKey       = []byte{0x04} // the address prefix of the locked order fee coins
	PrefixUserTokenKey = []byte{0x03} // the address prefix of the user-token relationship
)

func GetUserTokenPrefix(owner sdk.AccAddress) []byte {
	return append(PrefixUserTokenKey, owner.Bytes()...)
}

func GetUserTokenKey(owner sdk.AccAddress, symbol string) []byte {
	return append(GetUserTokenPrefix(owner), []byte(symbol)...)
}

func GetTokenAddress(symbol string) []byte {
	return append(TokenKey, []byte(symbol)...)
}

func GetLockAddress(addr sdk.AccAddress) []byte {
	return append(LockKey, addr.Bytes()...)
}

// GetLockFeeAddress gets the key for the lock fee information with address
func GetLockFeeAddress(addr sdk.AccAddress) []byte {
	return append(LockedFeeKey, addr.Bytes()...)
}

//// Key for getting a specific proposal from the store
//func keyDexListAsset(asset string) []byte {
//	return []byte(fmt.Sprintf("asset:%s", asset))
//}
