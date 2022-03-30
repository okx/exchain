package blockfee

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Pool interface {
	AddFeeBlockPool(amt sdk.Coins)
	GetFeeFromBlockPool() sdk.Coins
	ResetFeeBlockPool()
}

type collector struct {
	coins sdk.Coins
}

func NewCollector() *collector {
	return &collector{}
}

func (c *collector) ResetFeeBlockPool() {
	c.coins = sdk.NewCoins()
}

func (c *collector) AddFeeBlockPool(amt sdk.Coins) {

	c.coins = c.coins.Add2(amt)
}

func (c *collector) GetFeeFromBlockPool() sdk.Coins {
	return c.coins
}
