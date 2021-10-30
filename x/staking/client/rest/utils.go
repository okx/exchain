package rest

import (
	"fmt"
	"net/http"

	"github.com/okex/exchain/x/common"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/x/staking/types"
)

// contains checks if the a given query contains one of the tx types
func contains(stringSlice []string, txType string) bool {
	for _, word := range stringSlice {
		if word == txType {
			return true
		}
	}
	return false
}

// queries staking txs
func queryTxs(cliCtx context.CLIContext, action string, delegatorAddr string) (*sdk.SearchTxsResult, error) {
	page := 1
	limit := 100
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, action),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, delegatorAddr),
	}

	return utils.QueryTxsByEvents(cliCtx, events, page, limit)
}

func queryDelegator(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32DelAddr := mux.Vars(r)["delegatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelAddr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, types.CodeNoDelegatorExisted, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryDelegatorParams(delegatorAddr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, uint32(common.ErrorABCIQueryFails), err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryValidator(cliCtx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32ValAddr := mux.Vars(r)["validatorAddr"]

		validatorAddr, err := sdk.ValAddressFromBech32(bech32ValAddr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, types.CodeBadValidatorAddr, "validator address is invalid")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryValidatorParams(validatorAddr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
