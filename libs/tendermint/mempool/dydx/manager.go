package dydx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/dydx/contracts"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
)

type AccountRetriever interface {
	GetAccountNonce(address string) uint64
}

var (
	//Config = DydxConfig{
	//	PrivKeyHex:                 "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3",
	//	ChainID:                    "65",
	//	EthWsRpcUrl:                "wss://exchaintestws.okex.org:8443",
	//	EthHttpRpcUrl:              "https://exchaintestrpc.okex.org",
	//	PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
	//	P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
	//	P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
	//	P1MarginAddress:            "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
	//	VMode:                      false,
	//}

	Config = DydxConfig{
		PrivKeyHex:                 "2438019d3fccd8ffdff4d526c0f7fae4136866130affb3aa375d95835fa8f60f",
		ChainID:                    "8",
		EthWsRpcUrl:                "ws://localhost:8546",
		EthHttpRpcUrl:              "http://localhost:8545",
		PerpetualV1ContractAddress: "0xbc0Bf2Bf737344570c02d8D8335ceDc02cECee71",
		P1OrdersContractAddress:    "0x632D131CCCE01206F08390cB66D1AdEf9b264C61",
		P1MakerOracleAddress:       "0xF306F8B7531561d0f92BA965a163B6C6d422ade1",
		P1MarginAddress:            "0xeb95A3D1f7Ca2B8Ba61F326fC4dA9124b6C057b9",
	}

	//Config = DydxConfig{
	//	PrivKeyHex:                 "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3",
	//	ChainID:                    "8",
	//	EthWsRpcUrl:                "ws://localhost:8546",
	//	EthHttpRpcUrl:              "http://localhost:8545",
	//	PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
	//	P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
	//	P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
	//	P1MarginAddress:            "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
	//	VMode:                      true,
	//}
)

type OrderManager struct {
	orders    *clist.CList
	ordersMap sync.Map // orderKey => *clist.CElement

	historyMtx       sync.RWMutex
	addrTradeHistory map[common.Address][]*FilledP1Order
	tradeHistory     []*FilledP1Order
	trades           map[[32]byte]*FilledP1Order
	book             *DepthBook
	engine           *MatchEngine
	gServer          *OrderBookServer

	orderQueue   *OrderQueue
	waitDelete   []common.Hash
	waitUnfreeze []*MatchResult
}

func NewOrderManager(api PubSub, accRetriever AccountRetriever) *OrderManager {
	manager := &OrderManager{
		trades:           make(map[[32]byte]*FilledP1Order),
		addrTradeHistory: make(map[common.Address][]*FilledP1Order),
		orders:           clist.New(),
		book:             NewDepthBook(),
		orderQueue:       NewOrderQueue(),
	}

	me, err := NewMatchEngine(api, manager.book, Config, manager, log.NewTMLogger(os.Stdout))
	if err != nil {
		panic(err)
	}
	if accRetriever != nil {
		me.nonce = accRetriever.GetAccountNonce(me.from.String())
	} else {
		me.nonce, _ = me.httpCli.NonceAt(context.Background(), me.from, nil)
	}
	me.nonce--
	manager.engine = me

	manager.gServer = NewOrderBookServer(manager.book, log.NewTMLogger(os.Stdout))
	err = manager.gServer.Start("7070")
	if err != nil {
		panic(err)
	}
	go manager.ServeWeb()
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

	// TODO
	// should check order's filled amount from chain
	d.orderQueue.Enqueue(&wrapOdr)
	//result, err := d.engine.MatchAndTrade(&wrapOdr)
	//d.gServer.UpdateClient()
	//if err != nil {
	//	return err
	//}
	//
	//if result != nil {
	//	if result.NoSend {
	//		d.tradeTxsMtx.Lock()
	//		// d.tradeTxsMap[result.Tx.Hash()] = d.TradeTxs.PushBack(result)
	//		d.TradeTxs.PushBack(result)
	//		d.tradeTxsMtx.Unlock()
	//	}
	//	go d.book.Update(result)
	//}

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

func (d *OrderManager) updateOrderQueue(filled *contracts.P1OrdersLogOrderFilled) bool {
	if o := d.orderQueue.Get(filled.OrderHash); o != nil {
		o.LeftAmount.Sub(o.LeftAmount, filled.Fill.Amount)
		return true
	}
	return false
}

func (d *OrderManager) HandleOrderFilled(filled *contracts.P1OrdersLogOrderFilled) {
	if d.updateOrderQueue(filled) {
		return
	}

	var orderList *OrderList
	if filled.Flags[31]&FlagMaskIsBuy != FlagMaskNull {
		orderList = d.book.buyOrders
	} else {
		orderList = d.book.sellOrders
	}
	ele := orderList.Get(filled.OrderHash)
	if ele == nil {
		fmt.Println("element is nil, orderHash:", hex.EncodeToString(filled.OrderHash[:]))
		return
	}
	wodr := ele.Value.(*WrapOrder)
	wodr.Done(filled.Fill.Amount)
	d.historyMtx.Lock()
	defer d.historyMtx.Unlock()
	d.tradeHistory = append(d.tradeHistory, &FilledP1Order{
		Filled:        new(big.Int).Set(filled.Fill.Amount),
		Time:          time.Now(),
		P1OrdersOrder: wodr.P1OrdersOrder,
	})
	if filledOrder, ok := d.trades[filled.OrderHash]; ok {
		filledOrder.Filled.Add(filledOrder.Filled, filled.Fill.Amount)
		filledOrder.Time = time.Now()
	} else {
		filledOrder := &FilledP1Order{
			Filled:        new(big.Int).Set(filled.Fill.Amount),
			Time:          time.Now(),
			P1OrdersOrder: wodr.P1OrdersOrder,
		}
		d.addrTradeHistory[wodr.Maker] = append(d.addrTradeHistory[wodr.Maker], filledOrder)
		d.trades[filled.OrderHash] = filledOrder
	}

	if wodr.LeftAmount.Sign() == 0 && wodr.FrozenAmount.Sign() == 0 {
		orderList.Remove(ele)
		d.book.addrMtx.Lock()
		addrOrders := d.book.addrOrders[wodr.Maker]
		for i, order := range addrOrders {
			if order.Hash() == wodr.Hash() {
				addrOrders = append(addrOrders[:i], addrOrders[i+1:]...)
				break
			}
		}
		d.book.addrOrders[wodr.Maker] = addrOrders
		d.book.addrMtx.Unlock()
		//TODO delete broadcast queue
	}
	fmt.Println("debug filled", hex.EncodeToString(filled.OrderHash[:]), filled.TriggerPrice.String(), filled.Fill.Price.String(), filled.Fill.Amount.String())
}

func (d *OrderManager) ReapMaxBytesMaxGasMaxNum(maxBytes, maxGas, maxNum int64) (tradeTxs []types.Tx, totalBytes, totalGas int64) {
	if d == nil {
		return
	}
	if !types.HigherThanVenus(global.GetGlobalHeight()) {
		return
	}

	d.orderQueue.RLock()

	if orderQueueLen := int64(d.orderQueue.Len()); orderQueueLen < maxNum {
		maxNum = orderQueueLen
	}
	tradeTxs = make([]types.Tx, 0, maxNum)
	var shouldDelete []common.Hash

	nonce := d.engine.nonce + 1

	iter := d.orderQueue.NewIterator()
	for order := iter.Next(); order != nil; order = iter.Next() {
		mre, err := d.engine.MatchAndTrade(order)
		if err != nil || mre == nil {
			shouldDelete = append(shouldDelete, order.Hash())
			continue
		}

		if mre.Tx == nil {
			mre.Tx, _ = mre.tradeOps.Commit(&bind.TransactOpts{NoSend: true, Nonce: new(big.Int).SetUint64(nonce)})
			if mre.Tx == nil {
				shouldDelete = append(shouldDelete, order.Hash())
				mre.Unfreeze()
				continue
			}
			d.engine.logger.Debug("reap tx", "tx", mre.Tx.Hash().String())
		}
		tx := mre.Tx
		txBz, err := tx.MarshalBinary()
		if err != nil {
			shouldDelete = append(shouldDelete, order.Hash())
			mre.Unfreeze()
			continue
		}
		d.waitDelete = append(d.waitDelete, order.Hash())
		d.waitUnfreeze = append(d.waitUnfreeze, mre)
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
		nonce++
	}

	d.orderQueue.RUnlock()

	for _, hash := range shouldDelete {
		d.orderQueue.Delete(hash)
	}
	return
}

func (d *OrderManager) UpdateAddress(sender string, nonce uint64, code uint32) {
	if sender == d.engine.from.String() &&
		(code == abci.CodeTypeOK || code > abci.CodeTypeNonceInc) {
		d.engine.nonce = nonce
	}
}

func (d *OrderManager) Update(txsResps []*abci.ResponseDeliverTx) {
	if d == nil {
		return
	}

	for _, hash := range d.waitDelete {
		d.orderQueue.Delete(hash)
	}
	d.waitDelete = d.waitDelete[:0]

	for _, mre := range d.waitUnfreeze {
		mre.Unfreeze()
	}
	d.waitUnfreeze = d.waitUnfreeze[:0]

	d.engine.UpdateState(txsResps)
}

func (d *OrderManager) OrderQueueLen() int {
	if d == nil {
		return 0
	}
	return d.orderQueue.Len()
}

// order -> (check) fifo (broadcast)

// block -> mempool.Update(tx resp) -> fifo + orderBook

// propose foreach fifo -> match -> orderbook -> txs -> block
