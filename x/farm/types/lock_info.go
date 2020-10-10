package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// LockInfo is locked info of an address
type LockInfo struct {
	Owner            sdk.AccAddress `json:"owner"`
	PoolName         string         `json:"pool_name"`
	Amount           sdk.DecCoin    `json:"amount"`
	StartBlockHeight int64          `json:"start_block_height"`
	ReferencePeriod  uint64         `json:"reference_period"`
}

// NewLockInfo creates a new instance of LockInfo
func NewLockInfo(owner sdk.AccAddress, poolName string, amount sdk.DecCoin, startBlockHeight int64, referencePeriod uint64) LockInfo {
	return LockInfo{
		Owner:            owner,
		PoolName:         poolName,
		Amount:           amount,
		StartBlockHeight: startBlockHeight,
		ReferencePeriod:  referencePeriod,
	}
}
