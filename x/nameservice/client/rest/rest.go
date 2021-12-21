/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 5:47 上午
# @File : rest.go
# @Description :
# @Attention :
*/
package rest

import (
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

// RegisterRoutes registers ammswap-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
