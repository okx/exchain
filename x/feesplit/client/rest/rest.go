package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/query"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
	comm "github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/feesplit/types"
	govRest "github.com/okx/okbchain/x/gov/client/rest"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/feesplit/contract/{contract}", contractHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/deployer/{deployer}", deployerHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/withdrawer/{withdrawer}", withdrawerHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func RegisterRoutesV2(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/feesplit/contract/{contract}", contractHandlerFnV2(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/deployer/{deployer}", deployerHandlerFnV2(cliCtx)).Methods("GET")
	r.HandleFunc("/feesplit/parameters", queryParamsHandlerFnV2(cliCtx)).Methods("GET")
}

// FeeSplitSharesProposalRESTHandler defines feesplit proposal handler
func FeeSplitSharesProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			types.RouterKey, types.QueryParameters), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func contractHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

func deployerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["deployer"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req := &types.QueryDeployerFeeSplitsRequest{
			DeployerAddress: addr,
			Pagination:      query.NewPaginateFromPageLimit(page, limit),
		}
		data, err := cliCtx.Codec.MarshalJSON(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			types.RouterKey, types.QueryDeployerFeeSplits), data)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		var result types.QueryDeployerFeeSplitsResponse
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func withdrawerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["withdrawer"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req := &types.QueryWithdrawerFeeSplitsRequest{
			WithdrawerAddress: addr,
			Pagination:        query.NewPaginateFromPageLimit(page, limit),
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

func contractHandlerFnV2(cliCtx context.CLIContext) http.HandlerFunc {
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
		// convert ex to 0x
		if addr, err := sdk.AccAddressFromBech32(result.FeeSplit.DeployerAddress); err == nil {
			result.FeeSplit.DeployerAddress = ethcmn.BytesToAddress(addr.Bytes()).String()
		}
		if addr, err := sdk.AccAddressFromBech32(result.FeeSplit.WithdrawerAddress); err == nil {
			result.FeeSplit.WithdrawerAddress = ethcmn.BytesToAddress(addr.Bytes()).String()
		}

		resultJson, err := json.Marshal(comm.GetBaseResponse(result.FeeSplit))
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resultJson)
	}
}

func deployerHandlerFnV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["deployer"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req := &types.QueryDeployerFeeSplitsRequest{
			DeployerAddress: addr,
			Pagination:      query.NewPaginateFromPageLimit(page, limit),
		}
		data, err := cliCtx.Codec.MarshalJSON(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to marshal query params: %s", err))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			types.RouterKey, types.QueryDeployerFeeSplitsDetail), data)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		var result types.QueryDeployerFeeSplitsResponseV2
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}
		// convert ex to 0x
		for i, fs := range result.FeeSplits {
			if addr, err := sdk.AccAddressFromBech32(fs.DeployerAddress); err == nil {
				result.FeeSplits[i].DeployerAddress = ethcmn.BytesToAddress(addr.Bytes()).String()
			}
			if addr, err := sdk.AccAddressFromBech32(fs.WithdrawerAddress); err == nil {
				result.FeeSplits[i].WithdrawerAddress = ethcmn.BytesToAddress(addr.Bytes()).String()
			}
		}

		resultJson, err := json.Marshal(comm.GetBaseResponse(result))
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resultJson)
	}
}

func queryParamsHandlerFnV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			types.RouterKey, types.QueryParameters), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result types.QueryParamsResponse
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeUnMarshalJSONFailed, err.Error())
			return
		}

		resultJson, err := json.Marshal(comm.GetBaseResponse(result))
		if err != nil {
			comm.HandleErrorMsg(w, cliCtx, comm.CodeMarshalJSONFailed, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resultJson)
	}
}
