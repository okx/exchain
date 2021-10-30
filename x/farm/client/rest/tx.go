package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/dependence/cosmos-sdk/client/context"
	// sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	// "github.com/okex/exchain/dependence/cosmos-sdk/types/rest"
	// "github.com/okex/exchain/dependence/cosmos-sdk/x/auth/client/utils"
	// "github.com/okex/exchain/x/farm/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// r.HandleFunc(
	// TODO: Define the Rest route ,
	// Call the function which should be executed for this route),
	// ).Methods("POST")
}

/*
// Action TX body
type <Action>Req struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
}

func <Action>RequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req <Action>Req
		vars := mux.Vars(r)

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// TODO: Define the module tx logic for this action

		utils.WriteGenerateStdTxResponse(w, cliCtx, BaseReq, []sdk.Msg{msg})
	}
}
*/
