package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// YieldedTokenInfo is the token excluding native token which can be yielded by locking other tokens including LPT and
// token issued
type YieldedTokenInfo struct {
	TotalAmount             sdk.DecCoin `json:"total_amount"`
	StartBlockHeightToYield int64       `json:"start_block_height_to_yield"`
	AmountYieldedPerBlock   sdk.Dec     `json:"amount_yielded_per_block"`
}

// String returns a human readable string representation of a YieldedTokenInfo
func (yti YieldedTokenInfo) String() string {
	return fmt.Sprintf(`YieldingCoin：
  Coin:								%s
  Start Block Height To Yield:		%d
  YieldAmountPerBlock:				%s`,
		yti.TotalAmount, yti.StartBlockHeightToYield, yti.AmountYieldedPerBlock)
}

// YieldedTokenInfos is a collection of YieldedTokenInfo
type YieldedTokenInfos []YieldedTokenInfo

// String returns a human readable string representation of YieldedTokenInfos
func (ytis YieldedTokenInfos) String() (out string) {
	for _, yti := range ytis {
		out += yti.String() + "\n"
	}
	return strings.TrimSpace(out)
}
