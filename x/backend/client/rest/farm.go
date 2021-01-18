package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func registerFarmQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r = r.PathPrefix("/farm").Subrouter()
	r.HandleFunc("/pools/{whitelistOrNormal}", farmPoolsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/dashboard/{address}", farmDashboardHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/whitelist/max_apy", farmWhitelistMaxApyHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/pools/{poolName}/staked_info", farmStakedInfoHandler(cliCtx)).Methods("GET")
}
