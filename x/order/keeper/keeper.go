package keeper

import (
	"log"
	"sync"

	"github.com/willf/bitset"

	"github.com/okex/exchain/x/common/monitor"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/params"

	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/order/types"
)

var onStartUp sync.Once

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the various parts of the state machine
type Keeper struct {
	// The reference to the TokenKeeper to modify balances
	tokenKeeper TokenKeeper
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
func NewKeeper(tokenKeeper TokenKeeper, supplyKeeper SupplyKeeper, dexKeeper DexKeeper,
	paramSpace params.Subspace, feeCollectorName string, ordersStoreKey sdk.StoreKey,
	cdc *codec.Codec,
	enableBackend bool, metrics *monitor.OrderMetric) Keeper {

	return Keeper{
		metric: metrics,

		enableBackend:    enableBackend,
		feeCollectorName: feeCollectorName,

		tokenKeeper:  tokenKeeper,
		supplyKeeper: supplyKeeper,
		dexKeeper:    dexKeeper,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),

		orderStoreKey: ordersStoreKey,

		cdc:       cdc,
		cache:     NewCache(),
		diskCache: newDiskCache(),
	}
}

// ResetCache is called in BeginBlock
func (k Keeper) ResetCache(ctx sdk.Context) {

	onStartUp.Do(func() {

		// init depth book map
		depthStore := ctx.KVStore(k.orderStoreKey)
		depthIter := sdk.KVStorePrefixIterator(depthStore, types.DepthBookKey)

		for ; depthIter.Valid(); depthIter.Next() {
			depthBook := &types.DepthBook{}
			k.cdc.MustUnmarshalBinaryBare(depthIter.Value(), depthBook)
			k.diskCache.addDepthBook(types.GetKey(depthIter), depthBook)
		}
		depthIter.Close()

		// init OrderIDs map
		bookStore := ctx.KVStore(k.orderStoreKey)
		bookIter := sdk.KVStorePrefixIterator(bookStore, types.OrderIDsKey)

		for ; bookIter.Valid(); bookIter.Next() {
			var orderIDs []string
			k.cdc.MustUnmarshalJSON(bookIter.Value(), &orderIDs)
			k.diskCache.addOrderIDs(types.GetKey(bookIter), orderIDs)
		}
		bookIter.Close()
	})

	// Reset cache
	k.cache.reset()

	// VERY IMPORTANT: always reset disk cache in BeginBlock
	k.diskCache.reset()
	k.diskCache.setOpenNum(k.GetOpenOrderNum(ctx))
	k.diskCache.setStoreOrderNum(k.GetStoreOrderNum(ctx))

}

// Cache2Disk flushes cached data into KVStore, called in EndBlock
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
}

// OrderOperationMetric records the order information in the depthBook
type OrderOperationMetric struct {
	FullFillNum    int64
	OpenNum        int64
	CancelNum      int64
	ExpireNum      int64
	PartialFillNum int64
}

// GetOperationMetric gets OperationMetric from keeper
func (k Keeper) GetOperationMetric() OrderOperationMetric {
	return OrderOperationMetric{
		FullFillNum:    k.cache.GetFullFillNum(),
		OpenNum:        k.diskCache.getOpenNum(),
		CancelNum:      k.cache.GetCancelNum(),
		ExpireNum:      k.cache.GetExpireNum(),
		PartialFillNum: k.cache.GetPartialFillNum(),
	}
}

// GetCache returns the memoryCache
func (k Keeper) GetCache() *Cache {
	return k.cache
}

// GetDiskCache returns the diskCache
func (k Keeper) GetDiskCache() *DiskCache {
	return k.diskCache
}

// nolint
func (k Keeper) GetTokenKeeper() TokenKeeper {
	return k.tokenKeeper
}

// nolint
func (k Keeper) GetDexKeeper() DexKeeper {
	return k.dexKeeper
}

// GetExpireBlockHeight gets a slice of ExpireBlockHeight from KVStore
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

// GetOrder gets order from KVStore
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

// nolint
func (k Keeper) GetLastPrice(ctx sdk.Context, product string) sdk.Dec {
	// get last price from cache111
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

// GetDepthBookCopy gets depth book copy from cache111, you are supposed to update the Depthbook if you change it
// create if not exist
func (k Keeper) GetDepthBookCopy(product string) *types.DepthBook {
	book := k.diskCache.getDepthBook(product)

	if nil == book {
		return &types.DepthBook{}
		//c.depthBookMap[product] = res // you should do it by yourself
	}

	return book.Copy()
}

// SetDepthBook updates depthBook in diskCache
func (k Keeper) SetDepthBook(product string, book *types.DepthBook) {
	k.diskCache.setDepthBook(product, book)
}

// GetDepthBookFromDB gets depthBook from KVStore
func (k Keeper) GetDepthBookFromDB(ctx sdk.Context, product string) *types.DepthBook {
	store := ctx.KVStore(k.orderStoreKey)
	bookBytes := store.Get(types.GetDepthBookKey(product))
	if bookBytes == nil {
		// Return an empty DepthBook instead of nil
		return &types.DepthBook{}
	}
	depthBook := &types.DepthBook{}
	k.cdc.MustUnmarshalBinaryBare(bookBytes, depthBook)
	return depthBook
}

// GetProductPriceOrderIDs gets OrderIDs from diskCache
func (k Keeper) GetProductPriceOrderIDs(key string) []string {
	if orderIDs := k.diskCache.getOrderIDs(key); orderIDs != nil {
		return orderIDs
	}
	return []string{}
}

// GetProductPriceOrderIDsFromDB gets OrderIDs from KVStore
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

// GetBlockOrderNum gets the num of orders in specific block
func (k Keeper) GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	key := types.GetOrderNumPerBlockKey(blockHeight)
	numBytes := store.Get(key)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

// GetLastExpiredBlockHeight gets LastExpiredBlockHeight from KVStore
// LastExpiredBlockHeight means that the block height of his expired height
// list has been processed by expired recently
func (k Keeper) GetLastExpiredBlockHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.LastExpiredBlockHeightKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

// GetOpenOrderNum gets OpenOrderNum from KVStore
// OpenOrderNum means the number of orders currently in the open state
func (k Keeper) GetOpenOrderNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.OpenOrderNumKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

// StoreOrderNum means the number of orders currently stored
// nolint
func (k Keeper) GetStoreOrderNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.orderStoreKey)
	numBytes := store.Get(types.StoreOrderNumKey)
	if numBytes == nil {
		return 0
	}
	return common.BytesToInt64(numBytes)
}

// GetUpdatedDepthbookKeys gets UpdatedDepthbookKeys from diskCache
func (k Keeper) GetUpdatedDepthbookKeys() []string {
	return k.diskCache.GetUpdatedDepthbookKeys()
}

// GetUpdatedOrderIDs gets UpdatedOrderIDs from memoryCache
func (k Keeper) GetUpdatedOrderIDs() []string {
	return k.cache.getUpdatedOrderIDs()
}

// GetTxHandlerMsgResult: be careful, only call by backend module, other module should never use it!
func (k Keeper) GetTxHandlerMsgResult() []bitset.BitSet {
	if !k.enableBackend {
		return nil
	}
	return k.cache.toggleCopyTxHandlerMsgResult()
}

// nolint
func (k Keeper) addUpdatedOrderID(orderID string) {
	if k.enableBackend {
		k.cache.addUpdatedOrderID(orderID)
	}
}

// GetLastClosedOrderIDs gets closed order ids in last block
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

// nolint
func (k Keeper) GetBlockMatchResult() *types.BlockMatchResult {
	return k.cache.getBlockMatchResult()
}

// nolint
func (k Keeper) SetBlockMatchResult(result *types.BlockMatchResult) {
	if k.enableBackend {
		k.cache.setBlockMatchResult(result)
	}
}

// LockCoins locks coins from the specified address,
func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.SysCoins, lockCoinsType int) error {
	if coins.IsZero() {
		return nil
	}
	return k.tokenKeeper.LockCoins(ctx, addr, coins, lockCoinsType)
}

// nolint
func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.SysCoins, lockCoinsType int) {
	if coins.IsZero() {
		return
	}
	if err := k.tokenKeeper.UnlockCoins(ctx, addr, coins, lockCoinsType); err != nil {
		log.Printf("User(%s) unlock coins(%s) failed\n", addr.String(), coins.String())
	}
}

// BalanceAccount burns the specified coin and obtains another coin
func (k Keeper) BalanceAccount(ctx sdk.Context, addr sdk.AccAddress,
	outputCoins sdk.SysCoins, inputCoins sdk.SysCoins) {

	if err := k.tokenKeeper.BalanceAccount(ctx, addr, outputCoins, inputCoins); err != nil {
		log.Printf("User(%s) burn locked coins(%s) failed\n", addr.String(), outputCoins.String())
	}
}

// nolint
func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.SysCoins {
	return k.tokenKeeper.GetCoins(ctx, addr)
}

// GetProductFeeReceiver gets the fee receiver of specified product from dexKeeper
func (k Keeper) GetProductFeeReceiver(ctx sdk.Context, product string) (sdk.AccAddress, error) {
	tokenPair := k.GetDexKeeper().GetTokenPair(ctx, product)
	if tokenPair == nil {
		return sdk.AccAddress{}, types.ErrTokenPairNotExist(product)
	}

	operator, exists := k.GetDexKeeper().GetOperator(ctx, tokenPair.Owner)
	if exists {
		return operator.HandlingFeeAddress, nil
	}
	return tokenPair.Owner, nil
}

// AddFeeDetail adds detail message of fee to tokenKeeper
func (k Keeper) AddFeeDetail(ctx sdk.Context, from sdk.AccAddress, coins sdk.SysCoins,
	feeType string) {
	if coins.IsZero() {
		return
	}
	k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType, "")
}

// SendFeesToProductOwner sends fees from the specified address to productOwner
func (k Keeper) SendFeesToProductOwner(ctx sdk.Context, coins sdk.SysCoins, from sdk.AccAddress,
	feeType string, product string) (feeReceiver string, err error) {
	if coins.IsZero() {
		return "", nil
	}
	to, err := k.GetProductFeeReceiver(ctx, product)
	if err != nil {
		return "", err
	}
	k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType, "")
	if err := k.tokenKeeper.SendCoinsFromAccountToAccount(ctx, from, to, coins); err != nil {
		log.Printf("Send fee(%s) to address(%s) failed\n", coins.String(), to.String())
		return "", types.ErrSendCoinsFailed(coins.String(), to.String())
	}
	return to.String(), nil
}

// AddCollectedFees adds fee to the feePool
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.SysCoins, from sdk.AccAddress,
	feeType string, hasFeeDetail bool) error {
	if coins.IsZero() {
		return nil
	}
	if hasFeeDetail {
		k.tokenKeeper.AddFeeDetail(ctx, from.String(), coins, feeType, k.feeCollectorName)
	}
	baseCoins := coins
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, k.feeCollectorName, baseCoins)
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) *types.Params {
	var param types.Params
	k.paramSpace.GetParamSet(ctx, &param)
	return &param
}

// SetParams sets inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	k.paramSpace.SetParamSet(ctx, params)
}

// nolint
func (k Keeper) GetMetric() *monitor.OrderMetric {
	return k.metric
}

// nolint
func (k Keeper) SetMetric() {
	k.metric.FullFilledNum.Set(float64(k.cache.fullFillNum))
	k.metric.PendingNum.Set(float64(k.diskCache.openNum))
	k.metric.CanceledNum.Set(float64(k.cache.cancelNum))
	k.metric.ExpiredNum.Set(float64(k.cache.expireNum))
	k.metric.PartialFilledNum.Set(float64(k.cache.partialFillNum))
}

// GetBestBidAndAsk gets the highest bidPrice and the lowest askPrice from depthBook
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

// RemoveOrderFromDepthBook removes order from depthBook, and updates cancelNum, expireNum, updatedOrderIDs from cache111
func (k Keeper) RemoveOrderFromDepthBook(order *types.Order, feeType string) {
	k.addUpdatedOrderID(order.OrderID)
	if feeType == types.FeeTypeOrderCancel {
		k.cache.IncreaseCancelNum()
	} else if feeType == types.FeeTypeOrderExpire {
		k.cache.IncreaseExpireNum()
	}

	k.diskCache.removeOrder(order)
}

// nolint
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

// nolint
func (k Keeper) InsertOrderIntoDepthBook(order *types.Order) {
	k.diskCache.insertOrder(order)
}

// FilterDelistedProducts deletes non-existent products from the specified products
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

// nolint
func (k Keeper) AddTxHandlerMsgResult(resultSet bitset.BitSet) {
	if k.enableBackend {
		k.cache.addTxHandlerMsgResult(resultSet)
	}
}
