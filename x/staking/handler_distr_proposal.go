package staking

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/staking/keeper"
	"github.com/okx/okbchain/x/staking/types"
)

func handleMsgEditValidatorCommissionRate(ctx sdk.Context, msg types.MsgEditValidatorCommissionRate, k keeper.Keeper) (*sdk.Result, error) {
	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return nil, ErrNoValidatorFound(msg.ValidatorAddress.String())
	}

	commission, err := k.UpdateValidatorCommission(ctx, validator, msg.CommissionRate)
	if err != nil {
		return nil, err
	}

	// call the before-modification hook since we're about to update the commission
	k.BeforeValidatorModified(ctx, msg.ValidatorAddress)

	validator.Commission = commission

	k.SetValidator(ctx, validator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeEditValidatorCommissionRate,
			sdk.NewAttribute(types.AttributeKeyCommissionRate, msg.CommissionRate.String()),
		),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
