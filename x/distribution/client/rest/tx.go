package rest

import (
	"fmt"
	"net/http"

	comm "github.com/okex/exchain/x/common"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/okex/exchain/x/distribution/client/common"
	"github.com/okex/exchain/x/distribution/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, _ string) {
	// Replace the rewards withdrawal address
	r.HandleFunc(
		"/distribution/delegators/{delegatorAddr}/withdraw_address",
		setDelegatorWithdrawalAddrHandlerFn(cliCtx),
	).Methods("POST")

	// Withdraw validator rewards and commission
	r.HandleFunc(
		"/distribution/validators/{validatorAddr}/rewards",
		withdrawValidatorRewardsHandlerFn(cliCtx),
	).Methods("POST")

}

type (
	withdrawRewardsReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}

	setWithdrawalAddrReq struct {
		BaseReq         rest.BaseReq   `json:"base_req" yaml:"base_req"`
		WithdrawAddress sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
	}
)

// Replace the rewards withdrawal address
func setDelegatorWithdrawalAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setWithdrawalAddrReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// read and validate URL's variables
		delAddr, ok := checkDelegatorAddressVar(w, r)
		if !ok {
			return
		}

		msg := types.NewMsgSetWithdrawAddress(delAddr, req.WithdrawAddress)
		if err := msg.ValidateBasic(); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// Withdraw validator rewards and commission
func withdrawValidatorRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// read and validate URL's variable
		valAddr, ok := checkValidatorAddressVar(w, r)
		if !ok {
			return
		}

		// prepare multi-message transaction
		msgs, err := common.WithdrawValidatorRewardsAndCommission(valAddr)
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeInvalidParam, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, msgs)
	}
}

// Auxiliary

func checkDelegatorAddressVar(w http.ResponseWriter, r *http.Request) (sdk.AccAddress, bool) {
	addr, err := sdk.AccAddressFromBech32(mux.Vars(r)["delegatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, int(types.CodeAccAddressFromBech32Failed), fmt.Sprintf("invalid address：%s", mux.Vars(r)["delegatorAddr"]))
		return nil, false
	}

	return addr, true
}

func checkValidatorAddressVar(w http.ResponseWriter, r *http.Request) (sdk.ValAddress, bool) {
	addr, err := sdk.ValAddressFromBech32(mux.Vars(r)["validatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, int(types.CodeValAddressFromBech32), fmt.Sprintf("invalid address：%s", mux.Vars(r)["validatorAddr"]))
		return nil, false
	}

	return addr, true
}
