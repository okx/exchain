package fee

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Fee interface {
	AddFee(amt sdk.Coins)
	GetFee() sdk.Coins
	ResetFee()
}

type collector struct {
	coins sdk.Coins
}

func NewCollector() *collector {
	return &collector{}
}

func (c *collector) ResetFee() {
	c.coins = sdk.NewCoins()
}

func (c *collector) AddFee(amt sdk.Coins) {
	c.coins = c.coins.Add2(amt)
}

func (c *collector) GetFee() sdk.Coins {
	return c.coins
}
