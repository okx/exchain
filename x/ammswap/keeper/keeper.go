package keeper

import (
	"fmt"

	"github.com/okex/okexchain/x/common"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap/types"
	tokentypes "github.com/okex/okexchain/x/token/types"
)

// Keeper of the swap store
type Keeper struct {
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper

	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	paramSpace     types.ParamSubspace
	ObserverKeeper []types.BackendKeeper
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
	if rawItem == nil {
		return types.SwapTokenPair{}, types.ErrUnexistswapTokenPair(types.DefaultCodespace, fmt.Sprintf("non-existent swapTokenPair: %s", tokenPairName))
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(rawItem, &item)
	if err != nil {
		return types.SwapTokenPair{}, common.ErrUnMarshalJSONFailed(err.Error())
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

func (k Keeper) GetSwapTokenPairs(ctx sdk.Context) []types.SwapTokenPair {
	var result []types.SwapTokenPair
	iterator := k.GetSwapTokenPairsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		tokenPair := types.SwapTokenPair{}
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &tokenPair)
		result = append(result, tokenPair)
	}
	return result
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
		return poolToken, types.ErrUnexistPoolToken(types.DefaultCodespace, fmt.Sprintf("Pool token %s does not exist", symbol))
	}
	return poolToken, nil
}

// GetPoolTokenAmount gets the amount of the specified poolToken name
func (k Keeper) GetPoolTokenAmount(ctx sdk.Context, poolTokenName string) sdk.Dec {
	return k.supplyKeeper.GetSupplyByDenom(ctx, poolTokenName)
}

// MintPoolCoinsToUser mints coins and send them to the specified user address
func (k Keeper) MintPoolCoinsToUser(ctx sdk.Context, coins sdk.SysCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return types.ErrCodeMinCoinsFailed(types.DefaultCodespace, err.Error())
	}
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

// BurnPoolCoinsFromUser sends coins to account module and burns them
func (k Keeper) BurnPoolCoinsFromUser(ctx sdk.Context, coins sdk.SysCoins, addr sdk.AccAddress) error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return types.ErrSendCoinsFromAccountToModule(types.DefaultCodespace, err.Error())
	}
	return k.supplyKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

// SendCoinsToPool sends coins from user account to module account
func (k Keeper) SendCoinsToPool(ctx sdk.Context, coins sdk.SysCoins, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
}

// SendCoinsFromPoolToAccount sends coins from module account to user account
func (k Keeper) SendCoinsFromPoolToAccount(ctx sdk.Context, coins sdk.SysCoins, addr sdk.AccAddress) error {
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

func (k Keeper) GetRedeemableAssets(ctx sdk.Context, baseAmountName, quoteAmountName string, liquidity sdk.Dec) (baseAmount, quoteAmount sdk.SysCoin, err error) {
	err = types.ValidateBaseAndQuoteAmount(baseAmountName, quoteAmountName)
	if err != nil {
		return baseAmount, quoteAmount, err
	}
	swapTokenPairName := types.GetSwapTokenPairName(baseAmountName, quoteAmountName)
	swapTokenPair, err := k.GetSwapTokenPair(ctx, swapTokenPairName)
	if err != nil {
		return baseAmount, quoteAmount, err
	}
	poolTokenAmount := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if poolTokenAmount.LT(liquidity) {
		return baseAmount, quoteAmount, types.ErrInsufficientPoolToken(types.DefaultCodespace, "insufficient pool token")
	}

	baseDec := common.MulAndQuo(swapTokenPair.BasePooledCoin.Amount, liquidity, poolTokenAmount)
	quoteDec := common.MulAndQuo(swapTokenPair.QuotePooledCoin.Amount, liquidity, poolTokenAmount)
	baseAmount = sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
	quoteAmount = sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)
	return baseAmount, quoteAmount, nil
}

//CalculateTokenToBuy calculates the amount to buy
func CalculateTokenToBuy(swapTokenPair types.SwapTokenPair, sellToken sdk.SysCoin, buyTokenDenom string, params types.Params) sdk.SysCoin {
	var inputReserve, outputReserve sdk.Dec
	if buyTokenDenom < sellToken.Denom {
		inputReserve = swapTokenPair.QuotePooledCoin.Amount
		outputReserve = swapTokenPair.BasePooledCoin.Amount
	} else {
		inputReserve = swapTokenPair.BasePooledCoin.Amount
		outputReserve = swapTokenPair.QuotePooledCoin.Amount
	}
	tokenBuyAmt := GetInputPrice(sellToken.Amount, inputReserve, outputReserve, params.FeeRate)
	tokenBuy := sdk.NewDecCoinFromDec(buyTokenDenom, tokenBuyAmt)

	return tokenBuy
}

func GetInputPrice(inputAmount, inputReserve, outputReserve, feeRate sdk.Dec) sdk.Dec {
	inputAmountWithFee := inputAmount.MulTruncate(sdk.OneDec().Sub(feeRate).MulTruncate(sdk.NewDec(1000)))
	denominator := inputReserve.MulTruncate(sdk.NewDec(1000)).Add(inputAmountWithFee)
	return common.MulAndQuo(inputAmountWithFee, outputReserve, denominator)
}

func (k *Keeper) SetObserverKeeper(bk types.BackendKeeper) {
	k.ObserverKeeper = append(k.ObserverKeeper, bk)
}

func (k Keeper) OnSwapToken(ctx sdk.Context, address sdk.AccAddress, swapTokenPair types.SwapTokenPair, sellAmount sdk.SysCoin, buyAmount sdk.SysCoin) {
	for _, observer := range k.ObserverKeeper {
		observer.OnSwapToken(ctx, address, swapTokenPair, sellAmount, buyAmount)
	}
}

func (k Keeper) OnCreateExchange(ctx sdk.Context, swapTokenPair types.SwapTokenPair) {
	for _, observer := range k.ObserverKeeper {
		observer.OnSwapCreateExchange(ctx, swapTokenPair)
	}
}
