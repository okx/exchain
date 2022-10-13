package dydx

import (
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/types/time"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"testing"
)

var testAmounts = []int64{10, 51, 21, 51, 37, 20, 58, 48, 16, 49, 84, 87, 74, 36, 15, 73, 68, 91, 90, 31, 73, 56, 11, 37, 78, 9, 72, 50, 88, 71, 44, 43, 23, 59, 3, 39, 83, 7, 32, 80, 15, 16, 20, 71, 52, 7, 19, 62, 10, 97, 5, 21, 17, 12, 49, 79, 33, 25, 46, 34, 18, 1, 4, 10, 24, 17, 69, 79, 69, 50, 57, 23, 12, 2, 55, 57, 43, 1, 19, 18, 94, 41, 88, 21, 11, 80, 43, 55, 41, 53, 31, 17, 31, 18, 32, 14, 4, 80, 0, 71}

func TestOrderList(t *testing.T) {
	orders := createTestOrder(100)
	sellList := NewOrderList(false)
	for _, order := range orders {
		sellList.Insert(order)
	}
	for ele := sellList.Front(); ele != nil && ele.Next() != nil; ele = ele.Next() {
		require.True(t, ele.Value.(*WrapOrder).Price().Uint64() <= ele.Next().Value.(*WrapOrder).Price().Uint64())
	}

	buyList := NewOrderList(true)
	for _, order := range orders {
		buyList.Insert(order)
	}
	for ele := buyList.Front(); ele != nil && ele.Next() != nil; ele = ele.Next() {
		require.True(t, ele.Value.(*WrapOrder).Price().Uint64() >= ele.Next().Value.(*WrapOrder).Price().Uint64())
	}
}

func createTestOrder(n int) []*WrapOrder {
	var orders []*WrapOrder
	for i := 0; i < n && i < len(testAmounts); i++ {
		orders = append(orders, newWrapOrder(testAmounts[i]))
	}
	return orders
}

func newWrapOrder(amount int64) *WrapOrder {
	return &WrapOrder{
		P1Order: P1Order{
			P1OrdersOrder: contracts.P1OrdersOrder{
				Amount:       big.NewInt(amount),
				LimitPrice:   big.NewInt(0),
				TriggerPrice: big.NewInt(0),
				LimitFee:     big.NewInt(0),
				Expiration:   big.NewInt(time.Now().Unix()*2 + rand.Int63()),
			},
		},
	}
}
