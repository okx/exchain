package ammswap

import (
	"fmt"
	"github.com/okex/okexchain/x/ammswap/keeper"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/common/perf"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the ammswap type messages
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
		case types.MsgTokenToToken:
			name = "handleMsgTokenToToken"
			handlerFun = func() sdk.Result {
				return handleMsgTokenToToken(ctx, k, msg)
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

func handleMsgTokenToToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToToken) sdk.Result {
	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: block time exceeded deadline",
		}
	}
	var routeList []string
	routeList = append(routeList, msg.SoldTokenAmount.Denom)
	routeList = append(routeList, msg.TokenRoute...)
	routeList = append(routeList, msg.MinBoughtTokenAmount.Denom)
	routeLength := len(routeList)

	for i := 1; i < routeLength; i++ {
		tokenPair := types.GetSwapTokenPairName(routeList[i], routeList[i-1])
		_, err := k.GetSwapTokenPair(ctx, tokenPair)
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeUnknownRequest,
				Log:  err.Error(),
			}
		}
	}

	soldTokenAmount := msg.SoldTokenAmount
	var minBoughtTokenAmount sdk.DecCoin
	var recipient sdk.AccAddress
	for i := 1; i < routeLength; i++ {
		if i < routeLength - 1 {
			minBoughtTokenAmount = sdk.NewDecCoinFromDec(routeList[i], sdk.ZeroDec())
			recipient = msg.Sender
		}else {
			minBoughtTokenAmount = msg.MinBoughtTokenAmount
			recipient = msg.Recipient
		}
		mediumMsg := types.NewMsgTokenToToken(soldTokenAmount, minBoughtTokenAmount, nil, msg.Deadline, recipient, msg.Sender)
		var result sdk.Result
		result, soldTokenAmount = swapToken(ctx, k, mediumMsg)
		if !result.IsOK() {
			return result
		}
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgCreateExchange(ctx sdk.Context, k Keeper, msg types.MsgCreateExchange) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	err := k.IsTokenExist(ctx, msg.NameTokenA)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	err = k.IsTokenExist(ctx, msg.NameTokenB)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	tokenPair := msg.GetSwapTokenPairName()

	swapTokenPair, err := k.GetSwapTokenPair(ctx, tokenPair)
	if err == nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: exchange already exists",
		}
	}

	poolName := types.GetPoolTokenName(msg.NameTokenA, msg.NameTokenB)
	nameBaseToken, nameQuoteToken := types.GetBaseQuoteToken(msg.NameTokenA, msg.NameTokenB)
	baseToken := sdk.NewDecCoinFromDec(nameBaseToken, sdk.ZeroDec())
	quoteToken := sdk.NewDecCoinFromDec(nameQuoteToken, sdk.ZeroDec())
	_, err = k.GetPoolTokenInfo(ctx, poolName)
	if err == nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed: pool token already exists",
		}
	}
	k.NewPoolToken(ctx, poolName)
	event = event.AppendAttributes(sdk.NewAttribute("pool-token", poolName))
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
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		createExchangeMsg := types.NewMsgCreateExchange(msg.MaxAmountTokenA.Denom, msg.AmountTokenB.Denom, msg.Sender)
		createExchangeResult := handleMsgCreateExchange(ctx, k, createExchangeMsg)
		if !createExchangeResult.IsOK() {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log: fmt.Sprintf("create exchange failed: %s", createExchangeResult.Log),
			}
		}
		swapTokenPair, err = k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log: "unexpected logic: failed to create token pair but returned success",
			}
		}
	}

	availableAmountTokenA := sdk.NewDecCoinFromDec(msg.MaxAmountTokenA.Denom, sdk.ZeroDec())
	amountTokenB := msg.AmountTokenB
	var poolAmountTokenA, poolAmountTokenB *sdk.DecCoin
	if availableAmountTokenA.Denom == swapTokenPair.BasePooledCoin.Denom {
		poolAmountTokenA = &swapTokenPair.BasePooledCoin
		poolAmountTokenB = &swapTokenPair.QuotePooledCoin
	}else if availableAmountTokenA.Denom == swapTokenPair.QuotePooledCoin.Denom {
		poolAmountTokenA = &swapTokenPair.QuotePooledCoin
		poolAmountTokenB = &swapTokenPair.BasePooledCoin
	}else {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log: "unexpected logic: invalid swap token pair",
		}
	}
	var availableLiquidity sdk.Dec
	poolToken, err := k.GetPoolTokenInfo(ctx, swapTokenPair.PoolTokenName)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("failed to get pool token %s : %s", swapTokenPair.PoolTokenName, err.Error()),
		}
	}
	if poolAmountTokenA.IsZero() && poolAmountTokenB.IsZero() {
		availableAmountTokenA.Amount = msg.MaxAmountTokenA.Amount
		availableLiquidity = sdk.NewDec(1)
	} else if poolAmountTokenA.IsPositive() && poolAmountTokenB.IsPositive() {
		availableAmountTokenA.Amount = common.MulAndQuo(amountTokenB.Amount, poolAmountTokenA.Amount, poolAmountTokenB.Amount)
		totalSupply := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
		if totalSupply.IsZero() {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("unexpected totalSupply in pool token %s", poolToken.String()),
			}
		}
		availableLiquidity = common.MulAndQuo(amountTokenB.Amount, totalSupply, poolAmountTokenB.Amount)
		if availableLiquidity.IsZero() {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("failed to add liquidity"),
			}
		}
	} else {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("invalid token pair %s", swapTokenPair.String()),
		}
	}
	if availableAmountTokenA.Amount.GT(msg.MaxAmountTokenA.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "The required base token amount are greater than MaxAmountTokenA",
		}
	}
	if availableLiquidity.LT(msg.MinLiquidity) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "The available liquidity is less than MinLiquidity",
		}
	}

	// transfer coins
	coins := sdk.DecCoins{
		amountTokenB,
		availableAmountTokenA,
	}

	coins = coinSort(coins)

	err = k.SendCoinsToPool(ctx, coins, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient coins %s", err.Error()),
		}
	}
	// update swapTokenPair
	*poolAmountTokenA = poolAmountTokenA.Add(availableAmountTokenA)
	*poolAmountTokenB = poolAmountTokenB.Add(amountTokenB)
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)

	// update poolToken
	poolCoins := sdk.NewDecCoinFromDec(poolToken.Symbol, availableLiquidity)
	err = k.MintPoolCoinsToUser(ctx, sdk.DecCoins{poolCoins}, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "failed to mint pool token",
		}
	}

	event = event.AppendAttributes(sdk.NewAttribute("liquidity", availableLiquidity.String()))
	event = event.AppendAttributes(sdk.NewAttribute("availableAmountTokenA", availableAmountTokenA.String()))
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
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}

	liquidity := msg.Liquidity
	poolTokenAmount := k.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if poolTokenAmount.LT(liquidity) {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient pool token",
		}
	}

	baseDec := common.MulAndQuo(swapTokenPair.BasePooledCoin.Amount, liquidity, poolTokenAmount)
	quoteDec := common.MulAndQuo(swapTokenPair.QuotePooledCoin.Amount, liquidity, poolTokenAmount)

	baseAmount := sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
	quoteAmount := sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)

	var minBaseAmount, minQuoteAmount sdk.DecCoin
	if msg.MinAmountTokenA.Denom == baseAmount.Denom {
		minBaseAmount = msg.MinAmountTokenA
		minQuoteAmount = msg.MinAmountTokenB
	}else if msg.MinAmountTokenB.Denom == baseAmount.Denom {
		minBaseAmount = msg.MinAmountTokenB
		minQuoteAmount = msg.MinAmountTokenA
	}else {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log: "unexpected logic: invalid swap token pair",
		}
	}
	if baseAmount.IsLT(minBaseAmount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: The available base Amount(%s) are less than min base Amount(%s)", baseAmount.String(), msg.MinAmountTokenA.String()),
		}
	}
	if quoteAmount.IsLT(minQuoteAmount) {
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
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)

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

func swapToken(ctx sdk.Context, k Keeper, msg types.MsgTokenToToken) (sdk.Result, sdk.DecCoin) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	zeroDecCoin := sdk.NewDecCoinFromDec(msg.MinBoughtTokenAmount.Denom, sdk.ZeroDec())
	if err := common.HasSufficientCoins(msg.Sender, k.GetTokenKeeper().GetCoins(ctx, msg.Sender),
		sdk.DecCoins{msg.SoldTokenAmount}); err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  err.Error(),
		}, zeroDecCoin
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPairName())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  err.Error(),
		}, zeroDecCoin
	}
	params := k.GetParams(ctx)
	tokenBuy := keeper.CalculateTokenToBuy(swapTokenPair, msg.SoldTokenAmount, msg.MinBoughtTokenAmount.Denom, params)
	if tokenBuy.IsZero() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: selled token amount is too little to buy any token"),
		}, zeroDecCoin
	}
	if tokenBuy.Amount.LT(msg.MinBoughtTokenAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("Failed: expected minimum token to buy is %s but got %s", msg.MinBoughtTokenAmount, tokenBuy),
		}, zeroDecCoin
	}

	res := swapBoughtTokenSoldToken(ctx, k, swapTokenPair, tokenBuy, msg)
	if !res.IsOK() {
		return res, zeroDecCoin
	}
	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}, tokenBuy
}

func swapBoughtTokenSoldToken(
	ctx sdk.Context, k Keeper, swapTokenPair SwapTokenPair, tokenBuy sdk.DecCoin,
	msg types.MsgTokenToToken,
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
	if msg.MinBoughtTokenAmount.Denom < msg.SoldTokenAmount.Denom {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Add(msg.SoldTokenAmount)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Sub(tokenBuy)
	} else {
		swapTokenPair.QuotePooledCoin = swapTokenPair.QuotePooledCoin.Sub(tokenBuy)
		swapTokenPair.BasePooledCoin = swapTokenPair.BasePooledCoin.Add(msg.SoldTokenAmount)
	}
	k.SetSwapTokenPair(ctx, msg.GetSwapTokenPairName(), swapTokenPair)
	return sdk.Result{}
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

