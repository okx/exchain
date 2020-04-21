package swap

import (
	"fmt"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the swap type messages
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
		case types.MsgTokenOKTSwap:
			name = "handleMsgTokenOKTSwap"
			handlerFun = func() sdk.Result {
				return handleMsgTokenOKTSwap(ctx, k, msg)
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

func handleMsgCreateExchange(ctx sdk.Context, k Keeper, msg types.MsgCreateExchange) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	err := k.IsTokenExits(ctx, msg.Token)
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
			Log:  "Failed to create Exchange: exchange is exit",
		}
	}

	poolName := "oip3-" + msg.Token
	baseToken := sdk.NewDecCoinFromDec(msg.Token, sdk.ZeroDec())
	quoteToken := sdk.NewDecCoinFromDec(common.NativeToken, sdk.ZeroDec())
	poolToken, err := k.GetPoolTokenInfo(ctx, poolName)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "pool token is not exist",
		}
	}
	if len(poolToken.Symbol) == 0 {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "Failed to create Exchange: Pool Token not exit",
		}
	}

	swapTokenPair.BasePooledCoin = baseToken
	swapTokenPair.QuotePooledCoin = quoteToken
	swapTokenPair.PoolTokenName = poolName

	k.SetSwapTokenPair(ctx, tokenPair, swapTokenPair)

	event.AppendAttributes(sdk.NewAttribute("tokenpair", tokenPair))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddLiquidity(ctx sdk.Context, k Keeper, msg types.MsgAddLiquidity) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	if msg.Deadline < ctx.BlockTime().Unix() {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "blockTime exceeded deadline",
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
			Log:  fmt.Sprintf("get poolToken %s failed: %s", swapTokenPair.PoolTokenName, err.Error()),
		}
	}
	if swapTokenPair.QuotePooledCoin.Amount.IsZero() && swapTokenPair.BasePooledCoin.Amount.IsZero() {
		baseTokens.Amount = msg.MaxBaseAmount.Amount
		liquidity = sdk.NewDec(1)
	} else if swapTokenPair.BasePooledCoin.IsPositive() && swapTokenPair.QuotePooledCoin.IsPositive() {
		baseTokens.Amount = msg.QuoteAmount.Amount.Mul(swapTokenPair.BasePooledCoin.Amount).Quo(swapTokenPair.QuotePooledCoin.Amount)
		if poolToken.TotalSupply.IsZero() {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("unexpected totalSupply in poolToken %s", poolToken.String()),
			}
		}
		liquidity = msg.QuoteAmount.Amount.Quo(swapTokenPair.QuotePooledCoin.Amount).Mul(poolToken.TotalSupply)
	} else {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("invalid swapTokenPair %s", swapTokenPair.String()),
		}
	}
	if baseTokens.Amount.GT(msg.MaxBaseAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("MaxBaseAmount is too high"),
		}
	}
	if liquidity.LT(msg.MinLiquidity) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("MinLiquidity is too low"),
		}
	}

	// transfer coins
	coins := sdk.DecCoins{
		msg.QuoteAmount,
		baseTokens,
	}
	err = k.SendCoinsToPool(ctx, coins, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient Coins"),
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
			Log:  fmt.Sprintf("fail to mint poolCoins"),
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
			Log:  "blockTime exceeded deadline",
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
			Log:  fmt.Sprintf("get poolToken %s failed: %s", swapTokenPair.PoolTokenName, err.Error()),
		}
	}
	if poolTokenAmount.LT(liquidity) {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  "insufficient poolToken",
		}
	}

	baseDec := swapTokenPair.BasePooledCoin.Amount.Mul(liquidity).Quo(poolTokenAmount)
	quoteDec := swapTokenPair.QuotePooledCoin.Amount.Mul(liquidity).Quo(poolTokenAmount)
	baseAmount := sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
	quoteAmount := sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)

	if baseAmount.IsLT(msg.MinBaseAmount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "MinBaseAmount is too high",
		}
	}
	if quoteAmount.IsLT(msg.MinQuoteAmount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  "MinQuoteAmount is too high",
		}
	}

	// transfer coins
	coins := sdk.DecCoins{
		baseAmount,
		quoteAmount,
	}
	err = k.SendCoinsFromPoolToAccount(ctx, coins, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient Coins"),
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
			Log:  fmt.Sprintf("fail to burn poolCoins"),
		}
	}

	event.AppendAttributes(sdk.NewAttribute("quoteAmount", quoteAmount.String()))
	event.AppendAttributes(sdk.NewAttribute("baseAmount", baseAmount.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTokenOKTSwap(ctx sdk.Context, k Keeper, msg types.MsgTokenOKTSwap) sdk.Result {
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
			Log:  "blockTime exceeded deadline",
		}
	}
	swapTokenPair, err := k.GetSwapTokenPair(ctx, msg.GetSwapTokenPair())
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  err.Error(),
		}
	}
	res, tokenBuy := calculateTokenToBuy(swapTokenPair, msg)
	if !res.IsOK() {
		return res
	}
	res = swapTokenOKT(ctx, k, swapTokenPair, tokenBuy, msg)
	if !res.IsOK() {
		return res
	}
	event.AppendAttributes(sdk.NewAttribute("bought_token_amount", tokenBuy.String()))
	event.AppendAttributes(sdk.NewAttribute("recipient", msg.Recipient.String()))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

//calculate the amount to buy
func calculateTokenToBuy(swapTokenPair SwapTokenPair, msg types.MsgTokenOKTSwap) (sdk.Result, sdk.DecCoin) {
	var inputReserve, outputReserve sdk.Dec
	if msg.SoldTokenAmount.Denom == sdk.DefaultBondDenom {
		inputReserve = swapTokenPair.QuotePooledCoin.Amount
		outputReserve = swapTokenPair.BasePooledCoin.Amount
	} else {
		inputReserve = swapTokenPair.BasePooledCoin.Amount
		outputReserve = swapTokenPair.QuotePooledCoin.Amount
	}
	tokenBuyAmt := getInputPrice(msg.SoldTokenAmount.Amount, inputReserve, outputReserve)
	tokenBuy := sdk.NewDecCoinFromDec(msg.MinBoughtTokenAmount.Denom, tokenBuyAmt)
	if tokenBuyAmt.LT(msg.MinBoughtTokenAmount.Amount) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("expected minimum token to buy is %s but got %s", msg.MinBoughtTokenAmount, tokenBuy),
		}, sdk.DecCoin{}
	}
	return sdk.Result{}, tokenBuy
}

func swapTokenOKT(
	ctx sdk.Context, k Keeper, swapTokenPair SwapTokenPair, tokenBuy sdk.DecCoin,
	msg types.MsgTokenOKTSwap,
) sdk.Result {
	// transfer coins
	err := k.SendCoinsToPool(ctx, sdk.DecCoins{msg.SoldTokenAmount}, msg.Sender)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient Coins"),
		}
	}

	err = k.SendCoinsFromPoolToAccount(ctx, sdk.DecCoins{tokenBuy}, msg.Recipient)
	if err != nil {
		return sdk.Result{
			Code: sdk.CodeInsufficientCoins,
			Log:  fmt.Sprintf("insufficient Coins"),
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

func getInputPrice(inputAmount, inputReserve, outputReserve sdk.Dec) sdk.Dec {
	if !inputReserve.IsPositive() || !outputReserve.IsPositive() {
		panic("should not happen")
	}
	inputAmountWithFee := inputAmount.Mul(sdk.OneDec().Sub(types.FeeRate))
	numerator := inputAmountWithFee.Mul(outputReserve)
	denominator := inputReserve.Add(inputAmountWithFee)
	return numerator.Quo(denominator)
}
