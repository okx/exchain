package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"

	comm "github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/gov"
	govrest "github.com/okx/okbchain/x/gov/client/rest"
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
