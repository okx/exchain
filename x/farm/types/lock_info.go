package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// LockInfo is locked info of an address
type LockInfo struct {
	Owner            sdk.AccAddress `json:"owner"`
	PoolName         string         `json:"pool_name"`
	Amount           sdk.DecCoin    `json:"amount"`
	StartBlockHeight int64          `json:"start_block_height"`
}
