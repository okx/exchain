package keeper

import (
	"fmt"
	stakingexported "github.com/okex/exchain/x/staking/exported"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	"github.com/okex/exchain/x/distribution/types"
)

// HandleCommunityPoolSpendProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolSpendProposal(ctx sdk.Context, k Keeper, p types.CommunityPoolSpendProposal) error {
	if k.blacklistedAddrs[p.Recipient.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is blacklisted from receiving external funds", p.Recipient)
	}

	err := k.distributeFromFeePool(ctx, p.Amount, p.Recipient)
	if err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("transferred %s from the community pool to recipient %s", p.Amount, p.Recipient))
	return nil
}

// distributeFromFeePool distributes funds from the distribution module account to
// a receiver address while updating the community pool
func (k Keeper) distributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	feePool := k.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool, negative := feePool.CommunityPool.SafeSub(amount)
	if negative {
		return types.ErrBadDistribution()
	}
	feePool.CommunityPool = newPool

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiveAddr, amount)
	if err != nil {
		return err
	}

	k.SetFeePool(ctx, feePool)
	return nil
}

// HandleChangeDistributionTypeProposal is a handler for executing a passed change distribution type proposal
func HandleChangeDistributionTypeProposal(ctx sdk.Context, k Keeper, p types.ChangeDistributionTypeProposal) error {
	logger := k.Logger(ctx)

	//1.check if it's the same
	if k.GetDistributionType(ctx) == p.Type {
		logger.Debug(fmt.Sprintf("do nothing, same distribution type, %d", p.Type))
		return nil
	}

	//2. if on chain, iteration validators and init val which has not outstanding
	if p.Type == types.DistributionTypeOnChain {
		if !k.HasInitAllocateValidator(ctx) {
			k.SetInitAllocateValidator(ctx, true)
			k.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingexported.ValidatorI) (stop bool) {
				if validator != nil {
					k.initValidatorWithoutOutstanding(ctx, validator)
				}
				return false
			})
		}
	}

	//3. set it
	k.SetDistributionType(ctx, p.Type)

	return nil
}
