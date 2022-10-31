package dydx

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	POST = "POST"
	GET  = "GET"

	addrKey = "addr"
)

type ContractCache struct {
	Balances map[string]int
}

var contractCache = ContractCache{
	Balances: make(map[string]int),
}

type Response struct {
	Succeed  bool   `json:"succeed"`
	ErrorMsg string `json:"errorMsg"`
}

func (o *OrderManager) ServeWeb() {
	r := mux.NewRouter()
	r.HandleFunc("/", EmptyHandler)
	r.HandleFunc("/placeorder", o.PlaceOrderHandler).Methods(POST)
	r.HandleFunc("/send", o.SendHandler).Methods(POST)

	r.HandleFunc("/orders", o.OrdersHandler).Methods(GET)
	r.HandleFunc("/trades", o.TradesHandler).Methods(GET)
	r.HandleFunc("/position/{addr}", o.PositionHandler).Methods(GET)
	r.HandleFunc("/self-orders/{addr}", o.SelfOrdersHandler).Methods(GET)
	r.HandleFunc("/self-fills/{addr}", o.SelfFillsHandler).Methods(GET)

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

func (o *OrderManager) OrdersHandler(w http.ResponseWriter, r *http.Request) {}

func (o *OrderManager) TradesHandler(w http.ResponseWriter, r *http.Request) {}

func (o *OrderManager) PositionHandler(w http.ResponseWriter, r *http.Request) {}

func (o *OrderManager) SelfOrdersHandler(w http.ResponseWriter, r *http.Request) {}

func (o *OrderManager) SelfFillsHandler(w http.ResponseWriter, r *http.Request) {}
