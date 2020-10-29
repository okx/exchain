package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	Owner           sdk.AccAddress `json:"owner"`
	Name            string         `json:"name"`
	MinLockedAmount sdk.DecCoin    `json:"min_locked_amount"`
	DepositAmount   sdk.DecCoin    `json:"deposit_amount"`
	// sum of LockInfo.Amount
	TotalValueLocked        sdk.DecCoin       `json:"total_value_locked"`
	YieldedTokenInfos       YieldedTokenInfos `json:"yielded_token_infos"`
	TotalAccumulatedRewards sdk.DecCoins      `json:"total_accumulated_rewards"`
}

// NewFarmPool creates a new instance of FarmPool
func NewFarmPool(
	owner sdk.AccAddress, name string, lockedSymbol, depositAmt, totalValueLocked sdk.DecCoin,
	yieldedTokenInfos YieldedTokenInfos, accumulatedRewards sdk.DecCoins,
) FarmPool {
	return FarmPool{
		Owner:                   owner,
		Name:                    name,
		MinLockedAmount:         lockedSymbol,
		DepositAmount:           depositAmt,
		TotalValueLocked:        totalValueLocked,
		YieldedTokenInfos:       yieldedTokenInfos,
		TotalAccumulatedRewards: accumulatedRewards,
	}
}

func (fp FarmPool) Finished() bool {
	for _, yieldedTokenInfo := range fp.YieldedTokenInfos {
		if yieldedTokenInfo.RemainingAmount.IsPositive() {
			return false
		}
	}
	return fp.TotalValueLocked.IsZero()
}

// String returns a human readable string representation of FarmPool
func (fp FarmPool) String() string {
	return fmt.Sprintf(`FarmPool:
  Pool Name:  					    %s	
  Owner:							%s
  Locked Symbol:      			    %s
  Deposit Amount:                   %s
  Total Value Locked:               %s
  Yielded Token Infos:			    %s
  Total Accumulated Rewards:        %s`,
		fp.Name, fp.Owner, fp.MinLockedAmount, fp.DepositAmount, fp.TotalValueLocked, fp.YieldedTokenInfos, fp.TotalAccumulatedRewards)
}

// FarmPools is a collection of FarmPool
type FarmPools []FarmPool

// String returns a human readable string representation of FarmPools
func (fps FarmPools) String() (out string) {
	for _, fp := range fps {
		out += fp.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// PoolNum is a wrapped structure of uint to display by cli query
type PoolNum struct {
	Number uint `json:"number"`
}

// NewPoolNum creates a new instance of PoolNum
func NewPoolNum(num uint) PoolNum {
	return PoolNum{
		Number: num,
	}
}

// String returns a human readable string representation of PoolNum
func (pn PoolNum) String() string {
	return fmt.Sprintf(`Number Of Pools:
  Number: 		%d`, pn.Number)
}
