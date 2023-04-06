package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"net/http"
)

func registerQueryRoutesV2(cliCtx clientCtx.CLIContext, r *mux.Router) {
	r.HandleFunc("/wasm/code/{codeID}", queryCodeHandlerFnV2(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}", queryContractHandlerFnV2(cliCtx)).Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/raw/{key}", queryContractStateRawHandlerFnV2(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
	r.HandleFunc("/wasm/contract/{contractAddr}/smart/{query}", queryContractStateSmartHandlerFnV2(cliCtx)).Queries("encoding", "{encoding}").Methods("GET")
}

func queryCodeHandlerFnV2(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryCodeHandler(cliCtx, w, r)
		if err != nil {
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

func queryContractStateRawHandlerFnV2(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractStateRawHandler(cliCtx, w, r)
		if err != nil {
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

func queryContractHandlerFnV2(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractHandler(cliCtx, w, r)
		if err != nil {
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

func queryContractStateSmartHandlerFnV2(cliCtx clientCtx.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queryContractStateSmartHandler(cliCtx, w, r)
		if err != nil {
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
