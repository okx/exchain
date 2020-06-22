package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/types"
)

// OrderIDsMap stores orderIDSlice with map.
// <product:price:side> -> <orderIDs>
type OrderIDsMap struct {
	Data         map[string][]string
	updatedItems map[string]struct{}
}

// DepthBookMap stores depthBook with map.
// <product> -> <depthBook>
type DepthBookMap struct {
	data         map[string]*types.DepthBook
	updatedItems map[string]struct{}
	newItems     map[string]struct{}
}

// DiskCache stores cache that will persist to disk at endBlock.
type DiskCache struct {
	depthBookMap *DepthBookMap
	orderIDsMap  *OrderIDsMap
	priceMap     map[string]sdk.Dec

	storeOrderNum  int64 // current stored order num
	openNum        int64 // current open orders num
	closedOrderIDs []string
}

func newDiskCache() *DiskCache {
	return &DiskCache{
		closedOrderIDs: []string{},
		orderIDsMap:    &OrderIDsMap{make(map[string][]string), make(map[string]struct{})},
		priceMap:       make(map[string]sdk.Dec),
		depthBookMap: &DepthBookMap{make(map[string]*types.DepthBook), make(map[string]struct{}),
			make(map[string]struct{})},
	}
}

// reset is invoked in begin block
func (c *DiskCache) reset() {
	c.closedOrderIDs = []string{}
	c.orderIDsMap.updatedItems = make(map[string]struct{})
	c.depthBookMap.updatedItems = make(map[string]struct{})
	c.depthBookMap.newItems = make(map[string]struct{})
}

// nolint
func (c *DiskCache) GetClosedOrderIDs() []string {
	return c.closedOrderIDs
}

func (c *DiskCache) setLastPrice(product string, price sdk.Dec) {
	c.priceMap[product] = price
}

func (c *DiskCache) getLastPrice(product string) sdk.Dec {
	if price, ok := c.priceMap[product]; ok {
		return price
	}
	return sdk.ZeroDec()
}

// GetOrderIDsMapCopy returns a new copy of OrderIDsMap
func (c *DiskCache) GetOrderIDsMapCopy() *OrderIDsMap {
	if c.orderIDsMap == nil {
		return nil
	}
	ret := make(map[string][]string)
	for k, v := range c.orderIDsMap.Data {
		if len(v) == 0 {
			ret[k] = []string{}
		}
		ret[k] = append(ret[k], v...)
	}
	return &OrderIDsMap{Data: ret}
}

func (c *DiskCache) getOrderIDs(key string) []string {
	if c.orderIDsMap == nil {
		return nil
	}

	return c.orderIDsMap.Data[key]
}

func (c *DiskCache) setStoreOrderNum(num int64) {
	c.storeOrderNum = num
}

// nolint
func (c *DiskCache) DecreaseStoreOrderNum(num int64) int64 {
	c.storeOrderNum -= num
	return c.storeOrderNum
}

func (c *DiskCache) setOpenNum(num int64) {
	c.openNum = num
}

func (c *DiskCache) getOpenNum() int64 {
	return c.openNum
}

func (c *DiskCache) addOrderIDs(key string, orderIDs []string) {
	c.orderIDsMap.Data[key] = orderIDs
}

func (c *DiskCache) addDepthBook(product string, book *types.DepthBook) {
	if book == nil {
		panic("failed. a nil pointer appears")
	}

	c.depthBookMap.data[product] = book
}

// setOrderIDs updates or removes unfilled order ids
func (c *DiskCache) setOrderIDs(key string, orderIDs []string) {

	if len(orderIDs) == 0 {
		// remove empty element immediately, not do it by the end of endblock
		delete(c.orderIDsMap.Data, key)
	} else {
		c.orderIDsMap.Data[key] = orderIDs
	}
	c.orderIDsMap.updatedItems[key] = struct{}{}
}

// setDepthBook updates or removes a depth book
func (c *DiskCache) setDepthBook(product string, book *types.DepthBook) {
	if book != nil && len(book.Items) > 0 {
		c.depthBookMap.data[product] = book
	} else {
		delete(c.depthBookMap.data, product)
	}
	c.depthBookMap.updatedItems[product] = struct{}{}
}

// UpdatedOrderIDKeys
// nolint
func (c *DiskCache) GetUpdatedOrderIDKeys() []string {
	updatedKeys := make([]string, 0, len(c.orderIDsMap.updatedItems))
	for key := range c.orderIDsMap.updatedItems {
		updatedKeys = append(updatedKeys, key)
	}
	sort.Strings(updatedKeys)
	return updatedKeys
}

func (c *DiskCache) getDepthBook(product string) *types.DepthBook {
	res := c.depthBookMap.data[product]
	return res
}

func (c *DiskCache) getProductsFromDepthBookMap() []string {
	products := make([]string, 0, len(c.depthBookMap.data))
	for product := range c.depthBookMap.data {
		products = append(products, product)
	}
	return products
}

// GetUpdatedDepthbookKeys returns a new copy of UpdatedDepthbookKeys
func (c *DiskCache) GetUpdatedDepthbookKeys() []string {
	updatedKeys := make([]string, 0, len(c.depthBookMap.updatedItems))
	for key := range c.depthBookMap.updatedItems {
		updatedKeys = append(updatedKeys, key)
	}
	sort.Strings(updatedKeys)
	return updatedKeys
}

// GetNewDepthbookKeys returns a new copy of NewDepthbookKeys
func (c *DiskCache) GetNewDepthbookKeys() []string {
	newAddKeys := make([]string, 0, len(c.depthBookMap.newItems))
	for key := range c.depthBookMap.newItems {
		newAddKeys = append(newAddKeys, key)
	}
	return newAddKeys
}

// insertOrder inserts a new order into orderIDsMap
func (c *DiskCache) insertOrder(order *types.Order) {
	// 1. update depthBookMap
	depthBook, ok := c.depthBookMap.data[order.Product]
	if !ok {
		depthBook = &types.DepthBook{}
		c.depthBookMap.data[order.Product] = depthBook
	}
	depthBook.InsertOrder(order)
	c.depthBookMap.updatedItems[order.Product] = struct{}{}
	c.depthBookMap.newItems[order.Product] = struct{}{}

	// 2. update orderIDsMap
	orderIDsMap := c.orderIDsMap
	key := types.FormatOrderIDsKey(order.Product, order.Price, order.Side)
	orderIDs, ok := orderIDsMap.Data[key]
	if !ok {
		orderIDs = []string{}
	}
	orderIDs = append(orderIDs, order.OrderID)
	orderIDsMap.Data[key] = orderIDs
	c.orderIDsMap.updatedItems[key] = struct{}{}

	c.openNum++
	c.storeOrderNum++
}

func (c *DiskCache) closeOrder(orderID string) {
	c.closedOrderIDs = append(c.closedOrderIDs, orderID)
	c.openNum--
}

// remove an order from orderIDsMap when order cancelled/expired
func (c *DiskCache) removeOrder(order *types.Order) {

	// update depth book map
	depthBook := c.getDepthBook(order.Product)
	if depthBook != nil {
		depthBook.RemoveOrder(order)
		c.setDepthBook(order.Product, depthBook)
	}

	// update order id map
	orderIDsMap := c.orderIDsMap
	key := types.FormatOrderIDsKey(order.Product, order.Price, order.Side)
	orderIDs := orderIDsMap.Data[key]
	orderIDsLen := len(orderIDs)
	for i := 0; i < orderIDsLen; i++ {
		if orderIDs[i] == order.OrderID {
			orderIDs = append(orderIDs[:i], orderIDs[i+1:]...)
			c.setOrderIDs(key, orderIDs)
			break
		}
	}

	c.closeOrder(order.OrderID)
}
