package rest

import (
	"net/http"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	comm "github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/gov"
	govrest "github.com/okx/okbchain/x/gov/client/rest"
)

type ManageSysContractAddressProposalReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	ContractAddr sdk.AccAddress `json:"contract_address" yaml:"contract_address"`
	IsAdded      bool           `json:"is_added" yaml:"is_added"`

	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit  sdk.SysCoins   `json:"deposit" yaml:"deposit"`
}

// ManageSysContractAddressProposalRESTHandler defines evm proposal handler
func ManageSysContractAddressProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "manage_system_contract_address",
		Handler:  postManageSysContractAddressProposalHandlerFn(cliCtx),
	}
}

func postManageSysContractAddressProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ManageSysContractAddressProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewManageSysContractAddressProposal(
			req.Title,
			req.Description,
			req.ContractAddr,
			req.IsAdded,
		)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
