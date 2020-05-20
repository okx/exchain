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

func (k Keeper) SetAccountAssetOnProduct(ctx sdk.Context, address sdk.AccAddress, product string, amt sdk.DecCoins, assetType int) {

	assetOnProduct, ok := k.GetAccountAssetOnProduct(ctx, address, product)
	// account info has exist
	if ok {
		switch assetType {
		case types.DepositType:
			assetOnProduct.Available = assetOnProduct.Available.Add(amt)
		case types.BorrowType:
			assetOnProduct.Available = assetOnProduct.Available.Add(amt)
			assetOnProduct.Borrowed = assetOnProduct.Borrowed.Add(amt)
		}
	} else {
		if assetType == types.DepositType {
			assetOnProduct = types.AccountAssetOnProduct{Product: product, Available: amt}
		}
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

// GetTradePair returns  the trade pair by product
func (k Keeper) GetTradePair(ctx sdk.Context, product string) *types.TradePair {
	var tradePair types.TradePair
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetTradePairKey(product))
	if bytes == nil {
		return nil
	}

	if k.cdc.UnmarshalBinaryBare(bytes, &tradePair) != nil {
		ctx.Logger().Error("decoding of token pair is failed", product)
		return nil
	}
	return &tradePair
}

func (k Keeper) SetTradePair(ctx sdk.Context, tradePair *types.TradePair) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTradePairKey(tradePair.Name)
	store.Set(key, k.cdc.MustMarshalBinaryBare(tradePair))
}

// Deposit deposits amount of tokens for a product
func (k Keeper) Deposit(ctx sdk.Context, address sdk.AccAddress, product string, amount sdk.DecCoins) sdk.Error {
	tradePair := k.GetTradePair(ctx, product)
	if tradePair == nil {
		tradePair = &types.TradePair{
			Owner:       address,
			Name:        product,
			Deposit:     amount,
			BlockHeight: ctx.BlockHeight(),
		}
	} else {
		tradePair.Deposit = tradePair.Deposit.Add(amount)
	}

	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, address, types.ModuleName, amount)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", amount.String()))
	}
	k.SetTradePair(ctx, tradePair)
	return nil
}
