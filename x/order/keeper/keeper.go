package keeper

import (
	"log"

	"github.com/okex/okchain/x/common/monitor"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/params"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/order/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the various parts of the state machine
type Keeper struct {
	// The reference to the TokenKeeper to modify balances
	tokenKeeper TokenKeeper
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper
	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace

	dexKeeper DexKeeper

	supplyKeeper     SupplyKeeper
	feeCollectorName string

	// Unexposed key to access name store from sdk.Context
	orderStoreKey sdk.StoreKey

	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	enableBackend bool         // whether open backend plugin
	metric        *monitor.OrderMetric

	// cache data in memory to avoid marshal/unmarshal too frequently
	// reset cache data in BeginBlock
	cache     *Cache
	diskCache *DiskCache
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(tokenKeeper TokenKeeper, supplyKeeper SupplyKeeper, paramsKeeper params.Keeper, dexKeeper DexKeeper,
	paramSpace params.Subspace, feeCollectorName string, ordersStoreKey sdk.StoreKey,
	cdc *codec.Codec,
	enableBackend bool, metrics *monitor.OrderMetric) Keeper {

	return Keeper{
		metric: metrics,

		enableBackend:    enableBackend,
		feeCollectorName: feeCollectorName,

		tokenKeeper:  tokenKeeper,
		supplyKeeper: supplyKeeper,
		paramsKeeper: paramsKeeper,
		dexKeeper:    dexKeeper,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),

		orderStoreKey: ordersStoreKey,

		cdc:       cdc,
		cache:     NewCache(),
		diskCache: newDiskCache(),
	}
}

// Reset cache, called in BeginBlock
func (k Keeper) ResetCache(ctx sdk.Context) {
	// Reset cache
	k.cache.reset()

	k.diskCache.reset()
	k.diskCache.setOpenNum(k.GetOpenOrderNum(ctx))
	k.diskCache.setStoreOrderNum(k.GetStoreOrderNum(ctx))

	// init depth book & items cache
	if len(k.diskCache.depthBookMap.data) == 0 {
		depthStore := ctx.KVStore(k.orderStoreKey)
		depthIter := sdk.KVStorePrefixIterator(depthStore, types.DepthbookKey)

		for ; depthIter.Valid(); depthIter.Next() {
			depthBook := &types.DepthBook{}
			k.cdc.MustUnmarshalBinaryBare(depthIter.Value(), depthBook)
			k.SetDepthBook(types.GetKey(depthIter), depthBook)
		}
		depthIter.Close()
	}
	if len(k.diskCache.orderIDsMap.Data) == 0 {
		bookStore := ctx.KVStore(k.orderStoreKey)
		bookIter := sdk.KVStorePrefixIterator(bookStore, types.OrderIDsKey)

		for ; bookIter.Valid(); bookIter.Next() {
			var orderIDs []string
			k.cdc.MustUnmarshalJSON(bookIter.Value(), &orderIDs)
			k.SetOrderIDs(types.GetKey(bookIter), orderIDs) // startup
		}
		bookIter.Close()
	}
}

// Flush cached data into KVStore, called in EndBlock
func (k Keeper) Cache2Disk(ctx sdk.Context) {

	closedOrderIDs := k.diskCache.GetClosedOrderIDs()

	k.SetLastClosedOrderIDs(ctx, closedOrderIDs)
	k.setOpenOrderNum(ctx, k.diskCache.openNum)
	k.setStoreOrderNum(ctx, k.diskCache.storeOrderNum)

	// update depth book to KVStore
	updatedBookKeys := k.diskCache.GetUpdatedDepthbookKeys()
	for _, key := range updatedBookKeys {
		k.StoreDepthBook(ctx, key, k.diskCache.getDepthBook(key))
	}

	updatedItemKeys := k.diskCache.GetUpdatedOrderIDKeys()
	for _, key := range updatedItemKeys {
		k.StoreOrderIDsMap(ctx, key, k.diskCache.getOrderIDs(key))
	}
	k.diskCache.flush()
}

type OrderOperationMetric struct {
	FullFillNum    int64
	OpenNum        int64
	CancelNum      int64
	ExpireNum      int64
	PartialFillNum int64
}

func (k Keeper) GetOperationMetric() OrderOperationMetric {
	return OrderOperationMetric{
		FullFillNum:    k.cache.GetFullFillNum(),
		OpenNum:        k.diskCache.getOpenNum(),
		CancelNum:      k.cache.GetCancelNum(),
		ExpireNum:      k.cache.GetExpireNum(),
		PartialFillNum: k.cache.GetPartialFillNum(),
	}
}

func (k Keeper) GetCache() *Cache {
	return k.cache
}

func (k Keeper) GetDiskCache() *DiskCache {
	return k.diskCache
}

func (k Keeper) GetTokenKeeper() TokenKeeper {
	return k.tokenKeeper
}

func (k Keeper) GetDexKeeper() DexKeeper {
	return k.dexKeeper
}

func (k Keeper) GetExpireBlockHeight(ctx sdk.Context, blockHeight int64) []int64 {
	store := ctx.KVStore(k.orderStoreKey)
	orderInfo := store.Get(types.GetExpireBlockHeightKey(blockHeight))
	if orderInfo == nil {
		return []int64{}
	}
	var expireBlockNumbers []int64
	k.cdc.MustUnmarshalBinaryBare(orderInfo, &expireBlockNumbers)
	return expireBlockNumbers
}

func (k Keeper) GetOrder(ctx sdk.Context, orderID string) *types.Order {
	store := ctx.KVStore(k.orderStoreKey)
	orderInfo := store.Get(types.GetOrderKey(orderID))
	if orderInfo == nil {
		return nil
	}
	order := &types.Order{}
	k.cdc.MustUnmarshalBinaryBare(orderInfo, order)
	return order
}

func (k Keeper) GetLastPrice(ctx sdk.Context, product string) sdk.Dec {
	// get last price from cache
	price := k.diskCache.getLastPrice(product)
	if price.IsPositive() {
		return price
	}

	// get last price from KVStore
	store := ctx.KVStore(k.orderStoreKey)
	priceBytes := store.Get(types.GetPriceKey(product))
	if priceBytes == nil {
		// If last price does not exist, set the init price of token pair as last price
		tokenPair := k.dexKeeper.GetTokenPair(ctx, product)

		if tokenPair == nil {
			return sdk.ZeroDec()
		}
		k.SetLastPrice(ctx, product, tokenPair.InitPrice)
		return tokenPair.InitPrice
	}

	k.cdc.MustUnmarshalBinaryBare(priceBytes, &price)
	return price
}

// Get depth book copy from cache, you are supposed to update the Depthbook if you change it
// create if not exist
func (k Keeper) GetDepthBookCopy(product string) *types.DepthBook {
	book := k.diskCache.getDepthBook(product)

	if nil == book {
		return &types.DepthBook{}
		//c.depthBookMap[product] = res // you should do it by yourself
	}

	return book.Copy()
}

func (k Keeper) SetDepthBook(product string, book *types.DepthBook) {
	k.diskCache.setDepthBook(product, book)
}

func (k Keeper) GetDepthBookFromDB(ctx sdk.Context, product string) *types.DepthBook {
	store := ctx.KVStore(k.orderStoreKey)
	bookBytes := store.Get(types.GetDepthbookKey(product))
	if bookBytes == nil {
		// Return an empty DepthBook instead of nil
		return &types.DepthBook{}
	}
	depthBook := &types.DepthBook{}
	k.cdc.MustUnmarshalBinaryBare(bookBytes, depthBook)
	return depthBook
}

func (k Keeper) GetProductPriceOrderIDs(key string) []string {
	if orderIDs := k.diskCache.getOrderIDs(key); orderIDs != nil {
		return orderIDs
	}
	return []string{}
}

func (k Keeper) GetProductPriceOrderIDsFromDB(ctx sdk.Context, key string) []string {
	store := ctx.KVStore(k.orderStoreKey)
	idsBytes := store.Get(types.GetOrderIDsKey(key))
	if idsBytes == nil {
		return nil
	}
	var orderIDs []string
	k.cdc.MustUnmarshalJSON(idsBytes, &orderIDs)
	return orderIDs
}

// get the num of orders in specific block
func (k Keeper) GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetOrderNumPerBlockKey(blockHeight)
	numBytes := store.Get(key)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

func (k Keeper) GetLastExpiredBlockHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.LastExpiredBlockHeightKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

func (k Keeper) GetOpenOrderNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.OpenOrderNumKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

func (k Keeper) GetStoreOrderNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.StoreOrderNumKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

func (k Keeper) GetUpdatedDepthbookKeys() []string {
	return k.diskCache.GetUpdatedDepthbookKeys()
}

func (k Keeper) GetUpdatedOrderIDs() []string {
	return k.cache.getUpdatedOrderIDs()
}

func (k Keeper) addUpdatedOrderID(orderID string) {
	if k.enableBackend {
		k.cache.addUpdatedOrderID(orderID)
	}
}

// get closed order ids in last block
func (k Keeper) GetLastClosedOrderIDs(ctx sdk.Context) []string {
	store := ctx.KVStore(k.orderStoreKey)
	bz := store.Get(types.RecentlyClosedOrderIDsKey)
	orderIDs := []string{}
	if bz == nil {
		return orderIDs
	}
	k.cdc.MustUnmarshalJSON(bz, &orderIDs)
	return orderIDs
}

func (k Keeper) GetBlockMatchResult() *types.BlockMatchResult {
	return k.cache.getBlockMatchResult()
}

func (k Keeper) SetBlockMatchResult(result *types.BlockMatchResult) {
	if k.enableBackend {
		k.cache.setBlockMatchResult(result)
	}
}

// use TokenKeeper
func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) error {
	return k.tokenKeeper.LockCoins(ctx, addr, coins, lockCoinsType)
}

func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) {
	if err := k.tokenKeeper.UnlockCoins(ctx, addr, coins, lockCoinsType); err != nil {
		log.Printf("User(%s) unlock coins(%s) failed\n", addr.String(), coins.String())
	}
}

func (k Keeper) BalanceAccount(ctx sdk.Context, addr sdk.AccAddress,
	outputCoins sdk.DecCoins, inputCoins sdk.DecCoins) {

	if err := k.tokenKeeper.BalanceAccount(ctx, addr, outputCoins, inputCoins); err != nil {
		log.Printf("User(%s) burn locked coins(%s) failed\n", addr.String(), outputCoins.String())
	}
}

func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.DecCoins {
	return k.tokenKeeper.GetCoins(ctx, addr)
}

func (k Keeper) GetProductOwner(ctx sdk.Context, product string) sdk.AccAddress {
	tokenPair := k.GetDexKeeper().GetTokenPair(ctx, product)
	return tokenPair.Owner
}

func (k Keeper) AddFeeDetail(ctx sdk.Context, from sdk.AccAddress, coins sdk.DecCoins,
	feeType string) {
	k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType)
}

func (k Keeper) SendFeesToProductOwner(ctx sdk.Context, coins sdk.DecCoins, from sdk.AccAddress,
	feeType string, product string) error {
	if coins.IsZero() {
		return nil
	}
	to := k.GetProductOwner(ctx, product)
	k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType)
	if err := k.tokenKeeper.SendCoinsFromAccountToAccount(ctx, from, to, coins); err != nil {
		log.Printf("Send fee(%s) to address(%s) failed\n", coins.String(), to.String())
		return err
	}
	return nil
}

// use feeCollectionKeeper
// AddCollectedFees - add to the fee pool
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.DecCoins, from sdk.AccAddress,
	feeType string, hasFeeDetail bool) error {
	if coins.IsZero() {
		return nil
	}
	if hasFeeDetail {
		k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType)
	}
	baseCoins := coins
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, k.feeCollectorName, baseCoins)
}

// get inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) *types.Params {
	// get params from cache
	cacheParams := k.cache.GetParams()
	if cacheParams != nil {
		return cacheParams
	}

	// if param not stored in cache, get param from KVStore and cache it
	var param types.Params
	k.paramSpace.GetParamSet(ctx, &param)
	k.cache.SetParams(&param)
	return &param
}

// set inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	k.paramSpace.SetParamSet(ctx, params)
	k.cache.SetParams(params)
}

func (k Keeper) GetMetric() *monitor.OrderMetric {
	return k.metric
}

func (k Keeper) SetMetric() {
	k.metric.FullFilledNum.Set(float64(k.cache.fullFillNum))
	k.metric.PendingNum.Set(float64(k.diskCache.openNum))
	k.metric.CanceledNum.Set(float64(k.cache.cancelNum))
	k.metric.ExpiredNum.Set(float64(k.cache.expireNum))
	k.metric.PartialFilledNum.Set(float64(k.cache.partialFillNum))
}

func (k Keeper) GetBestBidAndAsk(ctx sdk.Context, product string) (sdk.Dec, sdk.Dec) {
	bestBid := sdk.ZeroDec()
	bestAsk := sdk.ZeroDec()
	depthBook := k.GetDepthBookFromDB(ctx, product)

	for _, item := range depthBook.Items {
		if item.BuyQuantity.IsPositive() {
			if item.Price.GT(bestBid) {
				bestBid = item.Price
			}
		}
		if item.SellQuantity.IsPositive() {
			if item.Price.LT(bestAsk) {
				bestAsk = item.Price
			}
		}
	}
	return bestBid, bestAsk
}

func (k Keeper) RemoveOrderFromDepthBook(order *types.Order, feeType string) {
	k.addUpdatedOrderID(order.OrderID)
	if feeType == types.FeeTypeOrderCancel {
		k.cache.IncreaseCancelNum()
	} else if feeType == types.FeeTypeOrderExpire {
		k.cache.IncreaseExpireNum()
	}

	k.diskCache.removeOrder(order)
}

func (k Keeper) UpdateOrder(order *types.Order, ctx sdk.Context) {
	// update order to keeper
	k.SetOrder(ctx, order.OrderID, order)
	// record updated orderID
	k.addUpdatedOrderID(order.OrderID)
	if order.Status == types.OrderStatusFilled {
		k.diskCache.closeOrder(order.OrderID)
		k.cache.IncreaseFullFillNum()
	} else {
		k.cache.IncreasePartialFillNum()
	}
}

func (k Keeper) InsertOrderIntoDepthBook(order *types.Order) {
	k.diskCache.insertOrder(order)
}

func (k Keeper) FilterDelistedProducts(ctx sdk.Context, products []string) []string {
	var cleanProducts []string
	for _, product := range products {
		tokenPair := k.dexKeeper.GetTokenPair(ctx, product)
		if tokenPair != nil {
			cleanProducts = append(cleanProducts, product)
		}
	}
	return cleanProducts
}
