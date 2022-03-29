package feecollector

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"sync"
)

type FeeCollector interface {
	Add(amt sdk.Coins)
	Sub(amt sdk.Coins)
	Get() sdk.Coins
}

type feeCollectorCoins struct {
	mu          sync.Mutex
	cachedCoins sdk.Coins
}

func NewFeeCollectorCoins() *feeCollectorCoins {
	return &feeCollectorCoins{
		cachedCoins: sdk.NewCoins(),
	}
}

func (fc *feeCollectorCoins) Add(amt sdk.Coins) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.cachedCoins = fc.cachedCoins.Add2(amt)
}

func (fc *feeCollectorCoins) Sub(amt sdk.Coins) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.cachedCoins = fc.cachedCoins.Sub(amt)
}

func (fc *feeCollectorCoins) Get() sdk.Coins {
	return fc.cachedCoins
}
