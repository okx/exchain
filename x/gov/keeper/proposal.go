package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okex/okchain/x/gov/types"
)

// SubmitProposal creates new proposal given a content
func (keeper Keeper) SubmitProposal(ctx sdk.Context, content sdkGovTypes.Content) (sdkGovTypes.Proposal, sdk.Error) {
	if !keeper.Router().HasRoute(content.ProposalRoute()) {
		return sdkGovTypes.Proposal{}, types.ErrNoProposalHandlerExists(keeper.Codespace(), content)
	}

	proposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return sdkGovTypes.Proposal{}, err
	}
	// get the time now as the submit time
	submitTime := ctx.BlockHeader().Time
	// get params for special proposal
	var depositPeriod time.Duration
	if !keeper.proposalHandlerRouter.HasRoute(content.ProposalRoute()) {
		depositPeriod = keeper.GetDepositParams(ctx).MaxDepositPeriod
	} else {
		proposalParams := keeper.proposalHandlerRouter.GetRoute(content.ProposalRoute())
		depositPeriod = proposalParams.GetMaxDepositPeriod(ctx, content)
	}

	proposal := sdkGovTypes.NewProposal(ctx, keeper.totalPower(ctx), content, proposalID, submitTime,
		submitTime.Add(depositPeriod))

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	keeper.SetProposalID(ctx, proposalID+1)

	if keeper.proposalHandlerRouter.HasRoute(content.ProposalRoute()) {
		keeper.proposalHandlerRouter.GetRoute(content.ProposalRoute()).AfterSubmitProposalHandler(ctx, proposal)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return proposal, nil
}

// GetProposal get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (proposal sdkGovTypes.Proposal, ok bool) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return
	}
	keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal, true
}

// SetProposal set a proposal to store
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal sdkGovTypes.Proposal) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(types.ProposalKey(proposal.ProposalID), bz)
}

// DeleteProposal deletes a proposal from store
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.StoreKey())
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}
	keeper.RemoveFromInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	keeper.RemoveFromActiveProposalQueue(ctx, proposalID, proposal.VotingEndTime)
	store.Delete(types.ProposalKey(proposalID))
}

// GetProposals returns all the proposals from store
func (keeper Keeper) GetProposals(ctx sdk.Context) (proposals sdkGovTypes.Proposals) {
	keeper.IterateProposals(ctx, func(proposal sdkGovTypes.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

// GetProposalsFiltered get Proposals from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (keeper Keeper) GetProposalsFiltered(
	ctx sdk.Context, voterAddr sdk.AccAddress, depositorAddr sdk.AccAddress, status types.ProposalStatus,
	numLatest uint64,
) []sdkGovTypes.Proposal {

	maxProposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return []sdkGovTypes.Proposal{}
	}

	matchingProposals := []sdkGovTypes.Proposal{}

	if numLatest == 0 {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := keeper.GetVote(ctx, proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal, ok := keeper.GetProposal(ctx, proposalID)
		if !ok {
			continue
		}

		if types.ValidProposalStatus(status) && proposal.Status != status {
			continue
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

// GetProposalID gets the highest proposal ID
func (keeper Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(keeper.Codespace(), "initial proposal ID hasn't been set")
	}
	keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

// SetProposalID sets the proposal ID to gov store
func (keeper Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.ProposalIDKey, bz)
}

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, proposal *sdkGovTypes.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	var votingPeriod time.Duration
	if !keeper.proposalHandlerRouter.HasRoute(proposal.ProposalRoute()) {
		votingPeriod = keeper.GetVotingPeriod(ctx, proposal.Content)
	} else {
		phr := keeper.proposalHandlerRouter.GetRoute(proposal.ProposalRoute())
		votingPeriod = phr.GetVotingPeriod(ctx, proposal.Content)
	}
	// calculate the end time of voting
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = types.StatusVotingPeriod

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalID, proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
}
