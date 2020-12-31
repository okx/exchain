package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/okex/okexchain/x/distribution/client/common"
	"github.com/okex/okexchain/x/distribution/types"
	comm "github.com/okex/okexchain/x/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Get the rewards withdrawal address
	r.HandleFunc(
		"/distribution/delegators/{delegatorAddr}/withdraw_address",
		delegatorWithdrawalAddrHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// accumulated commission of a single validator
	r.HandleFunc(
		"/distribution/validators/{validatorAddr}/validator_commission",
		accumulatedCommissionHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get the current distribution parameter values
	r.HandleFunc(
		"/distribution/parameters",
		paramsHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get the amount held in the community pool
	r.HandleFunc(
		"/distribution/community_pool",
		communityPoolHandler(cliCtx, queryRoute),
	).Methods("GET")
}

// HTTP request handler to query a delegation rewards
func delegatorWithdrawalAddrHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delegatorAddr, ok := checkDelegatorAddressVar(w, r)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz := cliCtx.Codec.MustMarshalJSON(types.NewQueryDelegatorWithdrawAddrParams(delegatorAddr))
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/withdraw_addr", queryRoute), bz)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the distribution params values
func paramsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params, err := common.QueryParams(cliCtx, queryRoute)
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, types.CodeInvalidRoute, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, params)
	}
}

// HTTP request handler to query the community pool coins
func communityPoolHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/community_pool", queryRoute), nil)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		var result sdk.SysCoins
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// HTTP request handler to query the accumulated commission of one single validator
func accumulatedCommissionHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatorAddr, ok := checkValidatorAddressVar(w, r)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bin := cliCtx.Codec.MustMarshalJSON(types.NewQueryValidatorCommissionParams(validatorAddr))
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/validator_commission", queryRoute), bin)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
