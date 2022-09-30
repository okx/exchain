package dydx

import (
	"encoding/hex"
	"sync"
	"testing"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/stretchr/testify/require"
)

func TestOrderManager(t *testing.T) {
	const orderCount = 100

	orderBytes, err := hex.DecodeString(orderHex)
	require.NoError(t, err)
	var odr Order
	err = odr.DecodeFrom(orderBytes)
	require.NoError(t, err)

	book := NewOrderManager()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		var totalCount int
		defer wg.Done()
		var next *clist.CElement
		for {
			if next == nil {
				select {
				case <-book.WaitChan():
					next = book.Front()
				}
			}
			var signedOrder SignedOrder
			err = signedOrder.DecodeFrom(next.Value.(*MempoolOrder).raw)
			require.NoError(t, err)
			var odr Order
			err = odr.DecodeFrom(signedOrder.Msg)
			require.NoError(t, err)
			require.Equal(t, uint64(totalCount), odr.Amount.Uint64())
			totalCount++
			select {
			case <-next.NextWaitChan():
				// see the start of the for loop for nil check
				next = next.Next()
			case <-time.After(time.Millisecond * 10):
				require.Equal(t, orderCount, totalCount)
				return
			}
		}

	}()
	go func() {
		defer wg.Done()
		for i := 0; i < orderCount; i++ {
			if i%(orderCount/10) == 0 {
				time.Sleep(time.Millisecond)
			}
			odr.Amount.SetInt64(int64(i))
			orderBytes, err := orderTuple.Encode(odr)
			require.NoError(t, err)
			signedOrder := SignedOrder{
				Msg: orderBytes,
			}
			signedOrderBytes, err := signedTuple.Encode(signedOrder)
			require.NoError(t, err)

			memOrder := NewMempoolOrder(signedOrderBytes, 0)
			err = book.Insert(memOrder, 0)
			require.NoError(t, err)
		}
	}()
	wg.Wait()
}
