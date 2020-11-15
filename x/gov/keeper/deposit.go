package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/gov/types"
)

// SetDeposit sets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(deposit.ProposalID, deposit.Depositor), bz)
}

func tryEnterVotingPeriod(
	ctx sdk.Context, keeper Keeper, proposal *types.Proposal, depositAmount sdk.SysCoins, eventType string,
) {
	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount...)
	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	var minDeposit sdk.SysCoins
	if !keeper.proposalHandlerRouter.HasRoute(proposal.ProposalRoute()) {
		minDeposit = keeper.GetDepositParams(ctx).MinDeposit
	} else {
		phr := keeper.proposalHandlerRouter.GetRoute(proposal.ProposalRoute())
		minDeposit = phr.GetMinDeposit(ctx, proposal.Content)
	}

	if proposal.Status == types.StatusDepositPeriod && proposal.TotalDeposit.IsAllGTE(minDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
		proposal.DepositEndTime = ctx.BlockHeader().Time
	}
	keeper.SetProposal(ctx, *proposal)

	if activatedVotingPeriod {
		// execute the logic when the deposit period is passed
		if !keeper.ProposalHandlerRouter().HasRoute(proposal.Content.ProposalRoute()) {
			keeper.AfterDepositPeriodPassed(ctx, *proposal)
		} else {
			proposalHandler := keeper.ProposalHandlerRouter().GetRoute(proposal.Content.ProposalRoute())
			proposalHandler.AfterDepositPeriodPassed(ctx, *proposal)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
			),
		)
	}
}

func updateDeposit(
	ctx sdk.Context, keeper Keeper, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.SysCoins,
) {
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount...)
	} else {
		deposit = types.Deposit{
			ProposalID: proposalID,
			Depositor:  depositorAddr,
			Amount:     depositAmount,
		}
	}
	keeper.SetDeposit(ctx, deposit)
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(
	ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress,
	depositAmount sdk.SysCoins, eventType string,
) sdk.Error {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal(keeper.codespace, proposalID)
	}

	// Check if proposal is still depositable
	if proposal.Status != types.StatusDepositPeriod {
		return types.ErrInvalidateProposalStatus(keeper.codespace,
			fmt.Sprintf("The status of proposal %d is in %s can not be deposited.",
				proposal.ProposalID, proposal.Status))
	}
	depositCoinsAmount := depositAmount
	// update the governance module's account coins pool
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositCoinsAmount)
	if err != nil {
		return err
	}

	// try enter voting period according to proposal's total deposit
	tryEnterVotingPeriod(ctx, keeper, &proposal, depositAmount, eventType)

	// Add or update deposit object
	updateDeposit(ctx, keeper, proposalID, depositorAddr, depositAmount)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// RefundDeposits refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	deposits := keeper.GetDeposits(ctx, proposalID)
	for i := 0; i < len(deposits); i++ {
		deposit := deposits[i]
		err := keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, deposit.Depositor,
			deposit.Amount)
		if err != nil {
			panic(err)
		}
		keeper.deleteDeposit(ctx, proposalID, deposit.Depositor)
	}
}

// DistributeDeposits distributes and deletes all the deposits on a specific proposal
func (keeper Keeper) DistributeDeposits(ctx sdk.Context, proposalID uint64) {
	deposits := keeper.GetDeposits(ctx, proposalID)
	for i := 0; i < len(deposits); i++ {
		deposit := deposits[i]
		err := keeper.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, keeper.feeCollectorName,
			deposit.Amount)
		if err != nil {
			panic(err)
		}
		keeper.deleteDeposit(ctx, proposalID, deposit.Depositor)
	}
}

func (keeper Keeper) deleteDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.DepositKey(proposalID, depositorAddr))
}

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(proposalID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) setDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(proposalID, depositorAddr), bz)
}

// GetAllDeposits returns all the deposits from the store
func (keeper Keeper) GetAllDeposits(ctx sdk.Context) (deposits types.Deposits) {
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDeposits returns all the deposits from a proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDepositsIterator gets all the deposits on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetDepositsIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))
}

// DeleteDeposits deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		err := keeper.supplyKeeper.BurnCoins(ctx, types.ModuleName, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}
