package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/order/types"
)

// ===============================================
// 1.
func (k Keeper) SetBlockOrderNum(ctx sdk.Context, blockHeight int64, orderNum int64) {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetOrderNumPerBlockKey(blockHeight)
	store.Set(key, common.Int64ToBytes(orderNum))
}

func (k Keeper) DropBlockOrderNum(ctx sdk.Context, blockHeight int64) {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetOrderNumPerBlockKey(blockHeight)
	store.Delete(key)
}

// ===============================================
// 2.
func (k Keeper) SetExpireBlockHeight(ctx sdk.Context, blockHeight int64, expireBlockHeight []int64) {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetExpireBlockHeightKey(blockHeight)
	store.Set(key, k.cdc.MustMarshalBinaryBare(expireBlockHeight))
}

func (k Keeper) DropExpireBlockHeight(ctx sdk.Context, blockHeight int64) {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetExpireBlockHeightKey(blockHeight)
	store.Delete(key)
}

// ===============================================
// 3.
// NewOrder - place an order
func (k Keeper) SetOrder(ctx sdk.Context, orderID string, order *types.Order) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Set(types.GetOrderKey(orderID), k.cdc.MustMarshalBinaryBare(order))
}

func (k Keeper) DropOrder(ctx sdk.Context, orderID string) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Delete(types.GetOrderKey(orderID))
}

// ===============================================
// 4.
func (k Keeper) StoreDepthBook(ctx sdk.Context, product string, depthBook *types.DepthBook) {
	store := ctx.KVStore(k.orderStoreKey)
	if depthBook == nil || len(depthBook.Items) == 0 {
		store.Delete(types.GetDepthbookKey(product))
	} else {
		store.Set(types.GetDepthbookKey(product), k.cdc.MustMarshalBinaryBare(depthBook)) //
	}
}

// ===============================================
// 5.
func (k Keeper) SetLastPrice(ctx sdk.Context, product string, price sdk.Dec) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Set(types.GetPriceKey(product), k.cdc.MustMarshalBinaryBare(price))
	k.diskCache.setLastPrice(product, price)
}

// ===============================================
// 6.
func (k Keeper) StoreOrderIDsMap(ctx sdk.Context, key string, orderIDs []string) {
	store := ctx.KVStore(k.orderStoreKey)
	if len(orderIDs) == 0 {
		store.Delete(types.GetOrderIDsKey(key))
	} else {
		store.Set(types.GetOrderIDsKey(key), k.cdc.MustMarshalJSON(orderIDs)) //StoreOrderIDsMap
	}
}

// ===============================================
// 7.

//func (k Keeper) updateStoreProductLockMap(ctx sdk.Context, lockMap *types.ProductLockMap) {
//	store := ctx.KVStore(k.otherStoreKey)
//	bz, _ := json.Marshal(lockMap)
//	store.Set([]byte("productLockMap"), bz)
//}
//
//

// ===============================================
// 8.

func (k Keeper) SetLastExpiredBlockHeight(ctx sdk.Context, expiredBlockHeight int64) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Set(types.LastExpiredBlockHeightKey, common.Int64ToBytes(expiredBlockHeight)) //lastExpiredBlockHeight
}

// ===============================================
// 9.

func (k Keeper) setOpenOrderNum(ctx sdk.Context, orderNum int64) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Set(types.OpenOrderNumKey, common.Int64ToBytes(orderNum)) //openOrderNum
}

// ===============================================
// 10.
func (k Keeper) setStoreOrderNum(ctx sdk.Context, orderNum int64) {
	store := ctx.KVStore(k.orderStoreKey)
	store.Set(types.StoreOrderNumKey, common.Int64ToBytes(orderNum)) //StoreOrderNum
}

// ===============================================
// 11.
// set closed order ids in this block
func (k Keeper) SetLastClosedOrderIDs(ctx sdk.Context, orderIDs []string) {
	store := ctx.KVStore(k.orderStoreKey)
	if len(orderIDs) == 0 {
		store.Delete(types.RecentlyClosedOrderIDsKey)
	}
	store.Set(types.RecentlyClosedOrderIDsKey, k.cdc.MustMarshalJSON(orderIDs)) //recentlyClosedOrderIDs
}

func (k Keeper) SetOrderIDs(key string, orderIDs []string) {
	k.diskCache.setOrderIDs(key, orderIDs)
}

func (k Keeper) GetProductsFromDepthBookMap() []string {
	return k.diskCache.getProductsFromDepthBookMap()
}
