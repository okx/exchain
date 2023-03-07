package distribution

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/keeper"
	"github.com/okx/okbchain/x/distribution/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

// NewHandler manages all distribution tx
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx.SetEventManager(sdk.NewEventManager())

		if !k.GetWithdrawRewardEnabled(ctx) {
			return nil, types.ErrCodeDisabledWithdrawRewards()
		}

		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)

		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)

		case types.MsgWithdrawDelegatorReward:
			if k.CheckDistributionProposalValid(ctx) {
				return handleMsgWithdrawDelegatorReward(ctx, msg, k)
			}
			return nil, types.ErrUnknownDistributionMsgType()
		case types.MsgWithdrawDelegatorAllRewards:
			if k.CheckDistributionProposalValid(ctx) {
				return handleMsgWithdrawDelegatorAllRewards(ctx, msg, k)
			}
			return nil, types.ErrUnknownDistributionMsgType()

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
		return nil, err
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

func NewDistributionProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content *govtypes.Proposal) error {
		switch c := content.Content.(type) {
		case types.CommunityPoolSpendProposal:
			return keeper.HandleCommunityPoolSpendProposal(ctx, k, c)
		case types.ChangeDistributionTypeProposal:
			return keeper.HandleChangeDistributionTypeProposal(ctx, k, c)
		case types.WithdrawRewardEnabledProposal:
			return keeper.HandleWithdrawRewardEnabledProposal(ctx, k, c)
		case types.RewardTruncatePrecisionProposal:
			return keeper.HandleRewardTruncatePrecisionProposal(ctx, k, c)
		default:
			return types.ErrUnknownDistributionCommunityPoolProposaType()
		}
	}
}
