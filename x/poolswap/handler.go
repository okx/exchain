package poolswap

import (
	"fmt"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/poolswap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the poolswap type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() sdk.Result
		var name string
		switch msg := msg.(type) {
		case types.MsgAddLiquidity:
			name = "handleMsgAddLiquidity"
			handlerFun = func() sdk.Result {
				return handleMsgAddLiquidity(ctx, k, msg)
			}
		case types.MsgRemoveLiquidity:
			name = "handleMsgRemoveLiquidity"
			handlerFun = func() sdk.Result {
				return handleMsgRemoveLiquidity(ctx, k, msg)
			}
		case types.MsgCreateExchange:
			name = "handleMsgCreateExchange"
			handlerFun = func() sdk.Result {
				return handleMsgCreateExchange(ctx, k, msg)
			}
		case types.MsgTokenToNativeToken:
			name = "handleMsgTokenToNativeToken"
			handlerFun = func() sdk.Result {
				return handleMsgTokenToTokenExchange(ctx, k, msg)
			}
		default:
			errMsg := fmt.Sprintf("Invalid msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgTokenToTokenExchange(ctx sdk.Context, k Keeper, msg types.MsgTokenToNativeToken) sdk.Result {
	if msg.SoldTokenAmount.Denom != sdk.DefaultBondDenom && msg.MinBoughtTokenAmount.Denom != sdk.DefaultBondDenom {
		return handleMsgTokenToToken(ctx, k, msg)
	}
	return handleMsgTokenToNativeToken(ctx, k, msg)
}

func handleMsgCreateExchange(ctx sdk.Context, k Keeper, msg types.MsgCreateExchange) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	err := k.IsTokenExist(ctx, msg.Token)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	tokenPair := msg.Token + "_" + common.NativeToken

	swapTokenPair, err := k.GetSwapTokenPair(ctx, tokenPair)
	if err == nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: exchange already exists",
		}
	}

	poolName := types.PoolTokenPrefix + msg.Token
	baseToken := sdk.NewDecCoinFromDec(msg.Token, sdk.ZeroDec())
	quoteToken := sdk.NewDecCoinFromDec(common.NativeToken, sdk.ZeroDec())
	poolToken, err := k.GetPoolTokenInfo(ctx, poolName)
	if err == nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: pool token already exists",
		}
	}
	k.NewPoolToken(ctx, poolName)
	event = event.AppendAttributes(sdk.NewAttribute("pool-token", poolToken.OriginalSymbol))
	swapTokenPair.BasePooledCoin = baseToken
	swapTokenPair.QuotePooledCoin = quoteToken
	swapTokenPair.PoolTokenName = poolName

	k.SetSwapTokenPair(ctx, tokenPair, swapTokenPair)

	event = event.AppendAttributes(sdk.NewAttribute("token-pair", tokenPair))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg types.MsgAddLiquidity) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: block time exceeded deadline",
		}
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPair())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}
	baseTokens := sdk.NewDecCoinFromDec(msg.MaxBaseAmount.Denom, sdk.ZeroDec())
	var liquidity sdk.Dec
	poolToken, err := k.GetPoolTokenInfo(ctx, swapTokenPair.PoolTokenName)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("failed to get pool token %s : %s", swapTokenPair.PoolTokenName, err.Error()),
		}
	}
	if swapTokenPair.QuotePooledCoin.Amount.IsZero() && swapTokenPair.BasePooledCoin.Amount.IsZero() {
		baseTokens.Amount = msg.MaxBaseAmount.Amount
		liquidity = sdk.NewDec(1)
	} else if swapTokenPair.BasePooledCoin.IsPositive() && swapTokenPair.QuotePooledCoin.IsPositive() {
		baseTokens.Amount = mulAndQuo(msg.QuoteAmount.Amount, swapTokenPair.BasePooledCoin.Amount, swapTokenPair.QuotePooledCoin.Amount)
		if poolToken.TotalSupply.IsZero() {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("unexpected totalSupply in pool token %s", poolToken.String()),
			}
		}
		liquidity = mulAndQuo(msg.QuoteAmount.Amount, poolToken.TotalSupply, swapTokenPair.QuotePooledCoin.Amount)
	} else {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("invalid token pair %s", swapTokenPair.String()),
		}
	}
	if baseTokens.Amount.GT(msg.MaxBaseAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "The required base token amount are greater than MaxBaseAmount",
		}
	}
	if liquidity.LT(msg.MinLiquidity) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "The available liquidity is less than MinLiquidity",
		}
	}

	// transfer coins
	coins := sdk.DecCoins{
		msg.QuoteAmount,
		baseTokens,
	}

	coins = coinSort(coins)
	//TODO another coin connot send to pool
	err = k.SendCoinsToPool(ctx, coins, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient coins %s", err.Error()),
		}
	}
	// update swapTokenPair
	swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Add(msg.QuoteAmount)
	swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Add(baseTokens)
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPair(), swapTokenPair)

	// update poolToken
	poolCoins := sdk.NewDecCoinFromDec(poolToken.Symbol, liquidity)
	err = k.MintPoolCoinsToUser(ctx, sdk.DecCoins{poolCoins}, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "failed to mint pool token",
		}
	}

	event.AppendAttributes(sdk.NewAttribute("liquidity", liquidity.String()))
	event.AppendAttributes(sdk.NewAttribute("baseAmount", baseTokens.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgRemoveLiquidity(ctx sdk.Context, k Keeper, msg types.MsgRemoveLiquidity) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: block time exceeded deadline",
		}
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPair())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	liquidity := msg.Liquidity
	poolTokenAmount, err := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("failed to get pool token %s : %s", swapTokenPair.PoolTokenName, err.Error()),
		}
	}
	if poolTokenAmount.LT(liquidity) {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient pool token",
		}
	}

	baseDec := mulAndQuo(swapTokenPair.BasePooledCoin.Amount, liquidity, poolTokenAmount)
	quoteDec := mulAndQuo(swapTokenPair.QuotePooledCoin.Amount, liquidity, poolTokenAmount)
	baseAmount := sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
	quoteAmount := sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)

	if baseAmount.IsLT(msg.MinBaseAmount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: The available baseAmount(%s) are less than MinBaseAmount(%s)", baseAmount.String(), msg.MinBaseAmount.String()),
		}
	}
	if quoteAmount.IsLT(msg.MinQuoteAmount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: available quote amount are less than least quote amount",
		}
	}

	// transfer coins
	coins := sdk.DecCoins{
		baseAmount,
		quoteAmount,
	}
	coins = coinSort(coins)
	err = k.SendCoinsFromPoolToAccount(ctx, coins, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient coins",
		}
	}
	// update swapTokenPair
	swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Sub(quoteAmount)
	swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Sub(baseAmount)
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPair(), swapTokenPair)

	// update poolToken
	poolCoins := sdk.NewDecCoinFromDec(swapTokenPair.PoolTokenName, liquidity)
	err = k.BurnPoolCoinsFromUser(ctx, sdk.DecCoins{poolCoins}, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "failed to burn pool token",
		}
	}

	event.AppendAttributes(sdk.NewAttribute("quoteAmount", quoteAmount.String()))
	event.AppendAttributes(sdk.NewAttribute("baseAmount", baseAmount.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTokenToNativeToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToNativeToken) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if err := common.HasSufficientCoins(msg.Sender, k.GetTokenKeeper().GetCoins(ctx, msg.Sender),
		sdk.DecCoins{msg.SoldTokenAmount}); err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  err.Error(),
		}
	}
	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: block time exceeded deadline",
		}
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPair())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}
	params := k.GetParams(ctx)
	tokenBuy := calculateTokenToBuy(swapTokenPair, msg, params)
	if tokenBuy.Amount.LT(msg.MinBoughtTokenAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: expected minimum token to buy is %s but got %s", msg.MinBoughtTokenAmount, tokenBuy),
		}
	}

	res := swapTokenNativeToken(ctx, k, swapTokenPair, tokenBuy, msg)
	if !res.IsOK() {
		return res
	}
	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTokenToToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToNativeToken) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: block time exceeded deadline",
		}
	}
	if err := common.HasSufficientCoins(msg.Sender, k.GetTokenKeeper().GetCoins(ctx, msg.Sender),
		sdk.DecCoins{msg.SoldTokenAmount}); err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  err.Error(),
		}
	}
	tokenPairOne := msg.SoldTokenAmount.Denom + "_" + sdk.DefaultBondDenom
	swapTokenPairOne, err := k.GetSwapTokenPair(ctx, tokenPairOne)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}
	tokenPairTwo := msg.MinBoughtTokenAmount.Denom + "_" + sdk.DefaultBondDenom
	swapTokenPairTwo, err := k.GetSwapTokenPair(ctx, tokenPairTwo)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	nativeAmount := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("0"))
	params := k.GetParams(ctx)
	msgOne := msg
	msgOne.MinBoughtTokenAmount = nativeAmount
	tokenNative := calculateTokenToBuy(swapTokenPairOne, msgOne, params)

	msgTwo := msg
	msgTwo.SoldTokenAmount = tokenNative
	tokenBuy := calculateTokenToBuy(swapTokenPairOne, msgTwo, params)

	if tokenBuy.Amount.LT(msg.MinBoughtTokenAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: expected minimum token to buy is %s but got %s", msg.MinBoughtTokenAmount, tokenBuy),
		}
	}

	res := swapTokenNativeToken(ctx, k, swapTokenPairOne, tokenNative, msgOne)
	if !res.IsOK() {
		return res
	}
	//TODO if fail,revert last swap
	res = swapTokenNativeToken(ctx, k, swapTokenPairTwo, tokenBuy, msgTwo)
	if !res.IsOK() {
		return res
	}

	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

//calculate the amount to buy
func calculateTokenToBuy(swapTokenPair SwapTokenPair, msg types.MsgTokenToNativeToken, params types.Params) sdk.DecCoin {
	var inputReserve, outputReserve sdk.Dec
	if msg.SoldTokenAmount.Denom == sdk.DefaultBondDenom {
		inputReserve = swapTokenPair.QuotePooledCoin.Amount
		outputReserve = swapTokenPair.BasePooledCoin.Amount
	} else {
		inputReserve = swapTokenPair.BasePooledCoin.Amount
		outputReserve = swapTokenPair.QuotePooledCoin.Amount
	}
	tokenBuyAmt := getInputPrice(msg.SoldTokenAmount.Amount, inputReserve, outputReserve, params.FeeRate)
	tokenBuy := sdk.NewDecCoinFromDec(msg.MinBoughtTokenAmount.Denom, tokenBuyAmt)

	return tokenBuy
}

func swapTokenNativeToken(
	ctx sdk.Context, k Keeper, swapTokenPair SwapTokenPair, tokenBuy sdk.DecCoin,
	msg types.MsgTokenToNativeToken,
) sdk.Result {
	// transfer coins
	err := k.SendCoinsToPool(ctx, sdk.DecCoins{msg.SoldTokenAmount}, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient Coins",
		}
	}

	err = k.SendCoinsFromPoolToAccount(ctx, sdk.DecCoins{tokenBuy}, msg.Recipient)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient Coins",
		}
	}

	// update swapTokenPair
	if msg.SoldTokenAmount.Denom == sdk.DefaultBondDenom {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Add(msg.SoldTokenAmount)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Sub(tokenBuy)
	} else {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Sub(tokenBuy)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Add(msg.SoldTokenAmount)
	}
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPair(), swapTokenPair)
	return sdk.Result{}
}

func getInputPrice(inputAmount, inputReserve, outputReserve, feeRate sdk.Dec) sdk.Dec {
	//if !inputReserve.IsPositive() || !outputReserve.IsPositive() {
	//	panic("should not happen")
	//}
	inputAmountWithFee := inputAmount.Mul(sdk.OneDec().Sub(feeRate).Mul(sdk.NewDec(1000)))
	denominator := inputReserve.Mul(sdk.NewDec(1000)).Add(inputAmountWithFee)
	return mulAndQuo(inputAmountWithFee, outputReserve, denominator)
}

func coinSort(coins sdk.DecCoins) sdk.DecCoins {
	var newCoins sdk.DecCoins
	for _, coin := range coins {
		if coin.Amount.IsPositive() {
			newCoins = append(newCoins, coin)
		}
	}
	newCoins = newCoins.Sort()
	return newCoins
}

var (
	// 10^8
	auxiliaryDec = sdk.NewDec(100000000)
)

// mulAndQuo returns a * b / c
func mulAndQuo(a, b, c sdk.Dec) sdk.Dec {
	a = a.Mul(auxiliaryDec)
	return a.Mul(b).Quo(c).Quo(auxiliaryDec)
}
