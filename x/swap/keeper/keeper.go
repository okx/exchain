package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/swap/types"
	tokentypes "github.com/okex/okchain/x/token/types"
)

// Keeper of the swap store
type Keeper struct {
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper

	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a swap keeper
func NewKeeper(supplyKeeper types.SupplyKeeper, tokenKeeper types.TokenKeeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		tokenKeeper:  tokenKeeper,
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

// GetSwapTokenPair gets SwapTokenPair with quote token name
func (k Keeper) GetSwapTokenPair(ctx sdk.Context, tokenPairName string) (types.SwapTokenPair, error) {
	store := ctx.KVStore(k.storeKey)
	var item types.SwapTokenPair
	byteKey := []byte(tokenPairName)
	rawItem := store.Get(byteKey)
	if len(rawItem) == 0 && tokenPairName == types.TestSwapTokenPairName {
		item = types.GetTestSwapTokenPair()
		k.SetSwapTokenPair(ctx, tokenPairName, item)
	}else {
		err := k.cdc.UnmarshalBinaryLengthPrefixed(rawItem, &item)
		if err != nil {
			return types.SwapTokenPair{}, err
		}
	}

	return item, nil
}

// Sets the entire SwapTokenPair data struct for a quote token name
func (k Keeper) SetSwapTokenPair(ctx sdk.Context, tokenPairName string, swapTokenPair types.SwapTokenPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(swapTokenPair)
	store.Set([]byte(tokenPairName), bz)
}

// Deletes the entire SwapTokenPair data struct for a quote token name
func (k Keeper) DeleteSwapTokenPair(ctx sdk.Context, tokenPairName string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(tokenPairName))
}

// Get an iterator over all SwapTokenPairs in which the keys are the names and the values are the whois
func (k Keeper) GetSwapTokenPairsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// NewToken new token
func (k Keeper) NewPoolToken(ctx sdk.Context, token tokentypes.Token) {
	k.tokenKeeper.NewToken(ctx, token)
}

// GetTokenInfo gets the token's info
func (k Keeper) GetPoolTokenInfo(ctx sdk.Context, symbol string) (tokentypes.Token, error) {
	poolToken := k.tokenKeeper.GetTokenInfo(ctx, symbol)
	if poolToken.Owner == nil {
		poolToken = types.InitPoolToken(symbol)
		k.NewPoolToken(ctx, poolToken)
	}
	return poolToken, nil
}

func (k Keeper) GetPoolTokenAmount(ctx sdk.Context, symbol string) (sdk.Dec, error) {
	poolToken, err := k.GetPoolTokenInfo(ctx, symbol)
	return poolToken.TotalSupply, err
}

func (k Keeper) UpdatePoolToken(ctx sdk.Context, token tokentypes.Token) {
	k.tokenKeeper.UpdateToken(ctx, token)
}

func (k Keeper) MintPoolCoinsToUser(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

func (k Keeper) BurnPoolCoinsFromUser(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return k.supplyKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

func (k Keeper) SendCoinsToPool(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
}

func (k Keeper) SendCoinsFromPoolToAccount(ctx sdk.Context, coins sdk.DecCoins, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

func (k Keeper) GetTokenKeeper() types.TokenKeeper {
	return k.tokenKeeper
}

