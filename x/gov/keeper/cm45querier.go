package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/common"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov/types"
	paramstypes "github.com/okex/exchain/x/params/types"
)

func cm45QueryProposal(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryProposalParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	proposal, ok := keeper.GetProposal(ctx, params.ProposalID)
	if !ok {
		return nil, types.ErrUnknownProposal(params.ProposalID)
	}

	// Here is for compatibility with the standard cosmos REST API.
	// Note: The Height field in OKC's ParameterChangeProposal will be discarded.
	if pcp, ok := proposal.Content.(paramstypes.ParameterChangeProposal); ok {
		innerContent := pcp.GetParameterChangeProposal()
		newProposal := types.WrapProposalForCosmosAPI(proposal, innerContent)
		proposal = newProposal
	}

	if p, ok := proposal.Content.(evmtypes.ManageContractMethodBlockedListProposal); ok {
		p.FixShortAddr()
		newProposal := types.WrapProposalForCosmosAPI(proposal, p)
		proposal = newProposal
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposal)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func cm45QueryProposals(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryProposalsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	proposals := keeper.GetProposalsFiltered(ctx, params.Voter, params.Depositor, params.ProposalStatus, params.Limit)
	cosmosProposals := make([]types.Proposal, 0, len(proposals))
	for _, proposal := range proposals {
		if pcp, ok := proposal.Content.(paramstypes.ParameterChangeProposal); ok {
			// Here is for compatibility with the standard cosmos REST API.
			// Note: The Height field in OKC's ParameterChangeProposal will be discarded.
			innerContent := pcp.GetParameterChangeProposal()
			newProposal := types.WrapProposalForCosmosAPI(proposal, innerContent)
			cosmosProposals = append(cosmosProposals, newProposal)
		} else if p, ok := proposal.Content.(evmtypes.ManageContractMethodBlockedListProposal); ok {
			p.FixShortAddr()
			newProposal := types.WrapProposalForCosmosAPI(proposal, p)
			cosmosProposals = append(cosmosProposals, newProposal)
		} else {
			cosmosProposals = append(cosmosProposals, proposal)
		}
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, cosmosProposals)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
