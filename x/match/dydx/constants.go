package dydx

import "math/big"

type ORDER_FLAGS int

const (
	IS_BUY ORDER_FLAGS = 1 << iota
	IS_DECREASE_ONLY
	IS_NEGATIVE_LIMIT_FEE
)

var ONE_MINUTE_IN_SECONDS = big.NewInt(60)
var ONE_HOUR_IN_SECONDS = new(big.Int).Mul(ONE_MINUTE_IN_SECONDS, big.NewInt(60))
var ONE_DAY_IN_SECONDS = new(big.Int).Mul(ONE_HOUR_IN_SECONDS, big.NewInt(24))
var ONE_YEAR_IN_SECONDS = new(big.Int).Mul(ONE_DAY_IN_SECONDS, big.NewInt(365))
