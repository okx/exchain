package types

import (
	"math/big"
)

type OrderSide int

const (
	BuySide  OrderSide = 0
	SellSide OrderSide = 1
)

type Order interface {
	Symbol() string
	ID() string
	Side() OrderSide
	Price() *big.Float
	Amount() *big.Float
}
