package gov

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/gov/keeper"
	"github.com/okex/exchain/x/gov/types"
)

func TestInitGenesisState(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)

	initialDeposit := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	deposits := types.Deposits{
		{ProposalID: 1, Depositor: keeper.Addrs[0], Amount: initialDeposit},
	}
	votes := types.Votes{
		{ProposalID: 2, Voter: keeper.Addrs[1], Option: types.OptionYes},
	}
	proposals := types.Proposals{
		types.Proposal{
			ProposalID:       1,
			Status:           StatusDepositPeriod,
			FinalTallyResult: EmptyTallyResult(sdk.ZeroDec()),
		},
		types.Proposal{
			ProposalID:       2,
			Status:           StatusVotingPeriod,
			FinalTallyResult: EmptyTallyResult(sdk.ZeroDec()),
		},
	}
	waitingProposals := map[string]uint64{"2": 1234}

	data := GenesisState{
		StartingProposalID: 3,
		Deposits:           deposits,
		Votes:              votes,
		Proposals:          proposals,
		WaitingProposals:   waitingProposals,
	}

	InitGenesis(ctx, gk, gk.SupplyKeeper(), data)
	// 0x00
	proposal0, ok := gk.GetProposal(ctx, data.Proposals[0].ProposalID)
	require.True(t, ok)
	proposal1, ok := gk.GetProposal(ctx, data.Proposals[1].ProposalID)
	require.True(t, ok)
	require.Equal(t, data.Proposals[0], proposal0)
	require.Equal(t, data.Proposals[1], proposal1)
	// 0x01
	var activeProposal types.Proposal
	gk.IterateActiveProposalsQueue(ctx, time.Now(), func(proposal types.Proposal,
	) (stop bool) {
		activeProposal = proposal
		return false
	})
	require.Equal(t, data.Proposals[1], activeProposal)
	// 0x02
	gk.IterateInactiveProposalsQueue(ctx, time.Now(), func(proposal types.Proposal,
	) (stop bool) {
		activeProposal = proposal
		return false
	})
	require.Equal(t, data.Proposals[0], activeProposal)
	// 0x03
	pid, err := gk.GetProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, data.Proposals[1].ProposalID+1, pid)
	// 0x10
	deposit, ok := gk.GetDeposit(ctx, data.Deposits[0].ProposalID, data.Deposits[0].Depositor)
	require.True(t, ok)
	require.Equal(t, data.Deposits[0], deposit)
	// 0x11
	// getProposalDepositCnt not public method
	// Referenced by other methods,such as GetDeposit and GetProposal
	// 0x20
	require.Equal(t, types.Votes(nil), gk.GetVotes(ctx, data.Proposals[0].ProposalID))
	require.Equal(t, data.Votes, gk.GetVotes(ctx, data.Proposals[1].ProposalID))
	// 0x21
	// getProposalVoteCnt not public method
	// Referenced by other methods,such as GetVotes
	// 0x30
	var waitingProposal types.Proposal
	gk.IterateWaitingProposalsQueue(ctx, 1234, func(proposal types.Proposal,
	) (stop bool) {
		waitingProposal = proposal
		return false
	})
	require.Equal(t, data.Proposals[1], waitingProposal)

	inactiveQueue := gk.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())
	activeQueue.Close()

	exportGenesis := ExportGenesis(ctx, gk)
	require.Equal(t, data.Deposits, exportGenesis.Deposits)
	require.Equal(t, data.Votes, exportGenesis.Votes)

	newCtx, _, newgk, _, _ := keeper.CreateTestInput(t, false, 1000)
	InitGenesis(newCtx, newgk, newgk.SupplyKeeper(), exportGenesis)
	// 0x00
	proposal0, ok = newgk.GetProposal(newCtx, exportGenesis.Proposals[0].ProposalID)
	require.True(t, ok)
	proposal1, ok = newgk.GetProposal(newCtx, exportGenesis.Proposals[1].ProposalID)
	require.True(t, ok)
	require.Equal(t, exportGenesis.Proposals[0], proposal0)
	require.Equal(t, exportGenesis.Proposals[1], proposal1)
	// 0x01
	newgk.IterateActiveProposalsQueue(newCtx, time.Now(), func(proposal types.Proposal,
	) (stop bool) {
		activeProposal = proposal
		return false
	})
	require.Equal(t, exportGenesis.Proposals[1], activeProposal)
	// 0x02
	newgk.IterateInactiveProposalsQueue(newCtx, time.Now(), func(proposal types.Proposal,
	) (stop bool) {
		activeProposal = proposal
		return false
	})
	require.Equal(t, exportGenesis.Proposals[0], activeProposal)
	// 0x03
	pid, err = newgk.GetProposalID(newCtx)
	require.NoError(t, err)
	require.Equal(t, exportGenesis.Proposals[1].ProposalID+1, pid)
	// 0x10
	deposit, ok = newgk.GetDeposit(newCtx, exportGenesis.Deposits[0].ProposalID, exportGenesis.Deposits[0].Depositor)
	require.True(t, ok)
	require.Equal(t, exportGenesis.Deposits[0], deposit)
	// 0x11
	// getProposalDepositCnt not public method
	// Referenced by other methods,such as GetDeposit and GetProposal
	// 0x20
	require.Equal(t, types.Votes(nil), newgk.GetVotes(newCtx, exportGenesis.Proposals[0].ProposalID))
	require.Equal(t, exportGenesis.Votes, newgk.GetVotes(newCtx, exportGenesis.Proposals[1].ProposalID))
	// 0x21
	// getProposalVoteCnt not public method
	// Referenced by other methods,such as GetVotes
	// 0x30
	newgk.IterateWaitingProposalsQueue(newCtx, 1234, func(proposal types.Proposal,
	) (stop bool) {
		waitingProposal = proposal
		return false
	})
	require.Equal(t, exportGenesis.Proposals[1], waitingProposal)

}

func TestValidateGenesis(t *testing.T) {
	data := GenesisState{}
	var err sdk.Error
	data.TallyParams.Threshold, err = sdk.NewDecFromStr("-23")
	require.Nil(t, err)
	require.NotNil(t, ValidateGenesis(data))

	data.TallyParams.Threshold = sdk.NewDecWithPrec(334, 3)
	data.TallyParams.Veto, err = sdk.NewDecFromStr("-23")
	require.Nil(t, err)
	require.NotNil(t, ValidateGenesis(data))

	data.TallyParams.Veto = sdk.NewDecWithPrec(334, 3)
	data.TallyParams.Quorum, err = sdk.NewDecFromStr("-23")
	require.Nil(t, err)
	require.NotNil(t, ValidateGenesis(data))

	data.TallyParams.Quorum = sdk.NewDecWithPrec(334, 3)
	data.TallyParams.YesInVotePeriod, err = sdk.NewDecFromStr("-23")
	require.Nil(t, err)
	require.NotNil(t, ValidateGenesis(data))

	data.TallyParams.YesInVotePeriod = sdk.NewDecWithPrec(334, 3)
	coin, err := sdk.NewDecFromStr("-23")
	require.Nil(t, err)
	data.DepositParams.MinDeposit = sdk.SysCoins{sdk.SysCoin{Denom: sdk.DefaultBondDenom, Amount: coin}}
	require.NotNil(t, ValidateGenesis(data))
}

func TestGenesisState_Equal(t *testing.T) {
	var minDeposit = sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}
	expected := GenesisState{
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
	require.True(t, expected.equal(DefaultGenesisState()))
}

func TestGenesisState_IsEmpty(t *testing.T) {
	require.True(t, GenesisState{}.isEmpty())
}
