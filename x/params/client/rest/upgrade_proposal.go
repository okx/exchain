package rest

import (
	"github.com/okex/exchain/x/params/types"
	"net/http"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/x/gov"
	govrest "github.com/okex/exchain/x/gov/client/rest"
)

// UpgradeProposalReq defines a upgrade proposal request body
type UpgradeProposalReq struct {
	BaseReq  rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`

	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`

	Name   string            `json:"name" yaml:"name"`
	Height uint64            `json:"height" yaml:"height"`
	Config map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
}

func ProposalUpgradeRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "upgrade",
		Handler:  postUpgradeProposalHandlerFn(cliCtx),
	}
}

func postUpgradeProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpgradeProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewUpgradeProposal(req.Title, req.Description, req.Name, req.Height, req.Config)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
