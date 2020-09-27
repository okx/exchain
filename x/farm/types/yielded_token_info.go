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

// NewYieldedTokenInfo creates a new instance of YieldedTokenInfo
func NewYieldedTokenInfo(
	totalAmount sdk.DecCoin, startBlockHeightToYield int64, amountYieldedPerBlock sdk.Dec,
) YieldedTokenInfo {
	return YieldedTokenInfo{
		TotalAmount:             totalAmount,
		StartBlockHeightToYield: startBlockHeightToYield,
		AmountYieldedPerBlock:   amountYieldedPerBlock,
	}
}

// String returns a human readable string representation of a YieldedTokenInfo
func (yti YieldedTokenInfo) String() string {
	return fmt.Sprintf(`YieldedTokenInfo：
  Coin:								%s
  Start Block Height To Yield:		%d
  AmountYieldedPerBlock:				%s`,
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
