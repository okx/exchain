package keeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// DeleteStateDirectly commit one state into db with no cache
func (k Keeper) DeleteStateDirectly(ctx sdk.Context, addr ethcmn.Address, key ethcmn.Hash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AddressStoragePrefix(addr)) // todo: instead of global store
	store.Delete(key.Bytes())
}