package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	comm "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/feesplit/types"
	govRest "github.com/okex/exchain/x/gov/client/rest"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/feesplit/contract/{contract}", withdrawAddrHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/withdrawer-contracts/{addr}", withdrawerContractsHandlerFn(cliCtx)).Methods("GET")
}

// FeeSplitSharesProposalRESTHandler defines feesplit proposal handler
func FeeSplitSharesProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}

func withdrawAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contract := mux.Vars(r)["contract"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		req := &types.QueryFeeSplitRequest{ContractAddress: contract}
		data, err := cliCtx.Codec.MarshalJSON(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryFeeSplit), data)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		var result types.QueryFeeSplitResponse
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func withdrawerContractsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["addr"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		req := &types.QueryWithdrawerFeeSplitsRequest{
			WithdrawerAddress: addr,
			Pagination:        nil,
		}
		data, err := cliCtx.Codec.MarshalJSON(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			types.RouterKey, types.QueryWithdrawerFeeSplits), data)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		var result types.QueryWithdrawerFeeSplitsResponse
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
