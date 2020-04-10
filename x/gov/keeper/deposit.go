package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okex/okchain/x/gov/types"
)

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(
	ctx sdk.Context, proposalID uint64, depositorAddr sdk.Address,
) (types.Deposit, bool) {
	var depositNum, i uint64
	depositNum = keeper.getProposalDepositCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	for i = 0; i < depositNum; i++ {
		depositKey := types.DepositKey(proposalID, i)
		bz := store.Get(depositKey)
		if bz == nil {
			continue
		}
		var deposit types.Deposit
		keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
		if deposit.Depositor.Equals(depositorAddr) {
			return deposit, true
		}
	}
	return types.Deposit{}, false
}

// SetDeposit sets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(deposit.ProposalID, deposit.DepositID), bz)
	keeper.setProposalDepositCnt(ctx, deposit.ProposalID)
}

func (keeper Keeper) changeOldDeposit(ctx sdk.Context, proposalID uint64, deposit types.Deposit) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(proposalID, deposit.DepositID), bz)
}

func tryEnterVotingPeriod(
	ctx sdk.Context, keeper Keeper, proposal *types.Proposal, depositAmount sdk.DecCoins, eventType string,
) {
	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount)
	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	var minDeposit sdk.DecCoins
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
	ctx sdk.Context, keeper Keeper, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.DecCoins,
) {
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
		keeper.changeOldDeposit(ctx, proposalID, deposit)
	} else {
		newDeposit := types.Deposit{
			ProposalID: proposalID,
			Depositor:  depositorAddr,
			Amount:     depositAmount,
			DepositID:  keeper.getProposalDepositCnt(ctx, proposalID),
		}
		keeper.SetDeposit(ctx, newDeposit)
	}
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(
	ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress,
	depositAmount sdk.DecCoins, eventType string,
) sdk.Error {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal(keeper.Codespace(), proposalID)
	}

	// Check if proposal is still depositable
	if proposal.Status != types.StatusDepositPeriod {
		return types.ErrInvalidateProposalStatus(keeper.Codespace(),
			fmt.Sprintf("The status of proposal %d is in %s can not be deposited.",
				proposal.ProposalID, proposal.Status))
	}
	depositCoinsAmount := depositAmount
	// update the governance module's account coins pool
	err := keeper.SupplyKeeper().SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositCoinsAmount)
	if err != nil {
		return err
	}

	// try enter voting period according to proposal's total deposit
	tryEnterVotingPeriod(ctx, keeper, &proposal, depositAmount, eventType)

	// Add or update deposit object
	updateDeposit(ctx, keeper, proposalID, depositorAddr, depositAmount)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdkGovTypes.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(sdkGovTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// GetDeposits returns all the deposits from a proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits types.Deposits) {
	var depositNum, i uint64
	depositNum = keeper.getProposalDepositCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	for i = 0; i < depositNum; i++ {
		var deposit types.Deposit
		bz := store.Get(types.DepositKey(proposalID, i))
		if bz == nil {
			continue
		}
		keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
		deposits = append(deposits, deposit)
	}
	return deposits
}

// RefundDeposits refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	deposits := keeper.GetDeposits(ctx, proposalID)
	for i := 0; i < len(deposits); i++ {
		deposit := deposits[i]
		err := keeper.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, deposit.Depositor,
			deposit.Amount)
		if err != nil {
			panic(err)
		}
		keeper.deleteDeposit(ctx, proposalID, deposit.DepositID)
	}
	keeper.deleteDepositsCnt(ctx, proposalID)
}

// DistributeDeposits distributes and deletes all the deposits on a specific proposal
func (keeper Keeper) DistributeDeposits(ctx sdk.Context, proposalID uint64) {
	deposits := keeper.GetDeposits(ctx, proposalID)
	for i := 0; i < len(deposits); i++ {
		deposit := deposits[i]
		err := keeper.SupplyKeeper().SendCoinsFromModuleToModule(ctx, types.ModuleName, keeper.feeCollectorName,
			deposit.Amount)
		if err != nil {
			panic(err)
		}
		keeper.deleteDeposit(ctx, proposalID, deposit.DepositID)
	}
	keeper.deleteDepositsCnt(ctx, proposalID)
}

func (keeper Keeper) deleteDeposit(ctx sdk.Context, proposalID, depositID uint64) {
	store := ctx.KVStore(keeper.StoreKey())
	store.Delete(types.DepositKey(proposalID, depositID))
}

func (keeper Keeper) deleteDepositsCnt(ctx sdk.Context, proposalID uint64) {
	if depositsIterator := keeper.GetDeposits(ctx, proposalID); len(depositsIterator) == 0 {
		store := ctx.KVStore(keeper.StoreKey())
		store.Delete(types.DepositCntKey(proposalID))
	}
}

// setProposalDepositCnt save new count of deposit for a specific proposalID
func (keeper Keeper) setProposalDepositCnt(ctx sdk.Context, proposalID uint64) {
	cnt := keeper.getProposalDepositCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(cnt + 1)
	store.Set(types.DepositCntKey(proposalID), bz)
}

func (keeper Keeper) getProposalDepositCnt(ctx sdk.Context, proposalID uint64) uint64 {
	var cnt uint64
	store := ctx.KVStore(keeper.StoreKey())
	cntKey := types.DepositCntKey(proposalID)
	bz := store.Get(cntKey)
	if bz == nil {
		bz = keeper.Cdc().MustMarshalBinaryLengthPrefixed(0)
		store.Set(cntKey, bz)
		return 0
	}
	keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &cnt)
	return cnt
}
