package dydx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var book *DepthBook

type OrderShowList struct {
	Price  string
	Amount uint64
}

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
		//fmt.Fprintf(w, "total orders: %d\n", book.buyOrders.Len())
		m := make(map[uint64]uint64)
		list := []OrderShowList{{"0", 0}}
		for _, order := range book.buyOrders.List() {
			m[order.LimitPrice.Uint64()] += order.GetLeftAmount().Uint64()
			if order.GetLimitPrice().String() == list[len(list)-1].Price {
				list[len(list)-1].Amount += order.GetLeftAmount().Uint64()
			} else {
				list = append(list, OrderShowList{
					Price:  order.GetLimitPrice().String(),
					Amount: order.GetLeftAmount().Uint64(),
				})
			}
			//fmt.Fprintf(w, "orederHash: %s, amount: %d, left: %d, frozen: %d\n", order.Hash(), order.Amount, order.LeftAmount, order.FrozenAmount)
		}
		//data, err := json.MarshalIndent(m, "", "    ")
		//if err != nil {
		//	fmt.Fprintf(w, err.Error())
		//} else {
		//	fmt.Fprintf(w, string(data))
		//}
		//fmt.Fprintf(w, "\n\n")
		list = list[1:]
		data2, err := json.MarshalIndent(list, "", "    ")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, string(data2))
		}

	case "/sell":
		//fmt.Fprintf(w, "total orders: %d\n", book.sellOrders.Len())
		m := make(map[uint64]uint64)
		list := []OrderShowList{{"0", 0}}
		for _, order := range book.sellOrders.List() {
			m[order.LimitPrice.Uint64()] += order.GetLeftAmount().Uint64()
			if order.GetLimitPrice().String() == list[len(list)-1].Price {
				list[len(list)-1].Amount += order.GetLeftAmount().Uint64()
			} else {
				list = append(list, OrderShowList{
					Price:  order.GetLimitPrice().String(),
					Amount: order.GetLeftAmount().Uint64(),
				})
			}
			//fmt.Fprintf(w, "orederHash: %s, amount: %d, left: %d, frozen: %d\n", order.Hash(), order.Amount, order.LeftAmount, order.FrozenAmount)
		}
		//data, err := json.MarshalIndent(m, "", "    ")
		//if err != nil {
		//	fmt.Fprintf(w, err.Error())
		//} else {
		//	fmt.Fprintf(w, string(data))
		//}
		//fmt.Fprintf(w, "\n\n")
		list = list[1:]
		data2, err := json.MarshalIndent(list, "", "    ")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, string(data2))
		}
	default:
		fmt.Fprintf(w, "Invalid path")
	}
}
