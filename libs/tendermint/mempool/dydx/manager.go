package dydx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/dydx"
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
	AddressForOrder = "0xF1730217Bd65f86d2F0000000000000000000000"
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
		EthHttpRpcUrl:              "http://localhost:8545",
		PerpetualV1ContractAddress: "0xbc0Bf2Bf737344570c02d8D8335ceDc02cECee71",
		P1OrdersContractAddress:    "0x632D131CCCE01206F08390cB66D1AdEf9b264C61",
		P1MakerOracleAddress:       "0xF306F8B7531561d0f92BA965a163B6C6d422ade1",
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

	signals          chan struct{}
	marketPriceMtx   sync.RWMutex
	marketPrice      *big.Int
	balancesMtx      sync.RWMutex
	balances         map[common.Address]*contracts.P1TypesBalance
	historyMtx       sync.RWMutex
	addrTradeHistory map[common.Address][]*FilledP1Order
	tradeHistory     []*FilledP1Order
	trades           map[[32]byte]*FilledP1Order
	book             *DepthBook
	engine           *MatchEngine
	gServer          *OrderBookServer

	orderQueue   *OrderQueue
	waitUnfreeze []*MatchResult

	currentBlockTxs []types.Tx
	totalBytes      int64
	totalGas        int64

	filledOrCanceledOrders sync.Map

	logger log.Logger
}

func NewOrderManager(api PubSub, accRetriever AccountRetriever, logger log.Logger) *OrderManager {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	manager := &OrderManager{
		signals:          make(chan struct{}, 10),
		balances:         make(map[common.Address]*contracts.P1TypesBalance),
		trades:           make(map[[32]byte]*FilledP1Order),
		addrTradeHistory: make(map[common.Address][]*FilledP1Order),
		orders:           clist.New(),
		book:             NewDepthBook(),
		orderQueue:       NewOrderQueue(),
		logger:           logger,
	}

	me, err := NewMatchEngine(api, accRetriever, manager.book, Config, manager, manager.logger)
	if err != nil {
		panic(err)
	}
	manager.engine = me
	go manager.updateMarketPriceRoutine()
	manager.SendSignal()

	manager.gServer = NewOrderBookServer(manager.book, manager.logger)
	port := "7070"
	_ = manager.gServer.Start(port)
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

	d.logger.Debug("pre enqueue", "order", &wrapOdr)

	hash := [32]byte(wrapOdr.Hash())
	if _, ok := d.filledOrCanceledOrders.Load(hash); ok {
		d.logger.Debug("order is filled or canceled", "order", hash)
		return nil
	}

	ordersStatus, err := d.engine.contracts.P1Orders.GetOrdersStatus(nil, [][32]byte{hash})
	if err == nil {
		status := ordersStatus[0]
		if status.Status == 2 {
			d.filledOrCanceledOrders.Store(hash, struct{}{})
			d.logger.Debug("order is canceled", "order", hash)
			return nil
		}
		if status.FilledAmount.Sign() > 0 {
			wrapOdr.LeftAmount.Sub(wrapOdr.Amount, status.FilledAmount)
			if wrapOdr.LeftAmount.Sign() == 0 {
				d.filledOrCanceledOrders.Store(hash, struct{}{})
				d.logger.Debug("order is full filled", "order", hash)
				return nil
			}
			d.logger.Debug("order is partially filled", "order", hash, "left", wrapOdr.LeftAmount)
		}
	}

	err = d.checkBalance(&wrapOdr)
	if err != nil {
		return err
	}

	ok := d.orderQueue.Enqueue(&wrapOdr)
	d.logger.Debug("enqueue", "order", wrapOdr.Hash(), "ok", ok)

	return nil
}

func (d *OrderManager) SendSignal() {
	d.signals <- struct{}{}
}

func (d *OrderManager) GetMarketPrice() *big.Int {
	d.marketPriceMtx.RLock()
	defer d.marketPriceMtx.RUnlock()
	return d.marketPrice
}

func (d *OrderManager) updateMarketPriceRoutine() {
	for range d.signals {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		marketPrice, err := d.engine.contracts.P1MakerOracle.GetPrice(&bind.CallOpts{
			From:    d.engine.contracts.Addresses.PerpetualV1,
			Context: ctx,
		})
		cancel()
		if err != nil {
			d.logger.Error("UpdateMarketPrice", "GetPrice error", err)
			continue
		}
		d.marketPriceMtx.Lock()
		d.marketPrice = marketPrice
		d.marketPriceMtx.Unlock()
	}
}

func (d *OrderManager) checkBalance(order *WrapOrder) error {
	balance := d.getBalance(order.Maker)
	if balance == nil {
		p1Balance, err := d.engine.contracts.PerpetualV1.GetAccountBalance(nil, order.Maker)
		if err != nil {
			d.logger.Error("checkBalance", "GetAccountBalance error", err, "addr", order.Maker)
			return nil
		}
		d.setBalance(order.Maker, &p1Balance)
		balance = &p1Balance
	}
	marketPrice := d.GetMarketPrice()
	if marketPrice == nil {
		return nil
	}
	margin := negBig(new(big.Int).Mul(balance.Margin, exp18), balance.MarginIsPositive)
	position := negBig(new(big.Int).Mul(balance.Position, marketPrice), balance.PositionIsPositive)
	perpetualBalance := new(big.Int).Add(margin, position)
	if perpetualBalance.Sign() < 0 {
		d.logger.Debug("checkBalance", "addr", order.Maker, "perpetualBalance:", perpetualBalance)
		return ErrMarginNotEnough
	}

	cost := new(big.Int).Sub(order.LimitPrice, marketPrice)
	cost.Mul(cost, order.LeftAmount)
	if !order.isBuy() {
		cost.Neg(cost)
	}
	if perpetualBalance.Cmp(cost) < 0 {
		d.logger.Debug("checkBalance", "addr", order.Maker, "perpetualBalance:", perpetualBalance, "marketPrice", marketPrice, "cost", cost)
		return ErrMarginNotEnough
	}
	return nil
}

func (d *OrderManager) getBalance(addr common.Address) *contracts.P1TypesBalance {
	d.balancesMtx.RLock()
	defer d.balancesMtx.RUnlock()
	return d.balances[addr]
}

func (d *OrderManager) setBalance(addr common.Address, balance *contracts.P1TypesBalance) {
	d.balancesMtx.Lock()
	defer d.balancesMtx.Unlock()
	d.balances[addr] = balance
}

func (d *OrderManager) Remove(order OrderRaw) {
	ele, ok := d.ordersMap.LoadAndDelete(order.Key())
	if !ok {
		return
	}
	d.orders.Remove(ele.(*clist.CElement))
}

func (d *OrderManager) CancelOrder(order OrderRaw) {
	d.Remove(order)
	var p1Order P1Order
	err := p1Order.DecodeFrom(order)
	if err != nil {
		fmt.Println("decode order error:", err)
		return
	}
	d.book.DeleteByHash(p1Order.Hash())
	//TODO
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

func (d *OrderManager) updateOrderQueue(filled *contracts.P1OrdersLogOrderFilled) *WrapOrder {
	var o *WrapOrder
	if o = d.orderQueue.Get(filled.OrderHash); o != nil {
		o.LeftAmount.Sub(o.LeftAmount, filled.Fill.Amount)
		d.logger.Debug("update order queue", "order", o.Hash(), "filled", filled.Fill.Amount, "left", o.LeftAmount)
		if o.LeftAmount.Sign() == 0 {
			d.orderQueue.Delete(filled.OrderHash)
			d.logger.Debug("delete order queue", "order", o.Hash())
		}
		return o
	}
	return o
}

func (d *OrderManager) HandleTrade(trade *contracts.PerpetualV1LogTrade) {
	if trade != nil {
		var makerBalance = dydx.Bytes32ToBalance(&trade.MakerBalance)
		var takerBalance = dydx.Bytes32ToBalance(&trade.TakerBalance)
		d.logger.Debug("HandleTrade",
			"taker", trade.Taker,
			"takerBalance", dydx.P1TypesBalanceStringer(takerBalance),
			"maker", trade.Maker,
			"makerBalance", dydx.P1TypesBalanceStringer(makerBalance),
		)
	}
}

func (d *OrderManager) HandleWithdraw(withdraw *contracts.PerpetualV1LogWithdraw) {
	if withdraw != nil {
		balance := dydx.Bytes32ToBalance(&withdraw.Balance)
		d.logger.Debug("HandleWithdraw",
			"addr", withdraw.Account, "to", withdraw.Destination,
			"amount", withdraw.Amount, "balance", dydx.P1TypesBalanceStringer(balance),
		)
	}
}

func (d *OrderManager) HandleDeposit(deposit *contracts.PerpetualV1LogDeposit) {
	if deposit != nil {
		balance := dydx.Bytes32ToBalance(&deposit.Balance)
		d.logger.Debug("HandleDeposit",
			"addr", deposit.Account,
			"amount", deposit.Amount,
			"balance", dydx.P1TypesBalanceStringer(balance),
		)
	}
}

func (d *OrderManager) HandleIndex(index *contracts.PerpetualV1LogIndex) {
	if index != nil {
		index := dydx.Bytes32ToIndex(&index.Index)
		timeIndex := time.Unix(int64(index.Timestamp), 0).Local()
		valueStr := index.Value.String()
		if !index.IsPositive {
			valueStr = "-" + valueStr
		}
		d.logger.Debug("HandleIndex", "time", timeIndex, "value", valueStr)
	}
}

func (d *OrderManager) HandleAccountSettled(settled *contracts.PerpetualV1LogAccountSettled) {
	if settled != nil {
		d.logger.Debug("HandleAccountSettled", "addr", settled.Account)
	}
}

func (d *OrderManager) HandleOrderCanceled(canceled *contracts.P1OrdersLogOrderCanceled) {
	if canceled != nil {
		d.filledOrCanceledOrders.Store(canceled.OrderHash, nil)
		d.orderQueue.Delete(canceled.OrderHash)
	}
}

func (d *OrderManager) HandleOrderFilled(filled *contracts.P1OrdersLogOrderFilled) {
	wodr := d.updateOrderQueue(filled)
	if wodr == nil {
		var orderList *OrderList
		if filled.Flags[31]&FlagMaskIsBuy != FlagMaskNull {
			orderList = d.book.buyOrders
		} else {
			orderList = d.book.sellOrders
		}
		ele := orderList.Get(filled.OrderHash)
		if ele == nil {
			d.logger.Debug("element is nil, orderHash:", hex.EncodeToString(filled.OrderHash[:]))
			return
		}
		wodr = ele.Value.(*WrapOrder)
		wodr.Done(filled.Fill.Amount)

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
	}
	balance := d.getBalance(wodr.Maker)
	positionDiff := negBig(filled.Fill.Amount, balance.PositionIsPositive)
	marginDiff := negBig(new(big.Int).Mul(filled.Fill.Amount, filled.Fill.Price), balance.MarginIsPositive)
	marginDiff.Div(marginDiff, exp18)

	if wodr.isBuy() {
		balance.Position.Add(balance.Position, positionDiff)
		balance.Margin.Sub(balance.Margin, marginDiff)
	} else {
		balance.Position.Sub(balance.Position, positionDiff)
		balance.Margin.Add(balance.Margin, marginDiff)
	}
	if balance.Position.Sign() < 0 {
		balance.Position.Neg(balance.Position)
		balance.PositionIsPositive = !balance.PositionIsPositive
	}
	if balance.Margin.Sign() < 0 {
		balance.Margin.Neg(balance.Margin)
		balance.MarginIsPositive = !balance.MarginIsPositive
	}

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

	fmt.Println("debug filled", hex.EncodeToString(filled.OrderHash[:]), filled.TriggerPrice.String(), filled.Fill.Price.String(), filled.Fill.Amount.String())
}

func (d *OrderManager) ReapMaxBytesMaxGasMaxNum(maxBytes, maxGas, maxNum int64) (tradeTxs []types.Tx, totalBytes, totalGas int64) {
	if d == nil {
		return
	}
	if !types.HigherThanVenus(global.GetGlobalHeight()) {
		return
	}

	if len(d.currentBlockTxs) > 0 {
		return d.currentBlockTxs, d.totalBytes, d.totalGas
	}
	queueLen := d.orderQueue.Len()
	if queueLen == 0 {
		return
	}

	iterCount := 0
	defer func() {
		d.logger.Debug("finish reap order", "iterCount", iterCount, "totalBytes", totalBytes, "totalGas", totalGas)
		for i := 0; i < iterCount; i++ {
			d.orderQueue.Dequeue()
		}
		d.gServer.UpdateClient()
	}()

	d.logger.Debug("start reap order", "queue-size", queueLen)

	preMakeCap := maxNum
	if orderQueueLen := int64(queueLen); orderQueueLen < maxNum {
		preMakeCap = orderQueueLen
	}
	tradeTxs = make([]types.Tx, 0, preMakeCap)
	nonce := d.engine.nonce + 1

	d.orderQueue.Foreach(func(order *WrapOrder, index int, count int) bool {
		iterCount++
		mre, err := d.engine.MatchAndTrade(order)
		if err != nil || mre == nil {
			return true
		}

		if mre.Tx == nil {
			mre.Tx, err = mre.tradeOps.Commit(&bind.TransactOpts{NoSend: true, Nonce: new(big.Int).SetUint64(nonce)})
			if err != nil {
				d.logger.Error("commit trade tx failed", "err", err)
			}
			if mre.Tx == nil {
				mre.Unfreeze()
				return true
			}
			d.engine.logger.Debug("reap tx", "tx", mre.Tx.Hash().String())
		}
		tx := mre.Tx
		txBz, err := tx.MarshalBinary()
		if err != nil {
			mre.Unfreeze()
			return true
		}

		aminoOverhead := types.ComputeAminoOverhead(txBz, 1)
		if maxBytes > -1 && totalBytes+int64(len(txBz))+aminoOverhead > maxBytes {
			iterCount--
			d.engine.Rollback(mre)
			return false
		}
		totalBytes += int64(len(txBz)) + aminoOverhead
		newTotalGas := totalGas + int64(tx.Gas())
		if maxGas > -1 && newTotalGas > maxGas {
			iterCount--
			d.engine.Rollback(mre)
			return false
		}
		if len(tradeTxs) >= int(maxNum) {
			iterCount--
			d.engine.Rollback(mre)
			return false
		}
		totalGas = newTotalGas
		tradeTxs = append(tradeTxs, txBz)
		d.waitUnfreeze = append(d.waitUnfreeze, mre)
		nonce++
		return true
	})
	d.currentBlockTxs = tradeTxs
	d.totalBytes = totalBytes
	d.totalGas = totalGas
	return
}

func (d *OrderManager) UpdateAddress(sender string, nonce uint64, code uint32) {
	if d == nil {
		return
	}
	if sender == d.engine.from.String() &&
		(code == abci.CodeTypeOK || code > abci.CodeTypeNonceInc) {
		d.engine.nonce = nonce
	}
}

func (d *OrderManager) Update(txsResps []*abci.ResponseDeliverTx) {
	if d == nil {
		return
	}

	if len(d.currentBlockTxs) > 0 {
		d.currentBlockTxs = d.currentBlockTxs[:0]
		d.totalBytes = 0
		d.totalGas = 0
	}

	for _, mre := range d.waitUnfreeze {
		mre.Unfreeze()
	}
	d.waitUnfreeze = d.waitUnfreeze[:0]

	d.engine.UpdateState(txsResps)

	d.gServer.UpdateClient()
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
