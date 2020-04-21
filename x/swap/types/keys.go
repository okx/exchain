package types

const (
	// ModuleName is the name of the module
	ModuleName = "swap"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	TokenPairPrefixKey = []byte{0x01}
)

// nolint
func GetTokenPairKey(key string) []byte {
	return append(TokenPairPrefixKey, []byte(key)...)
}
