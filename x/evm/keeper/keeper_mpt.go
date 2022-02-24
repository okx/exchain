package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	FlagContractStateCache = "contract-state-cache"
)

var (
	ContractStateCache uint = 2048 // MB
)

func (k *Keeper) Commit(ctx sdk.Context) {
	// commit contract storage mpt trie
	k.EvmStateDb.WithContext(ctx).Commit(true)
}
