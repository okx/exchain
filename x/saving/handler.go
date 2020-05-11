package saving

import (
	"fmt"

	"github.com/okex/okchain/x/common/perf"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the saving type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", ModuleName)

		var handlerFun func() sdk.Result
		var name string
		switch msg := msg.(type) {
		case MsgDeposit:
			name = "handleMsgDeposit"
			handlerFun = func() sdk.Result {
				return handleMsgDeposit(ctx, k, msg, logger)
			}
		case MsgWithdraw:
			name = "handleMsgWithDraw"
			handlerFun = func() sdk.Result {
				return handleMsgWithDraw(ctx, k, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized dex message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit, logger log.Logger) sdk.Result {
	// TODO: validate token supports for saving
	if sdkErr := keeper.Deposit(ctx, msg.Address, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgDeposit: "+
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

func handleMsgWithDraw(ctx sdk.Context, keeper Keeper, msg MsgWithdraw, logger log.Logger) sdk.Result {
	if sdkErr := keeper.Withdraw(ctx, msg.Address, msg.Amount); sdkErr != nil {
		return sdkErr.Result()
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgWithDraw: "+
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
