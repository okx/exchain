package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	comm "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/distribution/client/common"
	"github.com/okex/exchain/x/distribution/types"
)

func cm45AccumulatedCommissionHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatorAddr := mux.Vars(r)["validatorAddr"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bin := cliCtx.Codec.MustMarshalJSON(types.NewQueryValidatorCommissionRequest(validatorAddr))
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/commission", queryRoute), bin)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}
		var commission types.QueryValidatorCommissionResponse
		cliCtx.Codec.MustUnmarshalJSON(res, &commission)
		wrappedCommission := types.NewWrappedCommission(commission)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, wrappedCommission)
	}
}

func cm45DelegatorRewardsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		delegatorAddr, ok := checkDelegatorAddressVar(w, r)
		if !ok {
			return
		}

		params := types.NewQueryDelegatorParams(delegatorAddr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeMarshalJSONFailed, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorTotalRewards)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			if sdkErr.Code == types.CodeEmptyDelegationDistInfo {
				total := sdk.DecCoins{}
				delRewards := make([]types.DelegationDelegatorReward, 0)
				totalRewards := types.NewQueryDelegatorTotalRewardsResponse(delRewards, total)
				cliCtx = cliCtx.WithHeight(height)
				rest.PostProcessResponse(w, cliCtx, totalRewards)
			} else {
				comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			}
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func cm45ParamsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
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
		wrappedParams := types.NewWrappedParams(params)
		rest.PostProcessResponse(w, cliCtx, wrappedParams)
	}
}
