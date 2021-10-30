package slashing

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/slashing/internal/types"
)

// NewHandler creates an sdk.Handler for all the slashing type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgUnjail:
			return handleMsgUnjail(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg MsgUnjail, k Keeper) (*sdk.Result, error) {
	validator := k.GetStakingKeeper().Validator(ctx, msg.ValidatorAddr)
	if validator == nil {
		return nil, sdkerrors.Wrapf(ErrNoValidatorForAddress, "Unjail failed")
	}

	if validator.GetMinSelfDelegation().IsZero() {
		return nil, sdkerrors.Wrapf(ErrMissingSelfDelegation, "Unjail failed")
	}

	err := k.Unjail(ctx, msg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
