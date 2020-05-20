package keeper

import (
	"encoding/json"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/types"
)

// OrderIDsMap stores orderIDSlice with map.
// <product:price:side> -> <orderIDs>
type OrderIDsMap struct {
	Data         map[string][]string `json:"data"`
	UpdatedItems map[string]struct{} `json:"updated_items"`
}

// DepthBookMap stores depthBook with map.
// <product> -> <depthBook>
type DepthBookMap struct {
	Data         map[string]*types.DepthBook `json:"data"`
	UpdatedItems map[string]struct{}         `json:"updated_items"`
	NewItems     map[string]struct{}         `json:"new_items"`
}

// DiskCache stores cache that will persist to disk at endBlock.
type DiskCache struct {
	DepthBookMap *DepthBookMap      `json:"depth_book_map"`
	OrderIDsMap  *OrderIDsMap       `json:"order_ids_map"`
	PriceMap     map[string]sdk.Dec `json:"price_map"`

	StoreOrderNum  int64    `json:"store_order_number"` // current stored order num
	OpenNum        int64    `json:"open_number"`        // current open orders num
	ClosedOrderIDs []string `json:"closed_order_ids"`
}

func newDiskCache() *DiskCache {
	return &DiskCache{
		ClosedOrderIDs: []string{},
		OrderIDsMap:    &OrderIDsMap{make(map[string][]string), make(map[string]struct{})},
		PriceMap:       make(map[string]sdk.Dec),
		DepthBookMap: &DepthBookMap{make(map[string]*types.DepthBook), make(map[string]struct{}),
			make(map[string]struct{})},
	}
}

// reset is invoked in begin block
func (c *DiskCache) reset() {
	c.ClosedOrderIDs = []string{}
	c.OrderIDsMap.UpdatedItems = make(map[string]struct{})
	c.DepthBookMap.UpdatedItems = make(map[string]struct{})
	c.DepthBookMap.NewItems = make(map[string]struct{})
}

// nolint
func (c *DiskCache) GetClosedOrderIDs() []string {
	return c.ClosedOrderIDs
}

func (c *DiskCache) setLastPrice(product string, price sdk.Dec) {
	c.PriceMap[product] = price
}

func (c *DiskCache) getLastPrice(product string) sdk.Dec {
	if price, ok := c.PriceMap[product]; ok {
		return price
	}
	return sdk.ZeroDec()
}

// GetOrderIDsMapCopy returns a new copy of OrderIDsMap
func (c *DiskCache) GetOrderIDsMapCopy() *OrderIDsMap {
	if c.OrderIDsMap == nil {
		return nil
	}
	ret := make(map[string][]string)
	for k, v := range c.OrderIDsMap.Data {
		if len(v) == 0 {
			ret[k] = []string{}
		}
		ret[k] = append(ret[k], v...)
	}
	return &OrderIDsMap{Data: ret}
}

func (c *DiskCache) getOrderIDs(key string) []string {
	return c.OrderIDsMap.Data[key]
}

func (c *DiskCache) setStoreOrderNum(num int64) {
	c.StoreOrderNum = num
}

// nolint
func (c *DiskCache) DecreaseStoreOrderNum(num int64) int64 {
	c.StoreOrderNum -= num
	return c.StoreOrderNum
}

func (c *DiskCache) setOpenNum(num int64) {
	c.OpenNum = num
}

func (c *DiskCache) getOpenNum() int64 {
	return c.OpenNum
}

func (c *DiskCache) addOrderIDs(key string, orderIDs []string) {
	c.OrderIDsMap.Data[key] = orderIDs
}

func (c *DiskCache) addDepthBook(product string, book *types.DepthBook) {
	c.DepthBookMap.Data[product] = book
}

// setOrderIDs updates or removes unfilled order ids
func (c *DiskCache) setOrderIDs(key string, orderIDs []string) {

	if len(orderIDs) == 0 {
		// remove empty element immediately, not do it by the end of endblock
		delete(c.OrderIDsMap.Data, key)
	} else {
		c.OrderIDsMap.Data[key] = orderIDs
	}
	c.OrderIDsMap.UpdatedItems[key] = struct{}{}
}

// setDepthBook updates or removes a depth book
func (c *DiskCache) setDepthBook(product string, book *types.DepthBook) {
	if book != nil && len(book.Items) > 0 {
		c.DepthBookMap.Data[product] = book
	} else {
		delete(c.DepthBookMap.Data, product)
	}
	c.DepthBookMap.UpdatedItems[product] = struct{}{}
}

// UpdatedOrderIDKeys
// nolint
func (c *DiskCache) GetUpdatedOrderIDKeys() []string {
	updatedKeys := make([]string, 0, len(c.OrderIDsMap.UpdatedItems))
	for key := range c.OrderIDsMap.UpdatedItems {
		updatedKeys = append(updatedKeys, key)
	}
	sort.Strings(updatedKeys)
	return updatedKeys
}

func (c *DiskCache) getDepthBook(product string) *types.DepthBook {
	res := c.DepthBookMap.Data[product]
	return res
}

func (c *DiskCache) getProductsFromDepthBookMap() []string {
	products := make([]string, 0, len(c.DepthBookMap.Data))
	for product := range c.DepthBookMap.Data {
		products = append(products, product)
	}
	return products
}

// GetUpdatedDepthbookKeys returns a new copy of UpdatedDepthbookKeys
func (c *DiskCache) GetUpdatedDepthbookKeys() []string {
	updatedKeys := make([]string, 0, len(c.DepthBookMap.UpdatedItems))
	for key := range c.DepthBookMap.UpdatedItems {
		updatedKeys = append(updatedKeys, key)
	}
	sort.Strings(updatedKeys)
	return updatedKeys
}

// GetNewDepthbookKeys returns a new copy of NewDepthbookKeys
func (c *DiskCache) GetNewDepthbookKeys() []string {
	newAddKeys := make([]string, 0, len(c.DepthBookMap.NewItems))
	for key := range c.DepthBookMap.NewItems {
		newAddKeys = append(newAddKeys, key)
	}
	return newAddKeys
}

// insertOrder inserts a new order into orderIDsMap
func (c *DiskCache) insertOrder(order *types.Order) {
	// 1. update depthBookMap
	depthBook, ok := c.DepthBookMap.Data[order.Product]
	if !ok {
		depthBook = &types.DepthBook{}
		c.DepthBookMap.Data[order.Product] = depthBook
	}
	depthBook.InsertOrder(order)
	c.DepthBookMap.UpdatedItems[order.Product] = struct{}{}
	c.DepthBookMap.NewItems[order.Product] = struct{}{}

	// 2. update orderIDsMap
	orderIDsMap := c.OrderIDsMap
	key := types.FormatOrderIDsKey(order.Product, order.Price, order.Side)
	orderIDs, ok := orderIDsMap.Data[key]
	if !ok {
		orderIDs = []string{}
	}
	orderIDs = append(orderIDs, order.OrderID)
	orderIDsMap.Data[key] = orderIDs
	c.OrderIDsMap.UpdatedItems[key] = struct{}{}

	c.OpenNum++
	c.StoreOrderNum++
}

func (c *DiskCache) closeOrder(orderID string) {
	c.ClosedOrderIDs = append(c.ClosedOrderIDs, orderID)
	c.OpenNum--
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
	orderIDsMap := c.OrderIDsMap
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

// nolint
func (c *DiskCache) Clone() *DiskCache {
	cache := &DiskCache{}
	bytes, _ := json.Marshal(c)
	err := json.Unmarshal(bytes, cache)
	if err != nil {
		return c.DepthCopy()
	}

	return cache
}

// nolint
func (c *DiskCache) DepthCopy() *DiskCache {
	cache := DiskCache{
		DepthBookMap:   nil,
		OrderIDsMap:    nil,
		PriceMap:       nil,
		StoreOrderNum:  c.StoreOrderNum,
		OpenNum:        c.OpenNum,
		ClosedOrderIDs: nil,
	}

	if c.DepthBookMap != nil {
		cpData := make(map[string]*types.DepthBook)
		for k, v := range c.DepthBookMap.Data {
			cpItems := make([]types.DepthBookItem, 0, len(v.Items))
			cpItems = append(cpItems, v.Items...)
			cpData[k] = &types.DepthBook{Items: cpItems}
		}

		cpUpdatedItems := make(map[string]struct{})
		for k, _ := range c.DepthBookMap.UpdatedItems {
			cpUpdatedItems[k] = struct{}{}
		}

		cpNewItems := make(map[string]struct{})
		for k, _ := range c.DepthBookMap.NewItems {
			cpNewItems[k] = struct{}{}
		}

		cache.DepthBookMap = &DepthBookMap{
			Data:         cpData,
			UpdatedItems: cpUpdatedItems,
			NewItems:     cpNewItems,
		}
	}

	if c.OrderIDsMap != nil {
		cpData := make(map[string][]string)
		for k, v := range c.OrderIDsMap.Data {
			orderIDs := make([]string, 0, len(v))
			orderIDs = append(orderIDs, v...)
			cpData[k] = orderIDs
		}

		cpUpdateItems := make(map[string]struct{})
		for k, _ := range c.OrderIDsMap.UpdatedItems {
			cpUpdateItems[k] = struct{}{}
		}

		cache.OrderIDsMap = &OrderIDsMap{
			Data:         cpData,
			UpdatedItems: cpUpdateItems,
		}
	}

	if c.PriceMap != nil {
		cpPriceMap := make(map[string]sdk.Dec)
		for k, v := range c.PriceMap {
			cpPriceMap[k] = v
		}

		cache.PriceMap = cpPriceMap
	}

	if c.ClosedOrderIDs != nil {
		cpClosedOrderIDs := make([]string, 0, len(c.ClosedOrderIDs))
		cpClosedOrderIDs = append(cpClosedOrderIDs, c.ClosedOrderIDs...)

		cache.ClosedOrderIDs = cpClosedOrderIDs
	}

	return &cache
}
