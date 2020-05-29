package v0_8

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/common/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
)

func TestRegisterCodec(t *testing.T) {
	cdc := codec.New()
	RegisterCodec(cdc)
	var proposal Proposal
	proposalBytes := cdc.MustMarshalBinaryBare(&TextProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok := proposal.(*TextProposal)
	require.True(t, ok)

	proposal = nil
	proposalBytes = cdc.MustMarshalBinaryBare(&DexListProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok = proposal.(*DexListProposal)
	require.True(t, ok)

	proposal = nil
	proposalBytes = cdc.MustMarshalBinaryBare(&ParameterProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok = proposal.(*ParameterProposal)
	require.True(t, ok)

	proposal = nil
	proposalBytes = cdc.MustMarshalBinaryBare(&AppUpgradeProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok = proposal.(*AppUpgradeProposal)
	require.True(t, ok)
}

func TestBasicProposalImplement(t *testing.T) {
	basicProposal := &BasicProposal{}
	var proposal Proposal = basicProposal
	proposal.SetProposalID(1)
	proposal.SetTitle("text")
	proposal.SetDescription("text")
	proposal.SetProposalType(ProposalTypeText)
	proposal.SetStatus(StatusDepositPeriod)
	proposal.SetFinalTallyResult(TallyResult{})
	proposal.SetSubmitTime(time.Now())
	proposal.SetDepositEndTime(time.Now())
	proposal.SetTotalDeposit(sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.OneDec())})
	proposal.SetVotingStartTime(time.Time{})
	proposal.SetVotingEndTime(time.Time{})
	proposal.SetProtocolDefinition(proto.ProtocolDefinition{})
	basicProposal1 := &BasicProposal{
		proposal.GetProposalID(),
		proposal.GetTitle(),
		proposal.GetDescription(),
		proposal.GetProposalType(),
		proposal.GetStatus(),
		proposal.GetFinalTallyResult(),
		proposal.GetSubmitTime(),
		proposal.GetDepositEndTime(),
		proposal.GetTotalDeposit(),
		proposal.GetVotingStartTime(),
		proposal.GetVotingEndTime(),
	}
	require.Equal(t, basicProposal.String(), basicProposal1.String())
	require.Equal(t, proto.ProtocolDefinition{}, proposal.GetProtocolDefinition())
}

func TestProposalStatusImplement(t *testing.T) {
	// ProposalStatusFromString
	testCases := []struct {
		statusString string
		status       ProposalStatus
		err          error
	}{
		{"DepositPeriod", StatusDepositPeriod, nil},
		{"VotingPeriod", StatusVotingPeriod, nil},
		{"Passed", StatusPassed, nil},
		{"Rejected", StatusRejected, nil},
		{"", StatusNil, nil},
		{"xxx", ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", "xxx")},
	}
	for _, testCase := range testCases {
		status, err := ProposalStatusFromString(testCase.statusString)
		require.Equal(t, testCase.status, status)
		require.Equal(t, testCase.err, err)
	}

	testCases2 := []struct {
		statusString string
		status       ProposalStatus
	}{
		{"DepositPeriod", StatusDepositPeriod},
		{"VotingPeriod", StatusVotingPeriod},
		{"Passed", StatusPassed},
		{"Rejected", StatusRejected},
		{"", StatusNil},
	}
	for _, testCase := range testCases2 {
		require.Equal(t, testCase.status.String(), testCase.statusString)
	}

	proposalStatus := StatusDepositPeriod
	cdc := codec.New()
	statusBytes, err := cdc.MarshalJSON(proposalStatus)
	require.Nil(t, err)
	var proposalStatus2 ProposalStatus
	cdc.MustUnmarshalJSON(statusBytes, &proposalStatus2)
	require.Equal(t, proposalStatus, proposalStatus2)
}

func TestProposalKindImplement(t *testing.T) {
	testCases := []struct {
		kindString string
		kind       ProposalKind
		err        error
	}{
		{"Text", ProposalTypeText, nil},
		{"ParameterChange", ProposalTypeParameterChange, nil},
		{"AppUpgrade", ProposalTypeAppUpgrade, nil},
		{"DexList", ProposalTypeDexList, nil},
		{"xxx", ProposalKind(0xff), fmt.Errorf("'%s' is not a valid proposal type", "xxx")},
	}

	for _, testCase := range testCases {
		kind, err := ProposalTypeFromString(testCase.kindString)
		require.Equal(t, testCase.kind, kind)
		require.Equal(t, testCase.err, err)
	}

	testCases2 := []struct {
		kindString string
		kind       ProposalKind
	}{
		{"Text", ProposalTypeText},
		{"ParameterChange", ProposalTypeParameterChange},
		{"AppUpgrade", ProposalTypeAppUpgrade},
		{"DexList", ProposalTypeDexList},
		{"", ProposalTypeNil},
	}
	for _, testCase := range testCases2 {
		require.Equal(t, testCase.kind.String(), testCase.kindString)
	}

	proposalKind := ProposalTypeText
	cdc := codec.New()
	statusBytes, err := cdc.MarshalJSON(proposalKind)
	require.Nil(t, err)
	var proposalKind2 ProposalKind
	cdc.MustUnmarshalJSON(statusBytes, &proposalKind2)
	require.Equal(t, proposalKind, proposalKind2)
}

func TestVoteOptionImplement(t *testing.T) {
	testCases := []struct {
		optionString string
		option       VoteOption
		err          error
	}{
		{"Yes", OptionYes, nil},
		{"Abstain", OptionAbstain, nil},
		{"No", OptionNo, nil},
		{"NoWithVeto", OptionNoWithVeto, nil},
		{"xxx", VoteOption(0xff), fmt.Errorf("'%s' is not a valid vote option", "xxx")},
	}

	for _, testCase := range testCases {
		option, err := VoteOptionFromString(testCase.optionString)
		require.Equal(t, testCase.option, option)
		require.Equal(t, testCase.err, err)
	}

	testCases2 := []struct {
		optionString string
		option       VoteOption
	}{
		{"Yes", OptionYes},
		{"Abstain", OptionAbstain},
		{"No", OptionNo},
		{"NoWithVeto", OptionNoWithVeto},
		{"", OptionEmpty},
	}
	for _, testCase := range testCases2 {
		require.Equal(t, testCase.option.String(), testCase.optionString)
	}

	option := OptionYes
	cdc := codec.New()
	optionBytes, err := cdc.MarshalJSON(option)
	require.Nil(t, err)
	var option2 VoteOption
	cdc.MustUnmarshalJSON(optionBytes, &option2)
	require.Equal(t, option, option2)

	optionBytes, err = cdc.MarshalBinaryBare(option)
	require.Nil(t, err)
	err = cdc.UnmarshalBinaryBare(optionBytes, &option2)
	require.Nil(t, err)
	require.Equal(t, option, option2)

}
