package types

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
	TradePairKeyPrefix = []byte{0x01}
	MagrinAssetKey     = []byte{0x02}
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
