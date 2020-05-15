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
		// TODO: Define your msg cases
		//
		//Example:
		// case Msg<Action>:
		// 	return handleMsg<Action>(ctx, k, msg)
		case types.MsgMarginDeposit:
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
func handleMsgDexDeposit(ctx sdk.Context, keeper Keeper, msg MsgDexDeposit, logger log.Logger) sdk.Result {
	if sdkErr := keeper.Deposit(ctx, msg.Address, msg.Amount); sdkErr != nil {
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

// handle<Action> does x
func handleMsgMarginDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgMarginDeposit, logger log.Logger) (result sdk.Result) {
	tradePair := keeper.GetMarginTradePair(ctx, msg.Product)
	if nil == tradePair {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because non-exist product: %s", tradePair.Name())).Result()
	}

	if err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, types.ModuleName, msg.Amount.ToCoins()); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", msg.Amount.String())).Result()
	}

	//marginAcc := types.GetMarginAccount(msg.Address.String())

	keeper.SetMarginDepositOnProduct(ctx, msg.Address, msg.Product, sdk.DecCoins{msg.Amount})
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
