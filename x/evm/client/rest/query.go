package rest

import (
	"fmt"
	"net/http"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	clientCtx "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
	comm "github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/evm/client/utils"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

func registerQueryRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	r.HandleFunc("/evm/system/contract/address", QueryManageSysContractAddressFn(cliCtx)).Methods("GET")
}

// QueryManageSysContractAddressFn defines evm contract method blocked list handler
func QueryManageSysContractAddressFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := fmt.Sprintf("custom/%s/%s", evmtypes.ModuleName, evmtypes.QuerySysContractAddress)
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		addr, _, err := cliCtx.QueryWithData(path, nil)
		if err != nil {
			sdkErr := comm.ParseSDKError(err.Error())
			comm.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		ethAddr := ethcommon.BytesToAddress(addr).Hex()
		result := utils.ResponseSysContractAddress{Address: ethAddr}

		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}
