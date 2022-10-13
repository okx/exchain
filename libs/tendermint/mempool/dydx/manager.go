package dydx

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
)

type OrderManager struct {
	orders    *clist.CList
	ordersMap sync.Map // orderKey => *clist.CElement

	book *DepthBook
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: clist.New(),
		book:   NewDepthBook(),
	}
}

func (d *OrderManager) Insert(memOrder *MempoolOrder) error {
	var signedOdr SignedOrder
	if err := signedOdr.DecodeFrom(memOrder.raw); err != nil {
		return err
	}
	var odr P1Order
	if err := odr.DecodeFrom(signedOdr.Msg); err != nil {
		return err
	}
	if odr.Expiration.Uint64() == 0 {
		return fmt.Errorf("invalid expiration")
	}
	if err := odr.VerifySignature(signedOdr.Sig[:]); err != nil {
		return err
	}

	ele := d.orders.PushBack(memOrder)
	d.ordersMap.Store(memOrder.Key(), ele)

	wrapOdr := &WrapOrder{
		P1Order:    odr,
		LeftAmount: new(big.Int).Set(odr.Amount),
		Raw:        memOrder.Raw(),
		Sig:        signedOdr.Sig[:],
	}
	return d.book.Insert(wrapOdr)
}

func (d *OrderManager) Remove(order OrderRaw) {
	ele, ok := d.ordersMap.LoadAndDelete(order.Key())
	if !ok {
		return
	}
	d.orders.Remove(ele.(*clist.CElement))
}

func (d *OrderManager) Load(order OrderRaw) *clist.CElement {
	v, ok := d.ordersMap.Load(order.Key())
	if !ok {
		return nil
	}
	return v.(*clist.CElement)
}

func (d *OrderManager) WaitChan() <-chan struct{} {
	return d.orders.WaitChan()
}

func (d *OrderManager) Front() *clist.CElement {
	return d.orders.Front()
}