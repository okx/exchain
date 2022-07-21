package rest

import (
	"net/http"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"

	comm "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/gov"
	govrest "github.com/okex/exchain/x/gov/client/rest"
)

// ChangeDistributionTypeProposalRESTHandler returns a ChangeDistributionTypeProposal that exposes the change distribution type REST handler with a given sub-route.
func ChangeDistributionTypeProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
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

		content := types.NewChangeDistributionTypeProposal(req.Title, req.Description, req.Type)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// WithdrawRewardEnabledProposalRESTHandler returns a WithdrawRewardEnabledProposal that exposes the set withdraw reward enabled proposal REST handler with a given sub-route.
func WithdrawRewardEnabledProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "withdraw_reward_enabled",
		Handler:  postWithdrawRewardEnabledProposalHandlerFn(cliCtx),
	}
}

func postWithdrawRewardEnabledProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawRewardEnabledProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewWithdrawRewardEnabledProposal(req.Title, req.Description, req.Enabled)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
