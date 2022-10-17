package dydx

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/stretchr/testify/require"
)

func TestOrderManager(t *testing.T) {
	const orderCount = 100

	orderBytes, err := hex.DecodeString(orderHex)
	require.NoError(t, err)
	var odr P1Order
	err = odr.DecodeFrom(orderBytes)
	require.NoError(t, err)

	hexPriv := "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
	priv, err := crypto.HexToECDSA(hexPriv)
	addr := crypto.PubkeyToAddress(priv.PublicKey)

	manager := NewOrderManager(false)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		var totalCount int
		defer wg.Done()
		var next *clist.CElement
		for {
			if next == nil {
				select {
				case <-manager.WaitChan():
					next = manager.Front()
				case <-time.After(time.Second):
					panic("unexpected")
				}
			}

			var wrapOdr WrapOrder
			err = wrapOdr.DecodeFrom(next.Value.(*MempoolOrder).raw)
			require.NoError(t, err)
			fmt.Println("P1Order:", wrapOdr.P1Order, "sig", hex.EncodeToString(wrapOdr.Sig), len(wrapOdr.Sig), "raw:", hex.EncodeToString(next.Value.(*MempoolOrder).Raw()))
			err = wrapOdr.P1Order.VerifySignature(wrapOdr.Sig)
			require.NoError(t, err)
			require.Equal(t, uint64(totalCount), wrapOdr.P1Order.Amount.Uint64())
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
			odr := newRandP1Order(int64(i), addr.String())
			signedOrderBytes, err := newRawSignedOrder(odr, hexPriv)
			require.NoError(t, err)

			memOrder := NewMempoolOrder(signedOrderBytes, 0)
			err = manager.Insert(memOrder)
			require.NoError(t, err)
		}
	}()
	wg.Wait()
}

func TestInsert(t *testing.T) {
	hexPriv := "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
	priv, err := crypto.HexToECDSA(hexPriv)
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(priv.PublicKey)
	manager := NewOrderManager(false)
	odr := newRandP1Order(1, addr.String())
	signedOrderBytes, err := newRawSignedOrder(odr, hexPriv)
	require.NoError(t, err)

	memOrder := NewMempoolOrder(signedOrderBytes, 0)
	err = manager.Insert(memOrder)
	require.NoError(t, err)
}

func newRawSignedOrder(odr P1Order, hexPriv string) ([]byte, error) {
	sig, err := signOrder(odr, hexPriv, 65, contractAddress)
	if err != nil {
		return nil, err
	}
	orderBytes, err := odr.encodeOrder()
	if err != nil {
		return nil, err
	}
	return append(orderBytes, sig...), nil
}
