package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Earnings - structure for a earning query
type Earnings struct {
	TargetBlockHeight int64        `json:"target_block_height"`
	AmountLocked      sdk.DecCoin  `json:"amount_locked"`
	AmountYielded     sdk.DecCoins `json:"amount_yielded"`
}

// NewEarnings creates a new instance of Earnings
func NewEarnings(targetBlockHeight int64, amountLocked sdk.DecCoin, amountYielded sdk.DecCoins) Earnings {
	return Earnings{
		TargetBlockHeight: targetBlockHeight,
		AmountLocked:      amountLocked,
		AmountYielded:     amountYielded,
	}
}

// String returns a human readable string representation of Earnings
func (e Earnings) String() string {
	return fmt.Sprintf(`Earnings:
  Target Block Height: 		%d,
  Amount Locked:			%s,
  Amount Yielded:			%s`,
		e.TargetBlockHeight, e.AmountLocked, e.AmountYielded,
	)
}
