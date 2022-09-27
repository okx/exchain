package dydx

type Perpetual struct {
	contracts *Contracts
	orders    *Orders
}

func (p *Perpetual) NewTradeOperation() *TradeOperation {
	return &TradeOperation{
		contracts: p.contracts,
		orders:    p.orders,
	}
}
