package dydx

import (
	"fmt"
	"math/big"
)

type MatchEngine struct {
	depthBook *DepthBook
}

type MatchResult struct {
	MatchedRecords []*MatchRecord
	TakerOrder     *WrapOrder
}

func (r *MatchResult) AddMatchedRecord(price *big.Int, amount *big.Int, makerOrder *WrapOrder) {
	r.MatchedRecords = append(r.MatchedRecords, &MatchRecord{
		Price:  price,
		Amount: amount,
		Maker:  makerOrder,
		Taker:  r.TakerOrder,
	})
}

type MatchRecord struct {
	Price  *big.Int
	Amount *big.Int
	Taker  *WrapOrder
	Maker  *WrapOrder
}

func (m *MatchEngine) Match(order *WrapOrder) (*MatchResult, error) {
	if order.Type() == BuyOrderType {
		return processOrder(order, m.depthBook.sellOrders, m.depthBook.buyOrders)
	} else if order.Type() == SellOrderType {
		return processOrder(order, m.depthBook.buyOrders, m.depthBook.sellOrders)
	} else {
		return nil, fmt.Errorf("invalid order type")
	}
}

func processOrder(takerOrder *WrapOrder, makerBook *OrderList, takerBook *OrderList) (*MatchResult, error) {
	var matchResult = &MatchResult{
		TakerOrder: takerOrder,
	}
	for {
		makerOrderElem := makerBook.Front()
		if makerOrderElem == nil {
			break
		}
		makerOrder := makerOrderElem.Value.(*WrapOrder)
		if takerOrder.Type() == BuyOrderType && takerOrder.Price().Cmp(makerOrder.Price()) < 0 {
			break
		}
		if takerOrder.Type() == SellOrderType && takerOrder.Price().Cmp(makerOrder.Price()) > 0 {
			break
		}
		//marketPrice := makerOrder.Price()
		//matchAmount := takerOrder.Amount()
		//if takerOrder.Amount().Cmp(makerOrder.Amount()) > 0 {
		//	matchAmount = makerOrder.Amount()
		//}
		//matchResult.AddMatchedRecord(marketPrice, matchAmount, makerOrder)
		//takerOrder.SubAmount(matchAmount)
		//makerOrder.SubAmount(matchAmount)
		//if makerOrder.Amount().Cmp(big.NewInt(0)) == 0 {
		//	makerBook.Remove(makerOrderElem)
		//}
		//if takerOrder.Amount().Cmp(big.NewInt(0)) == 0 {
		//	break
		//}
	}
	//if takerOrder.Amount.Cmp(big.NewInt(0)) > 0 {
	//	takerBook.Insert(takerOrder)
	//}
	return matchResult, nil
}
