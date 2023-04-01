package rest

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"

	"github.com/gorilla/mux"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"

	"github.com/okex/exchain/x/wasm/keeper"
	"github.com/okex/exchain/x/wasm/types"
)

func registerQueryRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/code", listCodesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/code/{codeID}", queryCodeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/code/{codeID}/contracts", listContractsByCodeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}", queryContractHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/state", queryContractStateAllHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/history", queryContractHistoryFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/smart/{query}", queryContractStateSmartHandlerFn(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/raw/{key}", queryContractStateRawHandlerFn(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/blocked_methods", queryContractBlockedMethodsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/params", queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/whitelist", queryContractWhitelistHandlerFn(cliCtx)).Methods("GET")
}

func queryParamsHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParams)

		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryContractWhitelistHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParams)

		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var params types.Params
		cliCtx.Codec.MustUnmarshalJSON(res, &params)
		var whitelist []string
		whitelist = strings.Split(params.CodeUploadAccess.Address, ",")
		// When params.CodeUploadAccess.Address == "", whitelist == []string{""} and len(whitelist) == 1.
		// Above case should be avoided.
		if len(whitelist) == 1 && whitelist[0] == "" {
			whitelist = []string{}
		}
		response := types.NewQueryAddressWhitelistResponse(whitelist)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, response)
	}
}

func listCodesHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		queryClient := types.NewQueryClient(cliCtx)
		pageReq, err := rest.ParseGRPCWasmPageRequest(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var reverse bool
		reverseStr := r.FormValue("reverse")
		if reverseStr == "" {
			reverse = false
		} else {
			reverse, err = strconv.ParseBool(reverseStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		res, err := queryClient.Codes(
			context.Background(),
			&types.QueryCodesRequest{
				Pagination: pageReq,
			},
		)

		if reverse {
			for i, j := 0, len(res.CodeInfos)-1; i < j; i, j = i+1, j-1 {
				res.CodeInfos[i], res.CodeInfos[j] = res.CodeInfos[j], res.CodeInfos[i]
			}
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)

	}
}

func queryCodeHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		codeId, err := strconv.ParseUint(mux.Vars(r)["codeID"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "codeId should be a number")
			return
		}
		queryClient := types.NewQueryClient(cliCtx)
		res, err := queryClient.Code(
			context.Background(),
			&types.QueryCodeRequest{
				CodeId: codeId,
			},
		)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if res == nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "contract not found")
			return
		}

		out, err := cliCtx.CodecProy.GetProtocMarshal().MarshalJSON(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}

func listContractsByCodeHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeID, err := strconv.ParseUint(mux.Vars(r)["codeID"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryClient := types.NewQueryClient(cliCtx)
		pageReq, err := rest.ParseGRPCWasmPageRequest(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var reverse bool
		reverseStr := r.FormValue("reverse")
		if reverseStr == "" {
			reverse = false
		} else {
			reverse, err = strconv.ParseBool(reverseStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		res, err := queryClient.ContractsByCode(
			context.Background(),
			&types.QueryContractsByCodeRequest{
				CodeId:     codeID,
				Pagination: pageReq,
			},
		)

		if reverse {
			for i, j := 0, len(res.Contracts)-1; i < j; i, j = i+1, j-1 {
				res.Contracts[i], res.Contracts[j] = res.Contracts[j], res.Contracts[i]
			}
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryContractHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryClient := types.NewQueryClient(cliCtx)
		res, err := queryClient.ContractInfo(
			context.Background(),
			&types.QueryContractInfoRequest{
				Address: mux.Vars(r)["contractAddr"],
			},
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		result := fromGrpcContractInfo(res)

		if result == nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "contract not found")
			return
		}

		out, err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("json marshal error %v", err))
			return
		}

		rest.PostProcessResponse(w, cliCtx, out)
	}
}

func queryContractBlockedMethodsHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr, err := sdk.AccAddressFromBech32(mux.Vars(r)["contractAddr"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, keeper.QueryListContractBlockedMethod, addr.String())
		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(res))
	}
}

func queryContractStateAllHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["contractAddr"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryClient := types.NewQueryClient(cliCtx)
		pageReq, err := rest.ParseGRPCWasmPageRequest(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var reverse bool
		reverseStr := r.FormValue("reverse")
		if reverseStr == "" {
			reverse = false
		} else {
			reverse, err = strconv.ParseBool(reverseStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		res, err := queryClient.AllContractState(
			context.Background(),
			&types.QueryAllContractStateRequest{
				Address:    addr,
				Pagination: pageReq,
			},
		)

		if reverse {
			for i, j := 0, len(res.Models)-1; i < j; i, j = i+1, j-1 {
				res.Models[i], res.Models[j] = res.Models[j], res.Models[i]
			}
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryContractStateRawHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := newArgDecoder(hex.DecodeString)
		decoder.encoding = mux.Vars(r)["encoding"]
		queryData, err := decoder.DecodeString(mux.Vars(r)["key"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		queryClient := types.NewQueryClient(cliCtx)
		res, err := queryClient.RawContractState(
			context.Background(),
			&types.QueryRawContractStateRequest{
				Address:   mux.Vars(r)["contractAddr"],
				QueryData: queryData,
			},
		)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		out, err := cliCtx.CodecProy.GetProtocMarshal().MarshalJSON(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}

func queryContractStateSmartHandlerFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		decoder := newArgDecoder(hex.DecodeString)
		decoder.encoding = mux.Vars(r)["encoding"]
		queryData, err := decoder.DecodeString(mux.Vars(r)["query"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		queryClient := types.NewQueryClient(cliCtx)
		res, err := queryClient.SmartContractState(
			context.Background(),
			&types.QuerySmartContractStateRequest{
				Address:   mux.Vars(r)["contractAddr"],
				QueryData: queryData,
			},
		)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		out, err := cliCtx.CodecProy.GetProtocMarshal().MarshalJSON(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, json.RawMessage(out))
	}
}

func queryContractHistoryFn(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["contractAddr"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryClient := types.NewQueryClient(cliCtx)
		pageReq, err := rest.ParseGRPCWasmPageRequest(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var reverse bool
		reverseStr := r.FormValue("reverse")
		if reverseStr == "" {
			reverse = false
		} else {
			reverse, err = strconv.ParseBool(reverseStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		res, err := queryClient.ContractHistory(
			context.Background(),
			&types.QueryContractHistoryRequest{
				Address:    addr,
				Pagination: pageReq,
			},
		)

		if reverse {
			for i, j := 0, len(res.Entries)-1; i < j; i, j = i+1, j-1 {
				res.Entries[i], res.Entries[j] = res.Entries[j], res.Entries[i]
			}
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

type argumentDecoder struct {
	// dec is the default decoder
	dec      func(string) ([]byte, error)
	encoding string
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	switch a.encoding {
	case "hex":
		return hex.DecodeString(s)
	case "base64":
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}
