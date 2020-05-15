package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TradePair struct {
	Name        string       `json:"name"`
	Deposit     sdk.DecCoins `json:"deposit"`
	MaxLeverage int64        `json:"max-leverage"`
	BorrowRate  float64      `json:"borrow-rate"`
	RiskRate    float64      `json:"risk-rate"`
	BlockHeight int64        `json:"block_height"`
}

func DefaultTradePair() *TradePair {
	return &TradePair{}
}
