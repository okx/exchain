package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okexchain/x/distribution/keeper"
	"github.com/okex/okexchain/x/distribution/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
)

// NewHandler manages all distribution tx
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)

		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)

		default:
			return nil, types.ErrUnknownDistributionMsgType()
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)
func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.Keeper) (*sdk.Result, error) {
	err := k.SetWithdrawAddr(ctx, msg.DelegatorAddress, msg.WithdrawAddress)
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

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.MsgWithdrawValidatorCommission, k keeper.Keeper) (*sdk.Result, error) {
	_, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ERRWithdrawValidatorCommissionFailed()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewCommunityPoolSpendProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content *govtypes.Proposal) error {
		switch c := content.Content.(type) {
		case types.CommunityPoolSpendProposal:
			return keeper.HandleCommunityPoolSpendProposal(ctx, k, c)

		default:
			return types.ErrUnknownDistributionCommunityPoolProposaType()
		}
	}
}
