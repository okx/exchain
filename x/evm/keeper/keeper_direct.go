package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/evm/types"
)

// SetCodeDirectly commit code into db with no cache111
func (k Keeper) SetCodeDirectly(ctx sdk.Context, hash, code []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCode)
	store.Set(hash, code)
}

// SetStateDirectly commit one state into db with no cache111
func (k Keeper) SetStateDirectly(ctx sdk.Context, addr ethcmn.Address, key, value ethcmn.Hash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AddressStoragePrefix(addr))
	store.Set(key.Bytes(), value.Bytes())
}
