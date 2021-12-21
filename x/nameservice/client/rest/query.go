/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 5:47 上午
# @File : query.go
# @Description :
# @Attention :
*/
package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/x/nameservice/types"
	"net/http"
)

// RegisterRoutes registers nameservice-related REST handlers to a router
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// this line is used by starport scaffolding
	r.HandleFunc("/nameservice/whois", listWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/{key}", getWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/{key}/resolve", resolveNameHandler(cliCtx, "nameservice")).Methods("GET")
}

func listWhoisHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryListWhois), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getWhoisHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", storeName, types.QueryGetWhois, key), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolveNameHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars["key"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", storeName, types.QueryResolveName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
