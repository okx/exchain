package dydx

import (
	"math/big"
	"sync"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
)

type Matcher interface {
	MatchAndTrade(order *WrapOrder) (*MatchResult, error)
}

type emptyMatcher struct {
	book *DepthBook
}

func (e emptyMatcher) MatchAndTrade(order *WrapOrder) (*MatchResult, error) {
	err := e.book.Insert(order)
	return nil, err
}

func NewEmptyMatcher(book *DepthBook) Matcher {
	return emptyMatcher{
		book: book,
	}
}

type OrderManager struct {
	orders    *clist.CList
	ordersMap sync.Map // orderKey => *clist.CElement

	book   *DepthBook
	engine Matcher
}

func NewOrderManager(doMatch bool) *OrderManager {
	manager := &OrderManager{
		orders: clist.New(),
		book:   NewDepthBook(),
	}

	config := DydxConfig{
		PrivKeyHex:                 "fefac29bfa769d8a6c17b685816dadbd30e3f395e997ed955a5461914be75ed5",
		ChainID:                    "65",
		EthWsRpcUrl:                "wss://exchaintestws.okex.org:8443",
		PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
		P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
		P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
		P1MarginAddress:            "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
	}

	if doMatch {
		me, err := NewMatchEngine(manager.book, config, nil, nil)
		if err != nil {
			return nil
		}
		manager.engine = me
	} else {
		manager.engine = NewEmptyMatcher(manager.book)
	}
	go manager.Serve()
	return manager
}

func (d *OrderManager) Insert(memOrder *MempoolOrder) error {
	var wrapOdr WrapOrder
	if err := wrapOdr.DecodeFrom(memOrder.Raw()); err != nil {
		return err
	}
	if err := wrapOdr.P1Order.VerifySignature(wrapOdr.Sig); err != nil {
		return err
	}

	if wrapOdr.Expiration.Cmp(big.NewInt(time.Now().Unix())) <= 0 {
		return ErrExpiredOrder
	}

	ele := d.orders.PushBack(memOrder)
	d.ordersMap.Store(memOrder.Key(), ele)

	result, err := d.engine.MatchAndTrade(&wrapOdr)
	if err != nil {
		return err
	}

	go d.book.Update(result)
	return nil
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
