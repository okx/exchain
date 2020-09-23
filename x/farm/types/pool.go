package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type YieldOKTWhiteList []string

type FarmPool struct {
	PoolName          string `json:"pool_name"`
	LockedTokenSymbol string `json:"locked_token_symbol"`
	// S1
	TotalLockedToken       sdk.DecCoin  `json:"total_locked_token"`
	YieldingTokens         sdk.DecCoins `json:"releasing_tokens_holder"`
	YieldedTokens          sdk.DecCoins `json:"yielded_tokens_holder"`
	LastBlockHeightToYield int64        `json:"last_block_height_to_yield"`
	YieldAmountPerBlock    sdk.Dec      `json:"yield_amount_per_block"`
	// S2
	TotalLockedInfo sdk.Dec `json:"total_locked_info"`
}
