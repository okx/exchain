package gov

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/gov/keeper"
	"github.com/okex/exchain/x/gov/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID uint64            `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Deposits           Deposits          `json:"deposits" yaml:"deposits"`
	Votes              Votes             `json:"votes" yaml:"votes"`
	Proposals          []Proposal        `json:"proposals" yaml:"proposals"`
	WaitingProposals   map[string]uint64 `json:"waiting_proposals" yaml:"waiting_proposals"`
	DepositParams      DepositParams     `json:"deposit_params" yaml:"deposit_params"`
	VotingParams       VotingParams      `json:"voting_params" yaml:"voting_params"`
	TallyParams        TallyParams       `json:"tally_params" yaml:"tally_params"`
}

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	var minDeposit = sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}
	return GenesisState{
		StartingProposalID: 1,
		Proposals:          []types.Proposal{},
		DepositParams: DepositParams{
			MinDeposit:       minDeposit,
			MaxDepositPeriod: time.Hour * 24,
		},
		VotingParams: VotingParams{
			VotingPeriod: time.Hour * 72,
		},
		TallyParams: TallyParams{
			Quorum:          sdk.NewDecWithPrec(334, 3),
			Threshold:       sdk.NewDecWithPrec(5, 1),
			Veto:            sdk.NewDecWithPrec(334, 3),
			YesInVotePeriod: sdk.NewDecWithPrec(667, 3),
		},
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) equal(data2 GenesisState) bool {
	b1 := types.ModuleCdc.MustMarshalBinaryBare(data)
	b2 := types.ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) isEmpty() bool {
	emptyGenState := GenesisState{}
	return data.equal(emptyGenState)
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote Threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.TallyParams.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote Veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	quorum := data.TallyParams.Quorum
	if quorum.IsNegative() || quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote Quorum should be positive and less or equal to one, is %s",
			threshold.String())
	}

	yesInVotePeriod := data.TallyParams.YesInVotePeriod
	if yesInVotePeriod.IsNegative() || yesInVotePeriod.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote YesInVotePeriod should be positive and less or equal to one, is %s",
			threshold.String())
	}

	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k keeper.Keeper, supplyKeeper keeper.SupplyKeeper, data GenesisState) {
	k.SetProposalID(ctx, data.StartingProposalID)
	k.SetDepositParams(ctx, data.DepositParams)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	// check if the deposits pool account exists
	moduleAcc := k.GetGovernanceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	var totalDeposits sdk.SysCoins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote.ProposalID, vote)
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.ProposalID, proposal.DepositEndTime)
		case StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		}
		k.SetProposal(ctx, proposal)
	}

	for proposalIDStr, height := range data.WaitingProposals {
		proposalID, err := strconv.ParseUint(proposalIDStr, 10, 64)
		if err != nil {
			panic(err)
		}
		k.InsertWaitingProposalQueue(ctx, height, proposalID)
	}

	// add coins if not provided on genesis
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(totalDeposits); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	startingProposalID, err := k.GetProposalID(ctx)
	if err != nil {
		panic(err)
	}
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)

	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusNil, 0)

	var proposalsDeposits Deposits
	var proposalsVotes Votes
	for _, proposal := range proposals {
		deposits := k.GetDeposits(ctx, proposal.ProposalID)
		proposalsDeposits = append(proposalsDeposits, deposits...)

		votes := k.GetVotes(ctx, proposal.ProposalID)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	waitingProposals := make(map[string]uint64)
	k.IterateAllWaitingProposals(ctx, func(proposal types.Proposal, proposalID, height uint64) (stop bool) {
		waitingProposals[strconv.FormatUint(proposalID, 10)] = height
		return false
	})

	return GenesisState{
		StartingProposalID: startingProposalID,
		Deposits:           proposalsDeposits,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		WaitingProposals:   waitingProposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}
