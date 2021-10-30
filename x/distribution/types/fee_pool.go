package types

import (
	"fmt"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// FeePool is the struct of the global fee pool for distribution
type FeePool struct {
	CommunityPool sdk.SysCoins `json:"community_pool" yaml:"community_pool"` // pool for community funds yet to be spent
}

// InitialFeePool zero fee pool
func InitialFeePool() FeePool {
	return FeePool{
		CommunityPool: sdk.SysCoins{},
	}
}

// ValidateGenesis validates the fee pool for a genesis state
func (f FeePool) ValidateGenesis() error {
	if f.CommunityPool.IsAnyNegative() {
		return fmt.Errorf("negative CommunityPool in distribution fee pool, is %v",
			f.CommunityPool)
	}
	return nil
}
