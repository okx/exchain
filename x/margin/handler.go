package margin

import (
	"fmt"

	"github.com/okex/okchain/x/margin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// NewHandler creates an sdk.Handler for all the margin type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", ModuleName)

		var handlerFun func() sdk.Result
		switch msg := msg.(type) {
		case types.MsgDexDeposit:
			handlerFun = func() sdk.Result {
				return handleMsgDexDeposit(ctx, k, msg, logger)
			}
		case types.MsgDexWithdraw:
			handlerFun = func() sdk.Result {
				return handleMsgDexWithdraw(ctx, k, msg, logger)
			}
		case types.MsgDexSet:
			handlerFun = func() sdk.Result {
				return handleMsgDexSet(ctx, k, msg, logger)
			}
		case types.MsgDexSave:
			handlerFun = func() sdk.Result {
				return handleMsgDexSave(ctx, k, msg, logger)
			}
		case types.MsgDexReturn:
			handlerFun = func() sdk.Result {
				return handleMsgDexReturn(ctx, k, msg, logger)
			}
		case types.MsgDeposit:
			handlerFun = func() sdk.Result {
				return handleMsgDeposit(ctx, k, msg, logger)
			}
		case types.MsgWithdraw:
			handlerFun = func() sdk.Result {
				return handleMsgWithdraw(ctx, k, msg, logger)
			}

		case types.MsgBorrow:
			handlerFun = func() sdk.Result {
				return handleMsgBorrow(ctx, k, msg, logger)
			}
		case types.MsgRefund:
			handlerFun = func() sdk.Result {
				return handleMsgRefund(ctx, k, msg, logger)
			}

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
		return handlerFun()
	}
}

func handleMsgDexDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDexDeposit, logger log.Logger) sdk.Result {
	tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, msg.Product)
	if tokenPair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because non-exist product: %s", msg.Product)).Result()
	}
	if !tokenPair.Owner.Equals(msg.Address) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to deposit because %s is not the owner of product:%s", msg.Address.String(), msg.Product)).Result()
	}

	if sdkErr := keeper.DexDeposit(ctx, msg.Address, msg.Product, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}
	logger.Debug(fmt.Sprintf("successfully handleMsgDexDeposit: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDexWithdraw(ctx sdk.Context, keeper Keeper, msg types.MsgDexWithdraw, logger log.Logger) sdk.Result {
	if sdkErr := keeper.DexWithdraw(ctx, msg.Product, msg.Address, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgDexWithdraw: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDexSet(ctx sdk.Context, keeper Keeper, msg types.MsgDexSet, logger log.Logger) sdk.Result {
	if sdkErr := keeper.DexSet(ctx, msg.Address, msg.Product, msg.MaxLeverage, msg.BorrowRate, msg.MaintenanceMarginRatio); sdkErr != nil {
		return sdkErr.Result()
	}
	logger.Debug(fmt.Sprintf("successfully handleMsgDexSet: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDexSave(ctx sdk.Context, keeper Keeper, msg types.MsgDexSave, logger log.Logger) sdk.Result {
	tradePair := keeper.GetTradePair(ctx, msg.Product)
	if tradePair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to save because non-exist product: %s", msg.Product)).Result()
	}
	if !tradePair.Owner.Equals(msg.Address) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to save because %s is not the owner of product:%s", msg.Address.String(), msg.Product)).Result()
	}

	if sdkErr := keeper.DexSave(ctx, msg.Address, msg.Product, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}
	logger.Debug(fmt.Sprintf("successfully handleMsgDexSave: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDexReturn(ctx sdk.Context, keeper Keeper, msg types.MsgDexReturn, logger log.Logger) sdk.Result {
	tokenPair := keeper.GetTradePair(ctx, msg.Product)
	if tokenPair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to return because non-exist product: %s", msg.Product)).Result()
	}
	if !tokenPair.Owner.Equals(msg.Address) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to return because %s is not the owner of product:%s", msg.Address.String(), msg.Product)).Result()
	}

	if sdkErr := keeper.DexReturn(ctx, msg.Address, msg.Product, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}
	logger.Debug(fmt.Sprintf("successfully handleMsgDexReturn: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDeposit, logger log.Logger) (result sdk.Result) {
	tradePair := keeper.GetTradePair(ctx, msg.Product)
	if nil == tradePair {
		return types.ErrInvalidTradePair(types.Codespace, fmt.Sprintf("no such trade pair %s", msg.Product)).Result()
	}

	if err := keeper.Deposit(ctx, msg.Address, msg.Product, msg.Amount); err != nil {
		return err.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgDeposit: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute("amount", msg.Amount.String()),
			sdk.NewAttribute("product", msg.Product),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdraw(ctx sdk.Context, keeper Keeper, msg types.MsgWithdraw, logger log.Logger) sdk.Result {
	tokenPair := keeper.GetTradePair(ctx, msg.Product)
	if tokenPair == nil {
		return types.ErrInvalidTradePair(types.Codespace, fmt.Sprintf("no such trade pair %s", msg.Product)).Result()
	}

	if err := keeper.Withdraw(ctx, msg.Address, msg.Product, msg.Amount); err != nil {
		return err.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgWithdraw: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBorrow(ctx sdk.Context, keeper Keeper, msg types.MsgBorrow, logger log.Logger) (result sdk.Result) {
	tradePair := keeper.GetTradePair(ctx, msg.Product)
	if nil == tradePair {
		return types.ErrInvalidTradePair(types.Codespace, fmt.Sprintf("no such trade pair %s", msg.Product)).Result()
	}
	if msg.Leverage.GT(tradePair.MaxLeverage) {
		return types.ErrInvalidLeverage(types.Codespace, fmt.Sprintf("%s is more than the product max leverage %s", msg.Leverage, tradePair.MaxLeverage)).Result()
	}
	if err := keeper.Borrow(ctx, msg.Address, tradePair, msg.Amount, msg.Leverage); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
			sdk.NewAttribute("borrowed", sdk.NewCoins(msg.Amount).MulDec(msg.Leverage.Sub(sdk.NewDec(1))).String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgRefund(ctx sdk.Context, keeper Keeper, msg types.MsgRefund, logger log.Logger) (result sdk.Result) {
	tradePair := keeper.GetTradePair(ctx, msg.Product)
	if nil == tradePair {
		return types.ErrInvalidTradePair(types.Codespace, fmt.Sprintf("no such trade pair %s", msg.Product)).Result()
	}

	account := keeper.GetAccount(ctx, msg.Address, tradePair.Name)
	if account == nil {
		return types.ErrAccountNotExist(types.Codespace, fmt.Sprintf("failed to refund")).Result()
	}

	if account.Borrowed.AmountOf(msg.Amount.Denom).IsZero() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("refund amount:%s mismatch borrowed coins:%s", msg.Amount.String(), account.Borrowed.String())).Result()
	}

	actualRefund := keeper.Refund(ctx, account, msg.Address, tradePair, msg.Amount)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
			sdk.NewAttribute("actualRefund", actualRefund.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
