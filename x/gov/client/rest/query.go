package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	common "github.com/okex/exchain/x/gov/client/common"
	gcutils "github.com/okex/exchain/x/gov/client/utils"
	"github.com/okex/exchain/x/gov/types"
)

func depositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		var req DepositReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgDeposit(req.Depositor, proposalID, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func voteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		var req VoteReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		voteOption, err := types.VoteOptionFromString(gcutils.NormalizeVoteOption(req.Option))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgVote(req.Voter, proposalID, voteOption)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[RestParamsType]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := common.QueryParams(cliCtx, paramType)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var proposal types.Proposal
		cliCtx.Codec.MustUnmarshalJSON(res, &proposal)
		cm45p := proposal.ToCM45Proposal()
		wrappedProposal := types.NewWrappedProposal(*cm45p)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, wrappedProposal)
	}
}

func queryDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var proposal types.Proposal
		if err := cliCtx.Codec.UnmarshalJSON(res, &proposal); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// For inactive proposals we must query the txs directly to get the deposits
		// as they're no longer in state.
		propStatus := proposal.Status
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryDepositsByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData("custom/gov/deposits", bz)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var deposits types.Deposits
		cliCtx.Codec.MustUnmarshalJSON(res, &deposits)
		wrappedDeposits := types.NewWrappedDeposits(deposits)
		rest.PostProcessResponse(w, cliCtx, wrappedDeposits)
	}
}

func queryProposerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, err := gcutils.QueryProposerByTxQuery(cliCtx, proposalID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDepositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechDepositorAddr := vars[RestDepositor]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechDepositorAddr) == 0 {
			err := errors.New("depositor address required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryDepositParams(proposalID, depositorAddr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/deposit", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var deposit types.Deposit
		if err := cliCtx.Codec.UnmarshalJSON(res, &deposit); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// For an empty deposit, either the proposal does not exist or is inactive in
		// which case the deposit would be removed from state and should be queried
		// for directly via a txs query.
		if deposit.Empty() {
			bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposalParams(proposalID))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			res, _, err = cliCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = gcutils.QueryDepositByTxQuery(cliCtx, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryVoteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechVoterAddr := vars[RestVoter]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechVoterAddr) == 0 {
			err := errors.New("voter address required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryVoteParams(proposalID, voterAddr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/vote", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vote types.Vote
		if err := cliCtx.Codec.UnmarshalJSON(res, &vote); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// For an empty vote, either the proposal does not exist or is inactive in
		// which case the vote would be removed from state and should be queried for
		// directly via a txs query.
		if vote.Empty() {
			bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposalParams(proposalID))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			res, _, err = cliCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = gcutils.QueryVoteByTxQuery(cliCtx, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryVotesOnProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var proposal types.Proposal
		if err := cliCtx.Codec.UnmarshalJSON(res, &proposal); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// For inactive proposals we must query the txs directly to get the votes
		// as they're no longer in state.
		propStatus := proposal.Status
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryVotesByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData("custom/gov/votes", bz)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryProposalsWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechVoterAddr := r.URL.Query().Get(RestVoter)
		bechDepositorAddr := r.URL.Query().Get(RestDepositor)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		params := types.QueryProposalsParams{}

		if len(bechVoterAddr) != 0 {
			voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Voter = voterAddr
		}

		if len(bechDepositorAddr) != 0 {
			depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Depositor = depositorAddr
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err := types.ProposalStatusFromString(gcutils.NormalizeProposalStatus(strProposalStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ProposalStatus = proposalStatus
		}
		if len(strNumLimit) != 0 {
			numLimit, ok := rest.ParseUint64OrReturnBadRequest(w, strNumLimit)
			if !ok {
				return
			}
			params.Limit = numLimit
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/proposals", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryProposalsWithParameterCM45Fn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechVoterAddr := r.URL.Query().Get(RestVoter)
		bechDepositorAddr := r.URL.Query().Get(RestDepositor)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)
		strReverse := r.URL.Query().Get("pagination.reverse")

		params := types.QueryProposalsParams{}

		if len(bechVoterAddr) != 0 {
			voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Voter = voterAddr
		}

		if len(bechDepositorAddr) != 0 {
			depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Depositor = depositorAddr
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err := types.ProposalStatusFromString(gcutils.NormalizeProposalStatus(strProposalStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ProposalStatus = proposalStatus
		}
		if len(strNumLimit) != 0 {
			numLimit, ok := rest.ParseUint64OrReturnBadRequest(w, strNumLimit)
			if !ok {
				return
			}
			params.Limit = numLimit
		}
		needReverse := false
		if len(strReverse) != 0 {
			reverse, err := strconv.ParseBool(strReverse)
			if err != nil {
				return
			}
			needReverse = reverse
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/proposals", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var proposals []types.Proposal
		cliCtx.Codec.MustUnmarshalJSON(res, &proposals)
		if needReverse {
			for i, j := 0, len(proposals)-1; i < j; i, j = i+1, j-1 {
				proposals[i], proposals[j] = proposals[j], proposals[i]
			}
		}
		var cm45proposals []types.CM45Proposal
		for _, p := range proposals {
			cm45proposals = append(cm45proposals, *p.ToCM45Proposal())
		}
		wrappedProposals := types.NewWrappedProposals(cm45proposals)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, wrappedProposals)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryTallyOnProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/tally", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
