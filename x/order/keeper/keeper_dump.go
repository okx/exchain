package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/order/types"
	"github.com/tendermint/tendermint/libs/log"
)

// DumpStore dumps all key-value message from KVStore
func (k Keeper) DumpStore(ctx sdk.Context) {

	logger := ctx.Logger().With("module", "order")

	orderStore := ctx.KVStore(k.orderStoreKey)
	orderIDs := []string{}

	dumpStringHandler := func(key string, it sdk.Iterator, v interface{}) {
		logger.Error(fmt.Sprintf("%s: <%s> -> <%v>", key, types.GetKey(it), v))
	}

	dumpIntHandler := func(key string, it sdk.Iterator, v interface{}) {
		logger.Error(fmt.Sprintf("%s: <%v> -> <%v>", key, common.BytesToInt64(it.Key()[1:]), v))
	}

	dumpKvIntHandler := func(key string, it sdk.Iterator, v interface{}) {
		logger.Error(fmt.Sprintf("%s: <%v> -> <%v>", key, common.BytesToInt64(it.Key()[1:]),
			common.BytesToInt64(it.Value())))
	}

	unmarshalHandler := func(bz []byte, ptr interface{}) {
		k.cdc.MustUnmarshalBinaryBare(bz, ptr)
	}

	unmarshalJSONHanlder := func(bz []byte, ptr interface{}) {
		k.cdc.MustUnmarshalJSON(bz, ptr)
	}

	var order types.Order
	dumpKvs(orderStore, types.OrderKey, "OrderKey", &order, unmarshalHandler, dumpStringHandler)

	var depthBook types.DepthBook
	dumpKvs(orderStore, types.DepthBookKey, "DepthbookKey", &depthBook, unmarshalHandler, dumpStringHandler)

	var price sdk.Dec
	dumpKvs(orderStore, types.PriceKey, "PriceKey", &price, unmarshalHandler, dumpStringHandler)

	dumpKvs(orderStore, types.OrderNumPerBlockKey, "OrderNumPerBlockKey", nil, nil, dumpKvIntHandler)

	dumpKvs(orderStore, types.OrderIDsKey, "OrderIDsKey", &orderIDs, unmarshalJSONHanlder, dumpStringHandler)

	var expireBlockNumbers []int64
	dumpKvs(orderStore, types.ExpireBlockHeightKey, "ExpireBlockHeightKey", &expireBlockNumbers, unmarshalHandler, dumpIntHandler)

	dumpKv(orderStore, logger, types.LastExpiredBlockHeightKey, "LastExpiredBlockHeightKey")
	dumpKv(orderStore, logger, types.OpenOrderNumKey, "OpenOrderNumKey")
	dumpKv(orderStore, logger, types.StoreOrderNumKey, "StoreOrderNumKey")
	dumpKvJSON(orderStore, k, logger, types.RecentlyClosedOrderIDsKey, "RecentlyClosedOrderIDsKey", &orderIDs)
}

func dumpKvs(orderStore sdk.KVStore, k []byte, key string, v interface{},
	unmarshalHandler func([]byte, interface{}),
	dumpHandler func(string, sdk.Iterator, interface{})) {

	orderIter := sdk.KVStorePrefixIterator(orderStore, k)

	for ; orderIter.Valid(); orderIter.Next() {

		if unmarshalHandler != nil {
			unmarshalHandler(orderIter.Value(), v)
		}

		if dumpHandler != nil {
			dumpHandler(key, orderIter, v)
		}
	}
	orderIter.Close()
}

func dumpKvJSON(store sdk.KVStore, keeper Keeper, logger log.Logger, k []byte, key string, v interface{}) {
	bz := store.Get(k)
	if bz != nil {
		keeper.cdc.MustUnmarshalJSON(bz, v)
	}
	logger.Error(fmt.Sprintf("%s: -> <%v>", key, v))
}

func dumpKv(store sdk.KVStore, logger log.Logger, k []byte, key string) {
	bz := store.Get(k)
	var v interface{} = bz
	if bz != nil {
		v = common.BytesToInt64(bz)
	}
	logger.Error(fmt.Sprintf("%s: -> <%v>", key, v))
}
