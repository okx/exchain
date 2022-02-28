package simulation

import (
	"math/rand"

	simappparams "github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp/params"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	simtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/simulation"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/gov/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/simulation"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightMsgDeposit,
			simappparams.DefaultWeightTextProposal,
			SimulateTextProposalContent,
		),
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
	return types.NewTextProposal(
		simtypes.RandStringOfLength(r, 140),
		simtypes.RandStringOfLength(r, 5000),
	)
}
