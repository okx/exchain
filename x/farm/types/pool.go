package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type YieldOKTWhiteList []string

type FarmPool struct {
	PoolName          string `json:"pool_name"`
	LockedTokenSymbol string `json:"locked_token_symbol"`
	// sum of all lockedAmount
	TotalLockedCoin       sdk.DecCoin  `json:"total_locked_coin"`
	YieldingCoins         sdk.DecCoins `json:"yielding_coins"`
	YieldedCoins          sdk.DecCoins `json:"yielded_coins"`
	LastBlockHeightToYield int64        `json:"last_block_height_to_yield"`
	YieldAmountPerBlock    sdk.Dec      `json:"yield_amount_per_block"`
	// sum of all lockedAmount * lockedBlockHeight
	TotalLockedInfo sdk.Dec `json:"total_locked_info"`
}
