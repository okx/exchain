package types

const (
	// ModuleName is the name of the module
	ModuleName = "ammswap"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	// QuerySwapTokenPair query endpoints supported by the swap Querier
	QuerySwapTokenPair = "swapTokenPair"

	QuerySwapTokenPairs = "swapTokenPairs"

	QueryRedeemableAssets = "queryRedeemableAssets"

	QueryParams = "params"
)

var (
	// TokenPairPrefixKey to be used for KVStore
	TokenPairPrefixKey = []byte{0x01}
)

// nolint
func GetTokenPairKey(key string) []byte {
	return append(TokenPairPrefixKey, []byte(key)...)
}
