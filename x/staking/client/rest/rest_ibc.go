package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/x/staking/types"
)

func RegisterOriginRPCRoutersForGRPC(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cosmos/staking/v1beta1/delegators/{delegatorAddr}/unbonding_delegations",
		delegatorUnbondingDelegationsHandlerFn2(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cosmos/staking/v1beta1/delegations/{delegatorAddr}",
		delegatorDelegationsHandlerFn(cliCtx),
	).Methods("GET")

}

func delegatorDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBonds(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorDelegations))
}
