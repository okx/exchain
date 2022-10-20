package dydx

import (
	"fmt"
	"net/http"
)

var book *DepthBook

func (d *OrderManager) Serve() {
	book = d.book
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe("127.0.0.1:8555", nil)
	if err != nil {
		panic(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, " buy orders count: %d\n", book.buyOrders.Len())
		fmt.Fprintf(w, "sell orders count: %d\n", book.sellOrders.Len())

	case "/buy":
		fmt.Fprintf(w, "total orders: %d\n", book.buyOrders.Len())
		for _, order := range book.buyOrders.List() {
			fmt.Fprintf(w, "orederHash: %s, amount: %d, left: %d, frozen: %d\n", order.Hash(), order.Amount, order.LeftAmount, order.FrozenAmount)
		}

	case "/sell":
		fmt.Fprintf(w, "total orders: %d\n", book.sellOrders.Len())
		for _, order := range book.sellOrders.List() {
			fmt.Fprintf(w, "orederHash: %s, amount: %d, left: %d, frozen: %d\n", order.Hash(), order.Amount, order.LeftAmount, order.FrozenAmount)
		}
	default:
		fmt.Fprintf(w, "Invalid path")
	}
}
