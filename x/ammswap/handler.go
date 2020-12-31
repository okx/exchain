package ammswap

import (
	"github.com/okex/okexchain/x/ammswap/keeper"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/common/perf"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the ammswap type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() (*sdk.Result, error)
		var name string
		switch msg := msg.(type) {
		case types.MsgAddLiquidity:
			name = "handleMsgAddLiquidity"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgAddLiquidity(ctx, k, msg)
			}
		case types.MsgRemoveLiquidity:
			name = "handleMsgRemoveLiquidity"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgRemoveLiquidity(ctx, k, msg)
			}
		case types.MsgCreateExchange:
			name = "handleMsgCreateExchange"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgCreateExchange(ctx, k, msg)
			}
		case types.MsgTokenToToken:
			name = "handleMsgTokenToToken"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenToToken(ctx, k, msg)
			}
		default:
			return nil, types.ErrSwapUnknownMsgType()
		}
		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)

		res, err := handlerFun()
		common.SanityCheckHandler(res, err)
		return res, err
	}
}

func handleMsgTokenToToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToToken) (*sdk.Result, error) {
	_, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return swapTokenByRouter(ctx, k, msg)
	} else {
		return swapToken(ctx, k, msg)
	}
}

func handleMsgCreateExchange(ctx sdk.Context, k Keeper, msg types.MsgCreateExchange) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	// 0. check if 2 tokens exist
	err := k.IsTokenExist(ctx, msg.Token0Name)
	if err != nil {
		return nil, err
	}

	err = k.IsTokenExist(ctx, msg.Token1Name)
	if err != nil {
		return nil, err
	}

	// 1. check if the token pair exists
	tokenPairName := msg.GetSwapTokenPairName()
	_, err = k.GetSwapTokenPair(ctx, tokenPairName)
	if err == nil {
		return types.ErrSwapTokenPairExist().Result()
	}

	// 2. check if the pool token exists
	poolTokenName := types.GetPoolTokenName(msg.Token0Name, msg.Token1Name)
	_, err = k.GetPoolTokenInfo(ctx, poolTokenName)
	if err == nil {
		return types.ErrPoolTokenPairExist().Result()
	}

	// 3. create the pool token
	k.NewPoolToken(ctx, poolTokenName)

	// 4. create the token pair
	swapTokenPair := types.NewSwapPair(msg.Token0Name, msg.Token1Name)
	k.SetSwapTokenPair(ctx, tokenPairName, swapTokenPair)

	// 5. notify backend module
	k.OnCreateExchange(ctx, swapTokenPair)

	event = event.AppendAttributes(sdk.NewAttribute("pool-token-name", poolTokenName))
	event = event.AppendAttributes(sdk.NewAttribute("token-pair", tokenPairName))
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg types.MsgAddLiquidity) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	if msg.Deadline < ctx.BlockTime().Unix() {
		return types.ErrBlockTimeBigThanDeadline().Result()
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return nil, err
	}
	baseTokens := sdk.NewDecCoinFromDec(msg.MaxBaseAmount.Denom, sdk.ZeroDec())
	var liquidity sdk.Dec
	poolToken, err := k.GetPoolTokenInfo(ctx, swapTokenPair.PoolTokenName)
	if err != nil {
		return nil, err
	}
	if swapTokenPair.QuotePooledCoin.Amount.IsZero() && swapTokenPair.BasePooledCoin.Amount.IsZero() {
		baseTokens.Amount = msg.MaxBaseAmount.Amount
		liquidity = sdk.NewDec(1)
	} else if swapTokenPair.BasePooledCoin.IsPositive() && swapTokenPair.QuotePooledCoin.IsPositive() {
		baseTokens.Amount = common.MulAndQuo(msg.QuoteAmount.Amount, swapTokenPair.BasePooledCoin.Amount, swapTokenPair.QuotePooledCoin.Amount)
		totalSupply := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
		if baseTokens.IsZero() {
			baseTokens.Amount = sdk.NewDecWithPrec(1, sdk.Precision)
		}
		if totalSupply.IsZero() {
			return types.ErrIsZeroValue("totalSupply").Result()
		}
		liquidity = common.MulAndQuo(msg.QuoteAmount.Amount, totalSupply, swapTokenPair.QuotePooledCoin.Amount)
		if liquidity.IsZero() {
			return types.ErrIsZeroValue("liquidity").Result()
		}
	} else {
		return types.ErrInvalidTokenPair(swapTokenPair.String()).Result()
	}
	if baseTokens.Amount.GT(msg.MaxBaseAmount.Amount) {
		return types.ErrBaseTokensAmountBiggerThanMax().Result()
	}
	if liquidity.LT(msg.MinLiquidity) {
		return types.ErrLessThan("liquidity", "min liquidity").Result()
	}

	// transfer coins
	coins := sdk.SysCoins{
		msg.QuoteAmount,
		baseTokens,
	}

	coins = coinSort(coins)

	err = k.SendCoinsToPool(ctx, coins, msg.Sender)
	if err != nil {
		return types.ErrSendCoinsFailed(err).Result()
	}
	// update swapTokenPair
	swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Add(msg.QuoteAmount)
	swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Add(baseTokens)
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)

	// update poolToken
	poolCoins := sdk.NewDecCoinFromDec(poolToken.Symbol, liquidity)
	err = k.MintPoolCoinsToUser(ctx, sdk.SysCoins{poolCoins}, msg.Sender)
	if err != nil {
		return types.ErrMintPoolTokenFailed(err).Result()
	}

	event.AppendAttributes(sdk.NewAttribute("liquidity", liquidity.String()))
	event.AppendAttributes(sdk.NewAttribute("baseAmount", baseTokens.String()))
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveLiquidity(ctx sdk.Context, k Keeper, msg types.MsgRemoveLiquidity) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return types.ErrMsgDeadlineLessThanBlockTime().Result()
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return nil, err
	}

	liquidity := msg.Liquidity
	poolTokenAmount := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if poolTokenAmount.LT(liquidity) {
		return types.ErrLessThan("pool token amount", "liquidity").Result()
	}

	baseDec := common.MulAndQuo(swapTokenPair.BasePooledCoin.Amount, liquidity, poolTokenAmount)
	quoteDec := common.MulAndQuo(swapTokenPair.QuotePooledCoin.Amount, liquidity, poolTokenAmount)

	baseAmount := sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
	quoteAmount := sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)

	if baseAmount.IsLT(msg.MinBaseAmount) {
		return types.ErrLessThan("base amount", "min base amount").Result()
	}
	if quoteAmount.IsLT(msg.MinQuoteAmount) {
		return types.ErrLessThan("quote amount", "min quote amount").Result()
	}

	// transfer coins
	coins := sdk.SysCoins{
		baseAmount,
		quoteAmount,
	}
	coins = coinSort(coins)
	err = k.SendCoinsFromPoolToAccount(ctx, coins, msg.Sender)
	if err != nil {
		return types.ErrSendCoinsFromPoolToAccountFailed(err.Error()).Result()
	}
	// update swapTokenPair
	swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Sub(quoteAmount)
	swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Sub(baseAmount)
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)

	// update poolToken
	poolCoins := sdk.NewDecCoinFromDec(swapTokenPair.PoolTokenName, liquidity)
	err = k.BurnPoolCoinsFromUser(ctx, sdk.SysCoins{poolCoins}, msg.Sender)
	if err != nil {
		return types.ErrBurnPoolTokenFailed(err).Result()
	}

	event.AppendAttributes(sdk.NewAttribute("quoteAmount", quoteAmount.String()))
	event.AppendAttributes(sdk.NewAttribute("baseAmount", baseAmount.String()))
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func swapToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToToken) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if err := common.HasSufficientCoins(msg.Sender, k.GetTokenKeeper().GetCoins(ctx, msg.Sender),
		sdk.SysCoins{msg.SoldTokenAmount}); err != nil {
		return common.ErrInsufficientCoins(DefaultParamspace, err.Error()).Result()
	}
	if msg.Deadline < ctx.BlockTime().Unix() {
		return types.ErrBlockTimeBigThanDeadline().Result()
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return nil, err
	}
	if swapTokenPair.BasePooledCoin.IsZero() || swapTokenPair.QuotePooledCoin.IsZero() {
		return types.ErrIsZeroValue("base pooled coin or quote pooled coin").Result()
	}
	params := k.GetParams(ctx)
	tokenBuy := keeper.CalculateTokenToBuy(swapTokenPair, msg.SoldTokenAmount, msg.MinBoughtTokenAmount.Denom, params)
	if tokenBuy.IsZero() {
		return types.ErrIsZeroValue("token buy").Result()
	}
	if tokenBuy.Amount.LT(msg.MinBoughtTokenAmount.Amount) {
		return types.ErrLessThan("token buy amount", "min bought token amount").Result()
	}

	res, err := swapTokenNativeToken(ctx, k, swapTokenPair, tokenBuy, msg)
	if err != nil {
		return res, err
	}
	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func swapTokenByRouter(ctx sdk.Context, k Keeper, msg types.MsgTokenToToken) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return types.ErrBlockTimeBigThanDeadline().Result()
	}
	if err := common.HasSufficientCoins(msg.Sender, k.GetTokenKeeper().GetCoins(ctx, msg.Sender),
		sdk.SysCoins{msg.SoldTokenAmount}); err != nil {
		return common.ErrInsufficientCoins(DefaultParamspace, err.Error()).Result()
	}
	tokenPairOne := types.GetSwapTokenPairName(msg.SoldTokenAmount.Denom, sdk.DefaultBondDenom)
	swapTokenPairOne, err := k.GetSwapTokenPair(ctx, tokenPairOne)
	if err != nil {
		return nil, err
	}
	if swapTokenPairOne.BasePooledCoin.IsZero() || swapTokenPairOne.QuotePooledCoin.IsZero() {
		return types.ErrIsZeroValue("base pooled coin or quote pooled coin").Result()
	}
	tokenPairTwo := types.GetSwapTokenPairName(msg.MinBoughtTokenAmount.Denom, sdk.DefaultBondDenom)
	swapTokenPairTwo, err := k.GetSwapTokenPair(ctx, tokenPairTwo)
	if err != nil {
		return nil, err
	}
	if swapTokenPairTwo.BasePooledCoin.IsZero() || swapTokenPairTwo.QuotePooledCoin.IsZero() {
		return types.ErrIsZeroValue("base pooled coin or quote pooled coin").Result()
	}

	nativeAmount := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("0"))
	params := k.GetParams(ctx)
	msgOne := msg
	msgOne.MinBoughtTokenAmount = nativeAmount
	tokenNative := keeper.CalculateTokenToBuy(swapTokenPairOne, msgOne.SoldTokenAmount, msgOne.MinBoughtTokenAmount.Denom, params)
	if tokenNative.IsZero() {
		return types.ErrIsZeroValue("token native").Result()
	}
	msgTwo := msg
	msgTwo.SoldTokenAmount = tokenNative
	tokenBuy := keeper.CalculateTokenToBuy(swapTokenPairTwo, msgTwo.SoldTokenAmount, msgTwo.MinBoughtTokenAmount.Denom, params)
	// sanity check. user may set MinBoughtTokenAmount to zero on front end.
	// if set zero,this will not return err
	if tokenBuy.IsZero() {
		return types.ErrIsZeroValue("token to buy amount").Result()
	}
	if tokenBuy.Amount.LT(msg.MinBoughtTokenAmount.Amount) {
		return types.ErrLessThan("token buy amount", "min bought token amount").Result()
	}

	res, err := swapTokenNativeToken(ctx, k, swapTokenPairOne, tokenNative, msgOne)
	if err != nil {
		return res, err
	}
	res, err = swapTokenNativeToken(ctx, k, swapTokenPairTwo, tokenBuy, msgTwo)
	if err != nil {
		return res, err
	}

	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func swapTokenNativeToken(
	ctx sdk.Context, k Keeper, swapTokenPair SwapTokenPair, tokenBuy sdk.SysCoin,
	msg types.MsgTokenToToken,
) (*sdk.Result, error) {
	// transfer coins
	err := k.SendCoinsToPool(ctx, sdk.SysCoins{msg.SoldTokenAmount}, msg.Sender)
	if err != nil {
		return types.ErrSendCoinsToPoolFailed(err.Error()).Result()
	}

	err = k.SendCoinsFromPoolToAccount(ctx, sdk.SysCoins{tokenBuy}, msg.Recipient)
	if err != nil {
		return types.ErrSendCoinsFromPoolToAccountFailed(err.Error()).Result()
	}

	// update swapTokenPair
	if msg.MinBoughtTokenAmount.Denom < msg.SoldTokenAmount.Denom {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Add(msg.SoldTokenAmount)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Sub(tokenBuy)
	} else {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Sub(tokenBuy)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Add(msg.SoldTokenAmount)
	}
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)
	k.OnSwapToken(ctx, msg.Recipient, swapTokenPair, msg.SoldTokenAmount, tokenBuy)
	return &sdk.Result{}, nil
}

func coinSort(coins sdk.SysCoins) sdk.SysCoins {
	var newCoins sdk.SysCoins
	for _, coin := range coins {
		if coin.Amount.IsPositive() {
			newCoins = append(newCoins, coin)
		}
	}
	newCoins = newCoins.Sort()
	return newCoins
}
