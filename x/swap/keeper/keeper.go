package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/swap/types"
)

// Keeper of the swap store
type Keeper struct {
	bankKeeper   bank.Keeper
	supplyKeeper supply.Keeper

	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a swap keeper
func NewKeeper(bankKeeper bank.Keeper, supplyKeeper supply.Keeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		bankKeeper:   bankKeeper,
		supplyKeeper: supplyKeeper,
		storeKey:     key,
		cdc:          cdc,
		paramspace:   paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get SwapTokenPair with quote token name
func (k Keeper) GetSwapTokenPair(ctx sdk.Context, quote string) (types.SwapTokenPair, error) {
	store := ctx.KVStore(k.storeKey)
	var item types.SwapTokenPair
	byteKey := []byte(quote)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &item)
	if err != nil {
		return types.SwapTokenPair{}, err
	}
	return item, nil
}

// Sets the entire SwapTokenPair data struct for a quote token name
func (k Keeper) SetSwapTokenPair(ctx sdk.Context, quote string, swapTokenPair types.SwapTokenPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(swapTokenPair)
	store.Set([]byte(quote), bz)
}

// Deletes the entire SwapTokenPair data struct for a quote token name
func (k Keeper) DeleteSwapTokenPair(ctx sdk.Context, quote string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(quote))
}

// Get an iterator over all SwapTokenPairs in which the keys are the names and the values are the whois
func (k Keeper) GetSwapTokenPairsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
