package dydx

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/dydx/contracts"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/tendermint/mempool/placeorder"
	"github.com/spf13/viper"
)

const (
	POST = "POST"
	GET  = "GET"

	addrKey    = "addr"
	timeFormat = "15:04:05"

	// 0x4Ef308B36E9f75C97a38594acbFa9FBe1B847Da5 testnet
	// 0x2594E83A94F89Ffb923773ddDfF723BbE017b80D localnet
	placeOrderContractAddr = "0x2594E83A94F89Ffb923773ddDfF723BbE017b80D"
)

var oneWeekSeconds = int64(time.Hour/time.Second) * 24 * 7
var exp18, _ = new(big.Int).SetString("1000000000000000000", 10)

type Response struct {
	Succeed  bool   `json:"succeed"`
	ErrorMsg string `json:"errorMsg"`
}

func (o *OrderManager) ServeWeb() {
	r := mux.NewRouter()
	r.HandleFunc("/order", o.GenerateOrderHandler).Methods(GET).Queries("amount", "{amount}", "limitPrice", "{limitPrice}", "maker", "{maker}", "isBuy", "{isBuy}")
	r.HandleFunc("/placeorder", o.SendHandler).Methods(GET).Queries("signedOrder", "{signedOrder}")

	r.HandleFunc("/book", o.BookHandler).Methods(GET)
	r.HandleFunc("/trades", o.TradesHandler).Methods(GET)
	r.HandleFunc("/position", o.PositionHandler).Methods(GET).Queries("addr", "{addr}")
	r.HandleFunc("/orders", o.OrdersHandler).Methods(GET).Queries("addr", "{addr}")
	r.HandleFunc("/fills", o.FillsHandler).Methods(GET).Queries("addr", "{addr}")
	r.HandleFunc("/drop", o.DropHandler).Methods(GET).Queries("amount", "{amount}", "addr", "{addr}")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+viper.GetString("dydx.web-port"), r))
}

type OrderResponse struct {
	Order string `json:"order"`
	Hash  string `json:"hash"`
}

func (o *OrderManager) GenerateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	amount := vars["amount"]
	Amount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		fmt.Fprintf(w, "invalid amount")
		return
	}
	limitPrice := vars["limitPrice"]
	LimitPrice, ok := new(big.Int).SetString(limitPrice, 10)
	if !ok {
		fmt.Fprintf(w, "invalid limitPrice")
		return
	}
	LimitPrice = LimitPrice.Mul(LimitPrice, exp18)

	maker := vars["maker"]
	isBuy := vars["isBuy"]
	caller, err := placeorder.NewPlaceorderCaller(common.HexToAddress(placeOrderContractAddr), o.engine.httpCli)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	order := placeorder.OrdersOrder{
		Amount:       Amount,
		LimitPrice:   LimitPrice,
		TriggerPrice: big.NewInt(0),
		LimitFee:     big.NewInt(0),
		Maker:        common.HexToAddress(maker),
		Expiration:   big.NewInt(time.Now().Unix() + oneWeekSeconds),
	}
	if isBuy == "true" {
		order.Flags[31] = 1
	}
	msg, err := caller.GetOrderMessage(&bind.CallOpts{From: common.HexToAddress(maker), Context: context.Background()}, order)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	hash, err := caller.GetOrderHash(&bind.CallOpts{From: common.HexToAddress(maker), Context: context.Background()}, order)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	data, err := json.Marshal(OrderResponse{
		Order: hex.EncodeToString(msg),
		Hash:  hex.EncodeToString(hash[:]),
	})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))

}

func (o *OrderManager) SendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	hexSignedOrder := vars["signedOrder"]
	hexSignedOrder = strings.TrimPrefix(hexSignedOrder, "0x")
	signedOrder, err := hex.DecodeString(hexSignedOrder)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	err = o.Insert(NewMempoolOrder(signedOrder, 0))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "OK")
}

func (o *OrderManager) BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	levels := bookToLevel(o.orderQueue.book)
	data, err := json.Marshal(levels)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

type Trade struct {
	Size  int64  `json:"size"`
	Price string `json:"price"`
	Side  string `json:"side"`
	Time  string `json:"time"`
}

func (o *OrderManager) TradesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")

	trades := make([]*Trade, 0)
	o.historyMtx.RLock()
	for _, t := range o.tradeHistory {
		trade := &Trade{
			Size:  t.Filled.Int64(),
			Price: new(big.Int).Div(t.LimitPrice, exp18).String(),
			Time:  t.Time.Format(timeFormat),
		}
		if t.Flags[31] == 1 {
			trade.Side = "buy"
		} else {
			trade.Side = "sell"
		}
		trades = append(trades, trade)
	}
	o.historyMtx.RUnlock()
	data, err := json.Marshal(trades)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

type Balance struct {
	Margin        *big.Int `json:"margin"`
	Position      *big.Int `json:"position"`
	Erc20Balance  *big.Int `json:"erc20Balance"`
	MarginCache   *big.Int `json:"marginCache"`
	PositionCache *big.Int `json:"positionCache"`
}

func (o *OrderManager) PositionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])
	p1Balance, err := o.engine.contracts.PerpetualV1.GetAccountBalance(nil, addr)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	token, err := contracts.NewTestToken(o.engine.contracts.P1MarginAddress, o.engine.httpCli)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	balance, err := token.BalanceOf(nil, addr)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	p1BalanceCache := o.getBalance(addr)
	if p1BalanceCache == nil {
		p1BalanceCache = &contracts.P1TypesBalance{
			Margin:   big.NewInt(0),
			Position: big.NewInt(0),
		}
	}
	fmt.Println("p1 balance cache", p1BalanceCache)

	data, err := json.Marshal(&Balance{
		Margin:        negBig(p1Balance.Margin, p1Balance.MarginIsPositive),
		Position:      negBig(p1Balance.Position, p1Balance.PositionIsPositive),
		Erc20Balance:  balance,
		MarginCache:   negBig(p1BalanceCache.Margin, p1BalanceCache.MarginIsPositive),
		PositionCache: negBig(p1BalanceCache.Position, p1Balance.PositionIsPositive),
	})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

func negBig(n *big.Int, positive bool) *big.Int {
	if positive {
		return n
	}
	return new(big.Int).Neg(n)
}

type WebOrder struct {
	Order        string `json:"order"`
	OrderHash    string `json:"orderHash"`
	Status       string `json:"status"`
	IsBuy        bool   `json:"isBuy"`
	Amount       int64  `json:"amount"`
	FilledAmount int64  `json:"filledAmount"`
	Price        string `json:"price"`
	TriggerPrice string `json:"triggerPrice"`
	Expiration   string `json:"expiration"`
}

type ExportP1Order struct {
	Flags        [32]byte `json:"flags"`
	Amount       string   `json:"amount"`
	LimitPrice   string   `json:"limitPrice"`
	TriggerPrice string   `json:"triggerPrice"`
	LimitFee     string   `json:"limitFee"`
	Maker        string   `json:"maker"`
	Taker        string   `json:"taker"`
	Expiration   string   `json:"expiration"`
}

func (o *OrderManager) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])

	orders := make([]*WebOrder, 0)
	o.orderQueue.book.addrMtx.RLock()
	for _, order := range o.orderQueue.book.addrOrders[addr] {
		exportOrder := ExportP1Order{
			Flags:        order.Flags,
			Amount:       order.Amount.String(),
			LimitPrice:   order.LimitPrice.String(),
			TriggerPrice: order.TriggerPrice.String(),
			LimitFee:     order.LimitFee.String(),
			Maker:        order.Maker.String(),
			Taker:        order.Taker.String(),
			Expiration:   order.Expiration.String(),
		}
		data, _ := json.Marshal(exportOrder)
		orders = append(orders, &WebOrder{
			Order:        string(data),
			OrderHash:    order.orderHash.String(),
			Status:       "limit",
			IsBuy:        order.Flags[31] == 1,
			Amount:       order.Amount.Int64(),
			FilledAmount: new(big.Int).Sub(order.Amount, order.LeftAndFrozen()).Int64(),
			Price:        new(big.Int).Div(order.LimitPrice, exp18).String(),
			TriggerPrice: order.TriggerPrice.String(),
			Expiration:   fmt.Sprintf("%d hours", (order.Expiration.Int64()-time.Now().Unix())/3600),
		})
	}
	o.orderQueue.book.addrMtx.RUnlock()

	data, err := json.Marshal(orders)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

type Fills struct {
	Time   string `json:"time"`
	Type   string `json:"type"`
	IsBuy  bool   `json:"isBuy"`
	Amount int64  `json:"amount"`
	Filled int64  `json:"filled"`
	Price  string `json:"price"`
}

func (o *OrderManager) FillsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])

	o.historyMtx.RLock()
	defer o.historyMtx.RUnlock()

	fills := make([]*Fills, 0)
	for _, t := range o.addrTradeHistory[addr] {
		fills = append(fills, &Fills{
			Time:   t.Time.Format(timeFormat),
			Type:   "market",
			IsBuy:  t.P1OrdersOrder.Flags[31] == 1,
			Amount: t.Amount.Int64(),
			Filled: t.Filled.Int64(),
			Price:  new(big.Int).Div(t.LimitPrice, exp18).String(),
		})
	}
	data, err := json.Marshal(fills)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

func (o *OrderManager) DropHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])
	amount, err := strconv.ParseInt(vars["amount"], 10, 64)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	token, err := contracts.NewTestToken(common.HexToAddress(Config.P1MarginAddress), o.engine.httpCli)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	privAdmin, _ := crypto.HexToECDSA(Config.PrivKeyHex)
	chainID, _ := new(big.Int).SetString(Config.ChainID, 10)
	adminTxOps, _ := bind.NewKeyedTransactorWithChainID(privAdmin, chainID)
	adminTxOps.GasLimit = 1000000

	tx, err := token.Mint(adminTxOps, addr, big.NewInt(amount))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, tx.Hash().String())
}
