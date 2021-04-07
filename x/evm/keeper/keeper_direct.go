package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm/types"
)

// SetCodeDirectly commit code into db with no cache
func (k Keeper) SetCodeDirectly(ctx sdk.Context, hash, code []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCode)
	store.Set(hash, code)
}

// SetStateDirectly commit one state into db with no cache
func (k Keeper) SetStateDirectly(ctx sdk.Context, addr ethcmn.Address, key, value ethcmn.Hash) {
	store, err := types.LoadAccountStorageStore(addr, sdk.CommitID{})
	if err != nil {
		panic(err)
	}
	store.Set(key.Bytes(), value.Bytes())
}
