package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
	"github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/staking/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// query delegator info
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}",
		delegatorHandlerFn(cliCtx),
	).Methods("GET")

	// query delegator's unbonding delegation
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/unbonding_delegations",
		delegatorUnbondingDelegationsHandlerFn(cliCtx),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/validators",
		delegatorValidatorsHandlerFn(cliCtx),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/validators/{validatorAddr}",
		delegatorValidatorHandlerFn(cliCtx),
	).Methods("GET")

	// Get all delegations to a validator
	r.HandleFunc(
		"/staking/validators/{validatorAddr}/delegations",
		validatorDelegationsHandlerFn(cliCtx),
	).Methods("GET")

	// Queries delegate info for given validator delegator pair
	r.HandleFunc(
		"/staking/validators/{validatorAddr}/delegations/{delegatorAddr}",
		delegationHandlerFn(cliCtx),
	).Methods("GET")

	// query the all shares on a validator
	r.HandleFunc(
		"/staking/validators/{validatorAddr}/shares",
		validatorAllSharesHandlerFn(cliCtx),
	).Methods("GET")

	// get all validators
	r.HandleFunc(
		"/staking/validators",
		validatorsHandlerFn(cliCtx),
	).Methods("GET")

	// get a single validator info
	r.HandleFunc(
		"/staking/validators/{validatorAddr}",
		validatorHandlerFn(cliCtx),
	).Methods("GET")

	// Get HistoricalInfo at a given height
	r.HandleFunc(
		"/staking/historical_info/{height}",
		historicalInfoHandlerFn(cliCtx),
	).Methods("GET")

	// get the current state of the staking pool
	r.HandleFunc(
		"/staking/pool",
		poolHandlerFn(cliCtx),
	).Methods("GET")

	// get the current staking parameter values
	r.HandleFunc(
		"/staking/parameters",
		paramsHandlerFn(cliCtx),
	).Methods("GET")

	// get the current staking address values
	r.HandleFunc(
		"/staking/address",
		addressHandlerFn(cliCtx),
	).Methods("GET")

	// get the current staking address values
	r.HandleFunc(
		"/staking/address/{validatorAddr}/validator_address",
		validatorAddressHandlerFn(cliCtx),
	).Methods("GET")

	// get the current staking address values
	r.HandleFunc(
		"/staking/address/{validatorAddr}/account_address",
		accountAddressHandlerFn(cliCtx),
	).Methods("GET")

}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDelegator(cliCtx, "custom/staking/delegatorValidators")
}

// HTTP request handler to query a delegation
func delegationHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBondsInfo(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorDelegator))
}

// HTTP request handler to query the proxy relationship on a proxy delegator
func delegatorProxyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProxy))
}

// HTTP request handler to get information from a currently bonded validator
func delegatorValidatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBondsInfo(cliCtx, "custom/staking/delegatorValidator")
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryValidator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorDelegations))
}

// HTTP request handler to query the info of delegator's unbonding delegation
func delegatorUnbondingDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryUnbondingDelegation))
}

func delegatorUnbondingDelegationsHandlerFn2(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBonds(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryUnbondingDelegation2))
}

// HTTP request handler to query the info of a delegator
func delegatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegator))
}

// HTTP request handler to query the all shares added to a validator
func validatorAllSharesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryValidator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorAllShares))
}

// HTTP request handler to query historical info at a given height
func historicalInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx := cliCtx
		vars := mux.Vars(r)
		heightStr := vars["height"]
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil || height < 0 {
			common.HandleErrorMsg(w, cliCtx, common.CodeInternalError, fmt.Sprintf("Must provide non-negative integer for height: %v", err))
			return
		}

		params := types.NewQueryHistoricalInfoParams(height)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHistoricalInfo)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query list of validators
func validatorsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeArgsWithLimit, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		status := r.FormValue("status")
		if status == "" {
			status = sdk.BondStatusBonded
		}

		params := types.NewQueryValidatorsParams(page, limit, status)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidators)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryValidator(cliCtx, "custom/staking/validator")
}

// HTTP request handler to query the pool information
func poolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/staking/pool", nil)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the staking params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/staking/parameters", nil)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func addressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/staking/address", nil)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query one validator address
func validatorAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryValidatorAddr(cliCtx, "custom/staking/validatorAddress")
}

// HTTP request handler to query one validator's account address
func accountAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryValidatorAddr(cliCtx, "custom/staking/validatorAccAddress")
}

func queryValidatorAddr(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		validatorAddr := mux.Vars(r)["validatorAddr"]

		res, height, err := cliCtx.QueryWithData(endpoint, []byte(validatorAddr))
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
