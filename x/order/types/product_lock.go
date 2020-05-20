package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
type ProductLock struct {
	BlockHeight  int64 `json:"block_height"`
	Price        sdk.Dec `json:"price"`
	Quantity     sdk.Dec `json:"quantity"`
	BuyExecuted  sdk.Dec `json:"buy_executed"`
	SellExecuted sdk.Dec `json:"sell_executed"`
}

// nolint
type ProductLockMap struct {
	Data map[string]*ProductLock `json:"data"`
}

// nolint
func NewProductLockMap() *ProductLockMap {
	return &ProductLockMap{
		Data: make(map[string]*ProductLock),
	}
}
