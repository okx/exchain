package keeper

import (
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/types"
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
		return nil, types.ErrCodeCodeEmptyDelegationVoteValidator()
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

// FundCommunityPool allows an account to directly fund the community fund pool.
// The amount is first added to the distribution module account and then directly
// added to the pool. An error is returned if the amount cannot be sent to the
// module account.
func (k Keeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
		return err
	}

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	k.SetFeePool(ctx, feePool)

	return nil
}

func (k Keeper) CheckDistributionProposalValid(ctx sdk.Context) bool {
	return tmtypes.HigherThanVenus3(ctx.BlockHeight()) && k.HasInitAllocateValidator(ctx)
}
