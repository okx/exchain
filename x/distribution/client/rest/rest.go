package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"

	comm "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/gov"
	govrest "github.com/okex/exchain/x/gov/client/rest"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

// CommunityPoolSpendProposalRESTHandler returns a CommunityPoolSpendProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func CommunityPoolSpendProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "community_pool_spend",
		Handler:  postCommunityPoolSpendProposalHandlerFn(cliCtx),
	}
}

func postCommunityPoolSpendProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommunityPoolSpendProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewCommunityPoolSpendProposal(req.Title, req.Description, req.Recipient, req.Amount)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// ChangeDistributionModelProposalRESTHandler returns a ChangeDistributionTypeProposal that exposes the community pool spend REST handler with a given sub-route.
func ChangeDistributionModelProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "change_distribution_type",
		Handler:  postChangeDistributionTypeProposalHandlerFn(cliCtx),
	}
}

func postChangeDistributionTypeProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeDistributionTypeProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewChangeDistributionModelProposal(req.Title, req.Description, req.Type)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
