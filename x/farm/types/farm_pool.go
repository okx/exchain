package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	Owner             sdk.AccAddress    `json:"owner"`
	Name              string            `json:"name"`
	SymbolLocked      string            `json:"symbol_locked"`
	YieldedTokenInfos YieldedTokenInfos `json:"yielded_token_infos"`
	DepositAmount     sdk.DecCoin       `json:"deposit_amount"`
	// sum of LockInfo.Amount
	TotalValueLocked sdk.DecCoin  `json:"total_value_locked"`
	AmountYielded    sdk.DecCoins `json:"amount_yielded"`
}

// NewFarmPool creates a new instance of FarmPool
func NewFarmPool(name string, symbolLocked string, yieldedTokenInfos YieldedTokenInfos, totalValueLocked sdk.DecCoin,
	amountYielded sdk.DecCoins) FarmPool {
	return FarmPool{
		Name:              name,
		SymbolLocked:      symbolLocked,
		YieldedTokenInfos: yieldedTokenInfos,
		TotalValueLocked:  totalValueLocked,
		AmountYielded:     amountYielded,
	}
}

// CalculateAmountYieldedBetween is used for calculating how many tokens haven been yielding from LastClaimedBlockHeight to CurrentHeight
// Then transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
func (fp FarmPool) CalculateAmountYieldedBetween(currentHeight int64, startBlockHeight int64) (tokensYielded sdk.DecCoins) {
	for i := 0; i < len(fp.YieldedTokenInfos); i++ {
		startBlockHeightToYield := fp.YieldedTokenInfos[i].StartBlockHeightToYield
		if currentHeight > startBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if startBlockHeightToYield > startBlockHeight {
				blockInterval = sdk.NewDec(currentHeight - startBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(currentHeight - startBlockHeight)
			}

			var tokenYielded sdk.DecCoins
			// calculate how many coin have been yielded till the current block
			amount := blockInterval.MulTruncate(fp.YieldedTokenInfos[i].AmountYieldedPerBlock)
			remaining := fp.YieldedTokenInfos[i].RemainingAmount
			if amount.LT(remaining.Amount) {
				// add yielded amount
				tokenYielded = sdk.NewDecCoinsFromDec(remaining.Denom, amount)
				fp.AmountYielded = fp.AmountYielded.Add(tokenYielded)

				// subtract yielded_coin amount
				fp.YieldedTokenInfos[i].RemainingAmount.Amount = remaining.Amount.Sub(amount)
			} else {
				// add yielded amount
				tokenYielded = sdk.NewCoins(remaining)
				fp.AmountYielded = fp.AmountYielded.Add(tokenYielded)

				// initialize yieldedTokenInfo
				fp.YieldedTokenInfos[i] = NewYieldedTokenInfo(sdk.NewDecCoin(remaining.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())

				// TODO: remove the YieldedTokenInfo when its amount become zero
				// Currently, we support only one token of yield farming at the same time,
				// so, it is unnecessary to remove the element in slice
			}
			tokensYielded = tokensYielded.Add(tokenYielded)
		}
	}
	return
}

func (fp FarmPool) Finished() bool {
	for _, yieldedTokenInfo := range fp.YieldedTokenInfos {
		if yieldedTokenInfo.RemainingAmount.IsPositive() {
			return false
		}
	}
	return fp.TotalValueLocked.IsZero() && fp.AmountYielded.IsZero()
}

// String returns a human readable string representation of FarmPool
func (fp FarmPool) String() string {
	return fmt.Sprintf(`FarmPool:	
  Pool Name:  					    %s	
  Symbol Locked:      			    %s
  Yielded Token Infos:			    %s
  Total Value Locked:			    %s
  Amount Yielded Native Token:		%s`,
		fp.Name, fp.SymbolLocked, fp.YieldedTokenInfos, fp.TotalValueLocked, fp.AmountYielded)
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
