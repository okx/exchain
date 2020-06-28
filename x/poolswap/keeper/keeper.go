package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/poolswap/types"
	tokentypes "github.com/okex/okchain/x/token/types"
)

// Keeper of the swap store
type Keeper struct {
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper

	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace types.ParamSubspace
}

// NewKeeper creates a swap keeper
func NewKeeper(supplyKeeper types.SupplyKeeper, tokenKeeper types.TokenKeeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		tokenKeeper:  tokenKeeper,
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetSwapTokenPair gets SwapTokenPair with quote token name
func (k Keeper) GetSwapTokenPair(ctx sdk.Context, tokenPairName string) (types.SwapTokenPair, error) {
	store := ctx.KVStore(k.storeKey)
	var item types.SwapTokenPair
	byteKey := types.GetTokenPairKey(tokenPairName)
	rawItem := store.Get(byteKey)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(rawItem, &item)
	if err != nil {
		return types.SwapTokenPair{}, err
	}

	return item, nil
}

// SetSwapTokenPair sets the entire SwapTokenPair data struct for a quote token name
func (k Keeper) SetSwapTokenPair(ctx sdk.Context, tokenPairName string, swapTokenPair types.SwapTokenPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(swapTokenPair)
	store.Set(types.GetTokenPairKey(tokenPairName), bz)
}

// DeleteSwapTokenPair deletes the entire SwapTokenPair data struct for a quote token name
func (k Keeper) DeleteSwapTokenPair(ctx sdk.Context, tokenPairName string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTokenPairKey(tokenPairName))
}

// GetSwapTokenPairsIterator get an iterator over all SwapTokenPairs in which the keys are the names and the values are the whois
func (k Keeper) GetSwapTokenPairsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.TokenPairPrefixKey)
}

// NewPoolToken new token
func (k Keeper) NewPoolToken(ctx sdk.Context, symbol string) {
	poolToken := types.InitPoolToken(symbol)
	k.tokenKeeper.NewToken(ctx, poolToken)
}

// GetPoolTokenInfo gets the token's info
func (k Keeper) GetPoolTokenInfo(ctx sdk.Context, symbol string) (tokentypes.Token, error) {
	poolToken := k.tokenKeeper.GetTokenInfo(ctx, symbol)
	if poolToken.Owner == nil {
		return poolToken, fmt.Errorf("Pool token %s does not exist", symbol)
	}
	return poolToken, nil
}

// GetPoolTokenAmount gets the amount of the specified poolToken name
func (k Keeper) GetPoolTokenAmount(ctx sdk.Context, poolTokenName string) (sdk.Dec, error) {
	poolToken, err := k.GetPoolTokenInfo(ctx, poolTokenName)
	return poolToken.TotalSupply, err
}

// MintPoolCoinsToUser mints coins and send them to the specified user address
func (k Keeper) MintPoolCoinsToUser(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

// BurnPoolCoinsFromUser sends coins to account module and burns them
func (k Keeper) BurnPoolCoinsFromUser(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return k.supplyKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

// SendCoinsToPool sends coins from user account to module account
func (k Keeper) SendCoinsToPool(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
}

// SendCoinsFromPoolToAccount sends coins from module account to user account
func (k Keeper) SendCoinsFromPoolToAccount(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

// nolint
func (k Keeper) GetTokenKeeper() types.TokenKeeper {
	return k.tokenKeeper
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
