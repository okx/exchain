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
		case types.MsgDexSave:
			handlerFun = func() sdk.Result {
				return handleMsgDexSave(ctx, k, msg, logger)
			}
		case types.MsgDeposit:
			handlerFun = func() sdk.Result {
				return handleMsgMarginDeposit(ctx, k, msg, logger)
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

	if sdkErr := keeper.Deposit(ctx, msg.Address, msg.Product, msg.Amount); sdkErr != nil {
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
	if sdkErr := keeper.Withdraw(ctx, msg.Product, msg.Address, msg.Amount); sdkErr != nil {
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

func handleMsgDexSave(ctx sdk.Context, keeper Keeper, msg types.MsgDexSave, logger log.Logger) sdk.Result {
	tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, msg.Product)
	if tokenPair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to save because non-exist product: %s", msg.Product)).Result()
	}
	if !tokenPair.Owner.Equals(msg.Address) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to save because %s is not the owner of product:%s", msg.Address.String(), msg.Product)).Result()
	}

	if sdkErr := keeper.Save(ctx, msg.Address, msg.Product, msg.Amount); sdkErr != nil {
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

// handle<Action> does x
func handleMsgMarginDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDeposit, logger log.Logger) (result sdk.Result) {
	tradePair := keeper.GetMarginTradePair(ctx, msg.Product)
	if nil == tradePair {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because non-exist product: %s", tradePair.Name())).Result()
	}

	if err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, types.ModuleName, msg.Amount); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", msg.Amount.String())).Result()
	}

	//marginAcc := types.GetMarginAccount(msg.Address.String())

	keeper.SetAccountAssetOnProduct(ctx, msg.Address, msg.Product, msg.Amount)
	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute("deposit amount", msg.Amount.String()),
			sdk.NewAttribute("deposit product", msg.Product),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//// handle<Action> does x
//func handleMsgMarginDeposit<Action>(ctx sdk.Context, k Keeper, msg Msg<Action>) (*sdk.Result, error) {
//	err := k.<Action>(ctx, msg.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	// TODO: Define your msg events
//	ctx.EventManager().EmitEvent(
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
//			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
//		),
//	)
//
//	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
//}
