package keeper

import (
	"fmt"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
	govTypes "github.com/okx/okbchain/x/gov/types"
)

// withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	val := k.stakingKeeper.Validator(ctx, valAddr)
	if val == nil {
		return nil, types.ErrCodeEmptyValidatorDistInfo()
	}
	logger := k.Logger(ctx)

	del := k.stakingKeeper.Delegator(ctx, delAddr)
	if del == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	valAddressArray := del.GetShareAddedValidatorAddresses()
	exist := false
	for _, valAddress := range valAddressArray {
		if valAddress.Equals(valAddr) {
			exist = true
			break
		}
	}
	if !exist {
		return nil, types.ErrCodeEmptyDelegationVoteValidator()
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, delAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
		),
	)

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)
	logger.Debug("WithdrawDelegationRewards", "Validator", valAddr, "Delegator", delAddr)
	return rewards, nil
}

// withdraw all rewards
func (k Keeper) WithdrawDelegationAllRewards(ctx sdk.Context, delAddr sdk.AccAddress) error {
	del := k.stakingKeeper.Delegator(ctx, delAddr)
	if del == nil {
		return types.ErrCodeEmptyDelegationDistInfo()
	}

	valAddressArray := del.GetShareAddedValidatorAddresses()
	if len(valAddressArray) == 0 {
		return types.ErrCodeEmptyDelegationVoteValidator()
	}

	logger := k.Logger(ctx)
	for _, valAddr := range valAddressArray {
		val := k.stakingKeeper.Validator(ctx, valAddr)
		if val == nil {
			return types.ErrCodeEmptyValidatorDistInfo()
		}
		// withdraw rewards
		rewards, err := k.withdrawDelegationRewards(ctx, val, delAddr)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawRewards,
				sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			),
		)

		// reinitialize the delegation
		k.initializeDelegation(ctx, valAddr, delAddr)
		logger.Debug("WithdrawDelegationAllRewards", "Validator", valAddr, "Delegator", delAddr)
	}

	return nil
}

// GetTotalRewards returns the total amount of fee distribution rewards held in the store
func (k Keeper) GetTotalRewards(ctx sdk.Context) (totalRewards sdk.DecCoins) {
	k.IterateValidatorOutstandingRewards(ctx,
		func(_ sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool) {
			totalRewards = totalRewards.Add(rewards...)
			return false
		},
	)

	return totalRewards
}

// SetGovKeeper sets keeper of gov
func (k *Keeper) SetGovKeeper(gk types.GovKeeper) {
	k.govKeeper = gk
}

// CheckMsgSubmitProposal validates MsgSubmitProposal
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	err := k.govKeeper.CheckMsgSubmitProposal(ctx, msg)
	if err != nil {
		return err
	}

	log := ctx.Logger()
	switch content := msg.Content.(type) {
	case types.WithdrawRewardEnabledProposal, types.ChangeDistributionTypeProposal, types.RewardTruncatePrecisionProposal:
		log.Debug(fmt.Sprintf("proposal content type: %T", content))
		if !k.stakingKeeper.IsValidator(ctx, msg.Proposer) {
			return types.ErrCodeProposerMustBeValidator()
		}
	case types.CommunityPoolSpendProposal:
		return nil
	default:
		return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized %s proposal content type: %T", types.DefaultCodespace, content))
	}

	return nil
}

// nolint
func (keeper Keeper) GetMinDeposit(ctx sdk.Context, content govTypes.Content) (minDeposit sdk.SysCoins) {
	return keeper.govKeeper.GetDepositParams(ctx).MinDeposit
}

// nolint
func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, content govTypes.Content) time.Duration {
	return keeper.govKeeper.GetDepositParams(ctx).MaxDepositPeriod
}

// nolint
func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, content govTypes.Content) time.Duration {
	return keeper.govKeeper.GetVotingParams(ctx).VotingPeriod
}

// nolint
func (k Keeper) AfterSubmitProposalHandler(_ sdk.Context, _ govTypes.Proposal) {}
func (k Keeper) AfterDepositPeriodPassed(_ sdk.Context, _ govTypes.Proposal)   {}
func (k Keeper) RejectedHandler(_ sdk.Context, _ govTypes.Content)             {}
func (k Keeper) VoteHandler(_ sdk.Context, _ govTypes.Proposal, _ govTypes.Vote) (string, sdk.Error) {
	return "", nil
}
