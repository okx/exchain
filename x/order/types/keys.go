package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the order module
	ModuleName        = "order"
	DefaultParamspace = ModuleName
	DefaultCodespace  = ModuleName

	// QuerierRoute is the querier route for the order module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the order module
	RouterKey = ModuleName

	// query endpoints supported by the governance Querier
	QueryOrderDetail = "detail"
	QueryDepthBook   = "depthbook"
	QueryParameters  = "params"
	QueryStore       = "store"
	QueryDepthBookV2 = "depthbookV2"

	OrderStoreKey = ModuleName
)

var (
	// Keys for store prefixes

	// iterator keys
	OrderKey             = []byte{0x11}
	DepthbookKey         = []byte{0x12}
	OrderIDsKey          = []byte{0x13}
	PriceKey             = []byte{0x14}
	ExpireBlockHeightKey = []byte{0x15}
	OrderNumPerBlockKey  = []byte{0x16}

	// none iterator keys
	RecentlyClosedOrderIDsKey = []byte{0x17}
	LastExpiredBlockHeightKey = []byte{0x18}
	OpenOrderNumKey           = []byte{0x19}
	StoreOrderNumKey          = []byte{0x20}
)

func GetOrderKey(key string) []byte {
	return append(OrderKey, []byte(key)...)
}

func GetDepthbookKey(key string) []byte {
	return append(DepthbookKey, []byte(key)...)
}

func GetOrderIDsKey(key string) []byte {
	return append(OrderIDsKey, []byte(key)...)
}

func GetPriceKey(key string) []byte {
	return append(PriceKey, []byte(key)...)
}

func GetOrderNumPerBlockKey(blockHeight int64) []byte {
	return append(OrderNumPerBlockKey, sdk.Uint64ToBigEndian(uint64(blockHeight))...)
}

func GetExpireBlockHeightKey(blockHeight int64) []byte {
	return append(ExpireBlockHeightKey, sdk.Uint64ToBigEndian(uint64(blockHeight))...)
}

func FormatOrderIDsKey(product string, price sdk.Dec, side string) string {
	return fmt.Sprintf("%v:%v:%v", product, price.String(), side)
}

func GetKey(it sdk.Iterator) string {
	return string(it.Key()[1:])
}
