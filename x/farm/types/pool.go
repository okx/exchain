package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// LockInfo is locked info of an address
type LockInfo struct {
	Addr             sdk.AccAddress `json:"addr"`
	PoolName         string         `json:"pool_name"`
	Amount           sdk.DecCoin    `json:"amount"`
	StartBlockHeight int64          `json:"start_block_height"`
}

// YieldingCoin is the token excluding native token which can be yielded
// by locking other tokens including LPT and token issued
type YieldingCoin struct {
	Coin                    sdk.DecCoin `json:"coin"`
	StartBlockHeightToYield int64       `json:"start_block_height_to_yield"`
	YieldAmountPerBlock     sdk.Dec     `json:"yield_amount_per_block"`
}

// String returns a human readable string representation of a YieldingCoin
func (yc YieldingCoin) String() string {
	return fmt.Sprintf(`YieldingCoin：
  Coin:								%s
  Start Block Height To Yield:		%d
  YieldAmountPerBlock:				%s`,
		yc.Coin, yc.StartBlockHeightToYield, yc.YieldAmountPerBlock)
}

// YieldingCoins is a collection of YieldingCoin
type YieldingCoins []YieldingCoin

// String returns a human readable string representation of YieldingCoins
func (ycs YieldingCoins) String() (out string) {
	for _, yc := range ycs {
		out += yc.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	PoolName          string `json:"pool_name"`
	LockedTokenSymbol string `json:"locked_token_symbol"`
	// sum of LockInfo.Amount
	TotalLockedCoin        sdk.DecCoin   `json:"total_locked_coin"`
	YieldingCoins          YieldingCoins `json:"yielding_coins"`
	YieldedCoins           sdk.DecCoins  `json:"yielded_coins"`
	LastYieldedBlockHeight int64         `json:"last_yielded_block_height"`
	// sum of (LockInfo.Amount * LockInfo.StartBlockHeight)
	TotalLockedWeight sdk.Dec `json:"total_locked_weight"`
}

// String returns a human readable string representation of FarmPool
func (fp FarmPool) String() string {
	return fmt.Sprintf(`FarmPool:	
  Pool Name:  				%s	
  Locked Token Symbol:      %s
  Total Locked Coin:		%s
  Yielding Coins:			%s
  Yielded Coins:			%s
  LastYieldedBlockHeight:	%d
  TotalLockedWeight:		%s`,
		fp.PoolName, fp.LockedTokenSymbol, fp.TotalLockedCoin, fp.YieldingCoins, fp.YieldedCoins,
		fp.LastYieldedBlockHeight, fp.TotalLockedWeight)
}
