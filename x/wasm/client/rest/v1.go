package rest

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/x/wasm/types"
	"net/http"
)

func registerQueryRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/code", listCodesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/code/{codeID}", queryCodeHandlerFnV1(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/code/{codeID}/contracts", listContractsByCodeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}", queryContractHandlerFnV1(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/state", queryContractStateAllHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/history", queryContractHistoryFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/smart/{query}", queryContractStateSmartHandlerFnV1(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/raw/{key}", queryContractStateRawHandlerFnV1(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/blocked_methods", queryContractBlockedMethodsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/params", queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/whitelist", queryContractWhitelistHandlerFn(cliCtx)).Methods("GET")
}

type codeInfo struct {
	Id                    uint64 `json:"id"`
	Creator               string `json:"creator"`
	DataHash              string `json:"data_hash"`
	InstantiatePermission struct {
		Permission string `json:"permission"`
	} `json:"instantiate_permission"`
	Data []byte `json:"data"`
}

func fromGrpcCodeInfo(codeResp *types.QueryCodeResponse) *codeInfo {
	if codeResp == nil {
		return nil
	}
	var ret codeInfo
	ret.Id = codeResp.CodeID
	ret.Creator = codeResp.Creator
	ret.DataHash = codeResp.DataHash.String()
	ret.InstantiatePermission.Permission = codeResp.InstantiatePermission.Permission.String()
	ret.Data = codeResp.Data

	return &ret
}

func queryCodeHandlerFnV1(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryCodeHandler(cliCtx, w, r)
		if err != nil {
			return
		}

		result := fromGrpcCodeInfo(res)
		out, err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}

func queryContractStateRawHandlerFnV1(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractStateRawHandler(cliCtx, w, r)
		if err != nil {
			return
		}

		rest.PostProcessResponse(w, cliCtx, base64.StdEncoding.EncodeToString(res.Data))
	}
}

type contractInfo struct {
	Address string `json:"address"`
	CodeId  uint64 `json:"code_id"`
	Creator string `json:"creator"`
	Admin   string `json:"admin"`
	Label   string `json:"label"`
}

func fromGrpcContractInfo(contractResp *types.QueryContractInfoResponse) *contractInfo {
	if contractResp == nil {
		return nil
	}
	var ret contractInfo
	ret.Admin = contractResp.Admin
	ret.Label = contractResp.Label
	ret.CodeId = contractResp.CodeID
	ret.Address = contractResp.Address
	ret.Creator = contractResp.Creator

	return &ret
}

func queryContractHandlerFnV1(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractHandler(cliCtx, w, r)
		if err != nil {
			return
		}

		result := fromGrpcContractInfo(res)

		out, err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}

type smartInfo struct {
	Smart string `json:"smart"`
}

func fromGrpcSmartInfo(smartResp *types.QuerySmartContractStateResponse) *smartInfo {
	if smartResp == nil {
		return &smartInfo{}
	}
	var ret smartInfo
	ret.Smart = base64.StdEncoding.EncodeToString(smartResp.Data)

	return &ret
}

func queryContractStateSmartHandlerFnV1(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractStateSmartHandler(cliCtx, w, r)
		if err != nil {
			return
		}

		result := fromGrpcSmartInfo(res)
		out, err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}
