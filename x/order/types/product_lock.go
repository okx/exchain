package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProductLock struct {
	BlockHeight  int64
	Price        sdk.Dec
	Quantity     sdk.Dec
	BuyExecuted  sdk.Dec
	SellExecuted sdk.Dec
}

type ProductLockMap struct {
	Data map[string]*ProductLock
}

func NewProductLockMap() *ProductLockMap {
	return &ProductLockMap{
		Data: make(map[string]*ProductLock),
	}
}
