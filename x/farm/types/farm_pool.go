package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	Name              string             `json:"name"`
	SymbolLocked      string             `json:"symbol_locked"`
	YieldedTokenInfos []YieldedTokenInfo `json:"yieldied_token_infos"`

	// sum of LockInfo.Amount
	TotalValueLocked       sdk.DecCoin  `json:"total_value_locked"`
	AmountYielded          sdk.DecCoins `json:"amount_yielded"`
	LastClaimedBlockHeight int64        `json:"last_claimed_block_height"`
	// sum of (LockInfo.Amount * LockInfo.StartBlockHeight)
	TotalLockedWeight sdk.Dec `json:"total_locked_Weight"`
}

// String returns a human readable string representation of FarmPool
func (fp FarmPool) String() string {
	return fmt.Sprintf(`FarmPool:	
  Pool Name:  					%s	
  Symbol Locked:      			%s
  Yielded Token Infos:			%s
  Total Value Locked:			%s
  Amount Yielded:				%s
  Last Claimed Block Height:	%d
  Total Locked Weight:			%s`,
		fp.Name, fp.SymbolLocked, fp.YieldedTokenInfos, fp.TotalValueLocked, fp.AmountYielded,
		fp.LastClaimedBlockHeight, fp.TotalLockedWeight)
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
