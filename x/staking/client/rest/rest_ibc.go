package rest

import (
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

func RegisterOriginRPCRoutersForGRPC(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cosmos/staking/v1beta1/delegators/{delegatorAddr}/unbonding_delegations",
		delegatorUnbondingDelegationsHandlerFn(cliCtx),
	).Methods("GET")

}
