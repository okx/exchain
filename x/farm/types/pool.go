package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	PoolName          string `json:"pool_name"`
	LockedTokenSymbol string `json:"locked_token_symbol"`
	// sum of LockInfo.Amount
	TotalLockedCoin        sdk.DecCoin    `json:"total_locked_coin"`
	YieldingCoins          []YieldingCoin `json:"yielding_coins"`
	YieldedCoins           sdk.DecCoins   `json:"yielded_coins"`
	LastBlockHeightToYield int64          `json:"last_block_height_to_yield"`
	// sum of (LockInfo.Amount * LockInfo.StartBlockHeight)
	TotalLockedInfo sdk.Dec `json:"total_locked_info"`
}
