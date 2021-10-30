package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// nolint
type ProductLock struct {
	BlockHeight  int64
	Price        sdk.Dec
	Quantity     sdk.Dec
	BuyExecuted  sdk.Dec
	SellExecuted sdk.Dec
}

// nolint
type ProductLockMap struct {
	Data map[string]*ProductLock
}

// nolint
func NewProductLockMap() *ProductLockMap {
	return &ProductLockMap{
		Data: make(map[string]*ProductLock),
	}
}
