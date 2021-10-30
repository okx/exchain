package types

import (
	"fmt"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// LockInfo is locked info of an address
type LockInfo struct {
	Owner            sdk.AccAddress `json:"owner"`
	PoolName         string         `json:"pool_name"`
	Amount           sdk.SysCoin    `json:"amount"`
	StartBlockHeight int64          `json:"start_block_height"`
	ReferencePeriod  uint64         `json:"reference_period"`
}

// NewLockInfo creates a new instance of LockInfo
func NewLockInfo(owner sdk.AccAddress, poolName string, amount sdk.SysCoin, startBlockHeight int64, referencePeriod uint64) LockInfo {
	return LockInfo{
		Owner:            owner,
		PoolName:         poolName,
		Amount:           amount,
		StartBlockHeight: startBlockHeight,
		ReferencePeriod:  referencePeriod,
	}
}

// String returns a human readable string representation of LockInfo
func (li LockInfo) String() string {
	return fmt.Sprintf(`Lock Info:
  Owner:						%s	
  Pool Name:					%s
  Locked Amount:      			%s
  Start Block Height:           %d
  Reference Period:             %d`,
		li.Owner, li.PoolName, li.Amount, li.StartBlockHeight, li.ReferencePeriod)
}
