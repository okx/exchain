package dydx

import (
	"container/list"
	"math/big"
	"os"
	"sync"
	"time"

	ethcmm "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
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

	book    *DepthBook
	engine  Matcher
	gServer *OrderBookServer

	TradeTxs    *list.List
	TradeTxsMap map[ethcmm.Hash]*list.Element
	TradeTxsMtx sync.Mutex
}

func NewOrderManager(api PubSub, doMatch bool) *OrderManager {
	manager := &OrderManager{
		orders:      clist.New(),
		book:        NewDepthBook(),
		TradeTxs:    list.New(),
		TradeTxsMap: make(map[ethcmm.Hash]*list.Element),
	}

	config := DydxConfig{
		PrivKeyHex:                 "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3",
		ChainID:                    "65",
		EthWsRpcUrl:                "wss://exchaintestws.okex.org:8443",
		EthHttpRpcUrl:              "https://exchaintestrpc.okex.org",
		PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
		P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
		P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
		P1MarginAddress:            "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
		VMode:                      false,
	}

	if doMatch {
		me, err := NewMatchEngine(api, manager.book, config, nil, log.NewTMLogger(os.Stdout))
		if err != nil {
			panic(err)
		}
		manager.engine = me
	} else {
		manager.engine = NewEmptyMatcher(manager.book)
	}
	manager.gServer = NewOrderBookServer(manager.book, log.NewTMLogger(os.Stdout))
	err := manager.gServer.Start("7070")
	if err != nil {
		panic(err)
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
	d.gServer.UpdateClient()

	if result.OnChain != nil {
		go d.book.Update(result)
	} else {
		d.TradeTxsMtx.Lock()
		d.TradeTxsMap[result.Tx.Hash()] = d.TradeTxs.PushBack(result.Tx)
		d.TradeTxsMtx.Unlock()
	}
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

func (d *OrderManager) ReapMaxBytesMaxGasMaxNum(maxBytes, maxGas, maxNum int64) (tradeTxs []types.Tx, totalBytes, totalGas int64) {
	if d == nil {
		return
	}
	d.TradeTxsMtx.Lock()
	defer d.TradeTxsMtx.Unlock()

	if int64(d.TradeTxs.Len()) < maxNum {
		maxNum = int64(d.TradeTxs.Len())
	}
	tradeTxs = make([]types.Tx, 0, maxNum)

	for ele := d.TradeTxs.Front(); ele != nil; ele = ele.Next() {
		tx := ele.Value.(*ethtypes.Transaction)
		txBz, err := tx.MarshalBinary()
		if err != nil {
			continue
		}
		if maxBytes > -1 && totalBytes+int64(len(txBz)) > maxBytes {
			break
		}
		newTotalGas := totalGas + int64(tx.Gas())
		if maxGas > -1 && newTotalGas > maxGas {
			break
		}
		if len(tradeTxs) >= cap(tradeTxs) {
			break
		}
		totalGas = newTotalGas
		tradeTxs = append(tradeTxs, txBz)
	}
	return
}
