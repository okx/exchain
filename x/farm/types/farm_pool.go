package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FarmPool is the pool where an address can lock specified token to yield other tokens
type FarmPool struct {
	Owner         sdk.AccAddress `json:"owner"`
	Name          string         `json:"name"`
	SymbolLocked  string         `json:"symbol_locked"`
	DepositAmount sdk.DecCoin    `json:"deposit_amount"`
	// sum of LockInfo.Amount
	TotalValueLocked  sdk.DecCoin       `json:"total_value_locked"`
	YieldedTokenInfos YieldedTokenInfos `json:"yielded_token_infos"`
}

// NewFarmPool creates a new instance of FarmPool
func NewFarmPool(name string, symbolLocked string, yieldedTokenInfos YieldedTokenInfos, totalValueLocked sdk.DecCoin) FarmPool {
	return FarmPool{
		Name:              name,
		SymbolLocked:      symbolLocked,
		YieldedTokenInfos: yieldedTokenInfos,
		TotalValueLocked:  totalValueLocked,
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
  Symbol Locked:      			    %s
  Yielded Token Infos:			    %s
  Total Value Locked:			    %s`,
		fp.Name, fp.SymbolLocked, fp.YieldedTokenInfos, fp.TotalValueLocked)
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
