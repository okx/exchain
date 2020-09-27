package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LockInfo is locked info of an address
type LockInfo struct {
	Owner            sdk.AccAddress `json:"owner"`
	PoolName         string         `json:"pool_name"`
	Amount           sdk.DecCoin    `json:"amount"`
	StartBlockHeight int64          `json:"start_block_height"`
}

// YieldedToken is the token excluding native token which can be yielded
// by locking other tokens including LPT and token issued
type YieldedTokenInfo struct {
	TotalAmount             sdk.DecCoin `json:"total_amount"`
	StartBlockHeightToYield int64       `json:"start_block_height_to_yield"`
	AmountYieldedPerBlock   sdk.Dec     `json:"amount_yielded_per_block"`
}

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	Name              string `json:"name"`
	SymbolLocked      string `json:"symbol_locked"`
	YieldedTokenInfos      []YieldedTokenInfo `json:"yieldied_token_infos"`

	// sum of LockInfo.Amount
	TotalValueLocked       sdk.DecCoin        `json:"total_value_locked"`
	AmountYielded          sdk.DecCoins       `json:"amount_yielded"`
	LastUpdatedBlockHeight int64              `json:"last_updated_block_height"`
	// sum of (LockInfo.Amount * LockInfo.StartBlockHeight)
	TotalLockedInfo sdk.Dec `json:"total_locked_info"`
}

type RunTimePoolInfo struct {
	// sum of LockInfo.Amount
	TotalValueLocked       sdk.DecCoin        `json:"total_value_locked"`
	AmountYielded          sdk.DecCoins       `json:"amount_yielded"`
	LastUpdatedBlockHeight int64              `json:"last_updated_block_height"`
	// sum of (LockInfo.Amount * LockInfo.StartBlockHeight)
	TotalWeight sdk.Dec `json:"total_locked_info"`
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

// FarmPools is a collection of FarmPool
type FarmPools []FarmPool

// String returns a human readable string representation of FarmPools
func (fps FarmPools) String() (out string) {
	for _, fp := range fps {
		out += fp.String() + "\n"
	}
	return strings.TrimSpace(out)
}
