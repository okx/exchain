//nolint
package v0_9

import (
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
)

// Migrate accepts exported genesis state from v0.34 and migrates it to v0.36
// genesis state. This migration flattens the deposits and votes and updates the
// proposal content to the new
func Migrate(oldGenState v08gov.GenesisState) GenesisState {
	var deposits Deposits
	for _, deposit := range oldGenState.Deposits {
		deposits = append(deposits, Deposit{
			deposit.ProposalID,
			deposit.Depositor,
			deposit.Amount,
			deposit.DepositID,
		})
	}

	var votes Votes
	for _, vote := range oldGenState.Votes {
		votes = append(votes, Vote{
			vote.ProposalID,
			vote.Voter,
			vote.Option,
			vote.VoteID,
		})
	}

	var proposals []Proposal
	for _, proposal := range oldGenState.Proposals {
		if proposal.GetProposalType() == v08gov.ProposalTypeDexList {
			continue
		}
		newProposal := Proposal{
			Content:          migrateContent(proposal),
			ProposalID:       proposal.GetProposalID(),
			Status:           proposal.GetStatus(),
			FinalTallyResult: proposal.GetFinalTallyResult(),
			SubmitTime:       proposal.GetSubmitTime(),
			DepositEndTime:   proposal.GetDepositEndTime(),
			TotalDeposit:     proposal.GetTotalDeposit(),
			VotingStartTime:  proposal.GetVotingStartTime(),
			VotingEndTime:    proposal.GetVotingEndTime(),
		}
		proposals = append(proposals, newProposal)
	}

	depositParam, voteParam, tallyParams, tendermintParam := migrateParams(oldGenState.Params)

	return GenesisState{
		oldGenState.StartingProposalID,
		deposits, votes, proposals,
		depositParam, voteParam, tallyParams,
		tendermintParam,
	}
}

func migrateParams(oldParams v08gov.GovParams) (sdkGovTypes.DepositParams, sdkGovTypes.VotingParams,
	sdkGovTypes.TallyParams, sdkGovTypes.TendermintParams) {

	depositParam := sdkGovTypes.DepositParams{
		oldParams.TextMinDeposit,
		oldParams.TextMaxDepositPeriod,
	}

	voteParam := sdkGovTypes.VotingParams{
		oldParams.TextVotingPeriod,
	}

	tallyParams := sdkGovTypes.TallyParams{
		oldParams.Quorum,
		oldParams.Threshold,
		oldParams.Veto,
		oldParams.YesInVotePeriod,
	}

	tendermintParam := sdkGovTypes.TendermintParams{
		int(oldParams.MaxTxNumPerBlock),
	}

	return depositParam, voteParam, tallyParams, tendermintParam
}

func migrateContent(proposal v08gov.Proposal) (content sdkGovTypes.Content) {
	switch proposal.GetProposalType() {
	case v08gov.ProposalTypeText:
		return sdkGovTypes.NewTextProposal(proposal.GetTitle(), proposal.GetDescription())
	case v08gov.ProposalTypeParameterChange:
		paramChange, ok := proposal.(*v08gov.ParameterProposal)
		if !ok {
			panic("interface proposal failed to convert to ParameterProposal")
		}
		return ParameterChangeProposal{
			sdkparams.ParameterChangeProposal{
				paramChange.Title,
				paramChange.Description,
				convertParams(paramChange.Params),
			},
			uint64(paramChange.Height),
		}
	case v08gov.ProposalTypeAppUpgrade:
		appUpgrade, ok := proposal.(*v08gov.AppUpgradeProposal)
		if !ok {
			panic("interface proposal failed to convert to AppUpgradeProposal")
		}
		return AppUpgradeProposal{
			appUpgrade.Title,
			appUpgrade.Description,
			appUpgrade.ProtocolDefinition,
		}
	default:
		return nil
	}
}

func convertParams(v08params v08gov.Params) []sdkparams.ParamChange {
	var v09params []sdkparams.ParamChange
	for _, param := range v08params {
		v09params = append(v09params, sdkparams.ParamChange{
			Subspace: param.Subspace,
			Key:      param.Key,
			Value:    param.Value,
		})
	}
	return v09params
}
