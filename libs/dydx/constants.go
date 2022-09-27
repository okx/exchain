package dydx

import "math/big"

type ORDER_FLAGS int

const (
	IS_BUY ORDER_FLAGS = 1 << iota
	IS_DECREASE_ONLY
	IS_NEGATIVE_LIMIT_FEE
)

const (
	DEFAULT_EIP712_DOMAIN_NAME    = "P1Orders"
	DEFAULT_EIP712_DOMAIN_VERSION = "1.0"
	EIP712_ORDER_STRUCT_STRING    = "Order(" +
		"bytes32 flags," +
		"uint256 amount," +
		"uint256 limitPrice," +
		"uint256 triggerPrice," +
		"uint256 limitFee," +
		"address maker," +
		"address taker," +
		"uint256 expiration" +
		")"
	EIP712_CANCEL_ORDER_STRUCT_STRING = "CancelLimitOrder(" +
		"string action," +
		"bytes32[] orderHashes" +
		")"
)

var ONE_MINUTE_IN_SECONDS = big.NewInt(60)
var ONE_HOUR_IN_SECONDS = new(big.Int).Mul(ONE_MINUTE_IN_SECONDS, big.NewInt(60))
var ONE_DAY_IN_SECONDS = new(big.Int).Mul(ONE_HOUR_IN_SECONDS, big.NewInt(24))
var ONE_YEAR_IN_SECONDS = new(big.Int).Mul(ONE_DAY_IN_SECONDS, big.NewInt(365))
