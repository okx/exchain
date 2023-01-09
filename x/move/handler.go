package move

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/move/types"

	"github.com/okex/exchain/x/move/keeper"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx.SetEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgPublishMove:
			return handleMsgPublishMove(ctx, msg, k)

		default:
			return nil, nil
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)
func handleMsgPublishMove(ctx sdk.Context, msg types.MsgPublishMove, k keeper.Keeper) (*sdk.Result, error) {
	err := k.PublishMove(ctx, msg.DelegatorAddress, msg.MovePath)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
