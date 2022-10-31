package dydx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
)

const (
	POST = "POST"
	GET  = "GET"

	addrKey    = "addr"
	timeFormat = "15:00:00"
)

type Response struct {
	Succeed  bool   `json:"succeed"`
	ErrorMsg string `json:"errorMsg"`
}

func (o *OrderManager) ServeWeb() {
	r := mux.NewRouter()
	r.HandleFunc("/", EmptyHandler)
	r.HandleFunc("/placeorder", o.PlaceOrderHandler).Methods(POST)
	r.HandleFunc("/send", o.SendHandler).Methods(POST)

	r.HandleFunc("/book", o.BookHandler).Methods(GET)
	r.HandleFunc("/trades", o.TradesHandler).Methods(GET)
	r.HandleFunc("/position/{addr}", o.PositionHandler).Methods(GET)
	r.HandleFunc("/self-orders/{addr}", o.OrdersHandler).Methods(GET)
	r.HandleFunc("/self-fills/{addr}", o.FillsHandler).Methods(GET)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

func EmptyHandler(w http.ResponseWriter, r *http.Request) {}

func (o *OrderManager) PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {

}

func (o *OrderManager) SendHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hexSignedOrder := vars["signedOrder"]
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
	levels := bookToLevel(o.book)
	data, err := json.Marshal(levels)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

type Trade struct {
	Size  int64  `json:"size"`
	Price int64  `json:"price"`
	Time  string `json:"time"`
}

func (o *OrderManager) TradesHandler(w http.ResponseWriter, r *http.Request) {
	o.historyMtx.RLock()
	defer o.historyMtx.RUnlock()

	trades := make([]*Trade, 0)
	for _, t := range o.tradeHistory {
		fmt.Println("trade history", *t)
		trades = append(trades, &Trade{
			Size:  t.Amount.Int64(),
			Price: t.LimitPrice.Int64(),
			Time:  t.Time.Format(timeFormat),
		})
	}
	data, err := json.Marshal(trades)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}

type Balance struct {
	Margin   *big.Int `json:"margin"`
	Position *big.Int `json:"position"`
}

func (o *OrderManager) PositionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])
	p1Balance, err := o.engine.contracts.PerpetualV1.GetAccountBalance(nil, addr)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	data, err := json.Marshal(&Balance{
		Margin:   p1Balance.Margin,
		Position: p1Balance.Position,
	})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))

}

func (o *OrderManager) OrdersHandler(w http.ResponseWriter, r *http.Request) {

}

func (o *OrderManager) FillsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := common.HexToAddress(vars[addrKey])

	o.historyMtx.RLock()
	defer o.historyMtx.RUnlock()

	var trades []*Trade
	for _, t := range o.tradeHistory {
		if t.Maker != addr {
			continue
		}
		trades = append(trades, &Trade{
			Size:  t.Amount.Int64(),
			Price: t.LimitPrice.Int64(),
			Time:  t.Time.Format(timeFormat),
		})
	}
	data, err := json.Marshal(trades)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(data))
}
