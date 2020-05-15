package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dexTypes "github.com/okex/okchain/x/dex/types"

	"github.com/okex/okchain/x/margin/types"
	"github.com/okex/okchain/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the margin store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSubspace params.Subspace

	dexKeeper    types.DexKeeper
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper
}

// NewKeeper creates a margin keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace types.ParamSubspace, dexKeeper types.DexKeeper, tokenKeeper types.TokenKeeper, supplyKeeper types.SupplyKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSubspace: paramSubspace.WithKeyTable(types.ParamKeyTable()),

		dexKeeper:    dexKeeper,
		tokenKeeper:  tokenKeeper,
		supplyKeeper: supplyKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetMarginTradePair(ctx sdk.Context, product string) *dexTypes.TokenPair {
	// TODO : add margin token pair
	return k.dexKeeper.GetTokenPair(ctx, product)
}

func (k Keeper) GetAccountAssetOnProduct(ctx sdk.Context, addresses sdk.AccAddress, product string) (assetOnProduct types.AccountAssetOnProduct, ok bool) {
	bytes := ctx.KVStore(k.storeKey).Get(types.GetMarginProductAssetKey(addresses.String(), product))
	if bytes == nil {
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &assetOnProduct)
	return assetOnProduct, true
}

func (k Keeper) SetAccountAssetOnProduct(ctx sdk.Context, address sdk.AccAddress, product string, available sdk.DecCoins) {

	assetOnProduct, ok := k.GetAccountAssetOnProduct(ctx, address, product)
	if ok {
		assetOnProduct.Available = assetOnProduct.Available.Add(available)
	} else {
		assetOnProduct = types.AccountAssetOnProduct{Product: product, Available: available}
	}

	key := types.GetMarginProductAssetKey(address.String(), product)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(assetOnProduct)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

func (k Keeper) GetAccountDeposit(ctx sdk.Context, address sdk.AccAddress) (marginDeposit types.MarginProductAssets) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetMarginAllAssetKey(address.String()))
	defer iterator.Close()
	for i := int64(0); iterator.Valid(); iterator.Next() {
		var assetOnProduct types.AccountAssetOnProduct
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &assetOnProduct)
		marginDeposit = append(marginDeposit, assetOnProduct)
		i++
	}
	return
}

func (k Keeper) GetCDC() *codec.Codec {
	return k.cdc
}

// GetSupplyKeeper returns supply Keeper
func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

// GetSupplyKeeper returns token Keeper
func (k Keeper) GetTokenKeeper() types.TokenKeeper {
	return k.tokenKeeper
}

// GetDexKeeper returns dex Keeper
func (k Keeper) GetDexKeeper() types.DexKeeper {
	return k.dexKeeper
}

// GetParamSubspace returns paramSubspace
func (k Keeper) GetParamSubspace() params.Subspace {
	return k.paramSubspace
}

// SetParams sets inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.GetParamSubspace().SetParamSet(ctx, &params)
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.GetParamSubspace().GetParamSet(ctx, &params)
	return params
}

//// Get returns the pubkey from the adddress-pubkey relation
//func (k Keeper) Get(ctx sdk.Context, key string) (/* TODO: Fill out this type */, error) {
//	store := ctx.KVStore(k.storeKey)
//	var item /* TODO: Fill out this type */
//	byteKey := []byte(key)
//	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &item)
//	if err != nil {
//		return nil, err
//	}
//	return item, nil
//}
//
//func (k Keeper) set(ctx sdk.Context, key string, value /* TODO: fill out this type */ ) {
//	store := ctx.KVStore(k.storeKey)
//	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
//	store.Set([]byte(key), bz)
//}
//
//func (k Keeper) delete(ctx sdk.Context, key string) {
//	store := ctx.KVStore(k.storeKey)
//	store.Delete([]byte(key))
//}
