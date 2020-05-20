package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
type Deal struct {
	OrderID  string  `json:"order_id"`
	Side     string  `json:"side"`
	Quantity sdk.Dec `json:"quantity"`
	Fee      string  `json:"fee"`
}

// nolint
type MatchResult struct {
	BlockHeight int64   `json:"block_height"`
	Price       sdk.Dec `json:"price"`
	Quantity    sdk.Dec `json:"quantity"`
	Deals       []Deal  `json:"deals"`
}

// nolint
type BlockMatchResult struct {
	BlockHeight int64                  `json:"block_height"`
	ResultMap   map[string]MatchResult `json:"result_map"`
	TimeStamp   int64                  `json:"timestamp"`
}

// nolint
type DepthBookItem struct {
	Price        sdk.Dec `json:"price"`
	BuyQuantity  sdk.Dec `json:"buy_quantity"`
	SellQuantity sdk.Dec `json:"sell_quantity"`
}

// nolint
type DepthBook struct {
	Items []DepthBookItem `json:"items"`
}

// InsertOrder : Items in depth book are sorted by price desc
// insert a new order into depth book
func (depthBook *DepthBook) InsertOrder(order *Order) {
	bookLength := len(depthBook.Items)
	newItem := DepthBookItem{
		Price:        order.Price,
		BuyQuantity:  sdk.ZeroDec(),
		SellQuantity: sdk.ZeroDec(),
	}
	if order.Side == BuyOrder {
		newItem.BuyQuantity = order.RemainQuantity
	} else {
		newItem.SellQuantity = order.RemainQuantity
	}
	if bookLength == 0 || order.Price.LT(depthBook.Items[bookLength-1].Price) {
		depthBook.Items = append(depthBook.Items, newItem)
		return
	}

	// find first index, s.t. order.Price >= depthBook[index].Price
	index := sort.Search(bookLength, func(i int) bool {
		return order.Price.GTE(depthBook.Items[i].Price)
	})

	if order.Price.Equal(depthBook.Items[index].Price) {
		if order.Side == BuyOrder {
			depthBook.Items[index].BuyQuantity =
				depthBook.Items[index].BuyQuantity.Add(order.RemainQuantity)
		} else {
			depthBook.Items[index].SellQuantity =
				depthBook.Items[index].SellQuantity.Add(order.RemainQuantity)
		}
	} else { // order.InitPrice > depthBook[index].InitPrice
		rear := append([]DepthBookItem{newItem}, depthBook.Items[index:]...)
		depthBook.Items = append(depthBook.Items[:index], rear...)
	}
}

// RemoveOrder : remove an order from depth book when order cancelled/expired
func (depthBook *DepthBook) RemoveOrder(order *Order) {
	bookLen := len(depthBook.Items)
	// find first index, s.t. order.Price >= depthBook[index].Price
	// i.e. order.Price == depthBook[index].Price
	index := sort.Search(bookLen, func(i int) bool {
		return order.Price.GTE(depthBook.Items[i].Price)
	})

	if index < bookLen && depthBook.Items[index].Price.Equal(order.Price) {
		if order.Side == BuyOrder {
			depthBook.Items[index].BuyQuantity =
				depthBook.Items[index].BuyQuantity.Sub(order.RemainQuantity)
		} else if order.Side == SellOrder {
			depthBook.Items[index].SellQuantity =
				depthBook.Items[index].SellQuantity.Sub(order.RemainQuantity)
		}

		depthBook.RemoveIfEmpty(index)
	}
}

// Sub : subtract the buy or sell quantity
func (depthBook *DepthBook) Sub(index int, num sdk.Dec, side string) {
	if side == BuyOrder {
		depthBook.Items[index].BuyQuantity = depthBook.Items[index].BuyQuantity.Sub(num)
	} else if side == SellOrder {
		depthBook.Items[index].SellQuantity = depthBook.Items[index].SellQuantity.Sub(num)
	}
}

// RemoveIfEmpty : remove the filled or empty item
func (depthBook *DepthBook) RemoveIfEmpty(index int) bool {
	res := depthBook.Items[index].BuyQuantity.IsZero() && depthBook.Items[index].SellQuantity.IsZero()
	if res {
		depthBook.Items = append(depthBook.Items[:index], depthBook.Items[index+1:]...)
	}
	return res
}

// Copy : depth copy of depth book
func (depthBook *DepthBook) Copy() *DepthBook {
	itemList := make([]DepthBookItem, 0, len(depthBook.Items))
	itemList = append(itemList, depthBook.Items...)
	return &DepthBook{Items: itemList}
}
