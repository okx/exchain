package blockfee

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"sync"
)

type Pool interface {
	AddFee(amt sdk.Coins)
	SubFee(amt sdk.Coins)
	Get() sdk.Coins
	Reset()
}

type collector struct {
	mu    sync.Mutex
	coins sdk.Coins
}

func NewCollector() *collector {
	return &collector{}
}

func (c *collector) Reset() {
	c.coins = sdk.NewCoins()
}

func (c *collector) AddFee(amt sdk.Coins) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.coins = c.coins.Add2(amt)
}

func (c *collector) SubFee(amt sdk.Coins) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.coins = c.coins.Sub(amt)
}

func (c *collector) Get() sdk.Coins {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.coins
}
