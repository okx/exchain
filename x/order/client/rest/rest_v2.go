package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/order/keeper"
	ordertype "github.com/okex/okchain/x/order/types"
)

// nolint
func RegisterRoutesV2(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/instruments/{instrument_id}/book", depthBookHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/order/placeorder", broadcastPlaceOrderRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/order/cancelorder", broadcastCancelOrderRequest(cliCtx)).Methods("POST")
}

func depthBookHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["instrument_id"]
		sizeStr := r.URL.Query().Get("size")

		// validate request
		if product == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		var size int
		var err error
		if sizeStr != "" {
			size, err = strconv.Atoi(sizeStr)
			if err != nil {
				common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
				return
			}
		}
		if size < 0 {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := keeper.NewQueryDepthBookParams(product, size)
		req, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/order/%s", ordertype.QueryDepthBookV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	Tx   auth.StdTx `json:"tx"`
	Mode string     `json:"mode"`
}

type placeCancelOrderResponse struct {
	types.TxResponse
	OrderID      string `json:"order_id"`
	ClientOid    string `json:"client_oid"`
	Result       bool   `json:"result"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// BroadcastTxRequest implements a tx broadcasting handler that is responsible
// for broadcasting a valid and signed tx to a full node. The tx can be
// broadcasted via a sync|async|block mechanism.
func broadcastPlaceOrderRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		orderID := ""
		events := res.Events
		if len(events) > 0 {
			attributes := events[0].Attributes
			for i := 0; i < len(attributes); i++ {
				if attributes[i].Key == "orderID" {
					orderID = attributes[i].Value
				}
			}
		}
		res2 := placeCancelOrderResponse{
			res,
			orderID,
			"",
			true,
			"",
			"",
		}
		if res.Code != 0 {
			res2.Result = false
			res2.ErrorCode = strconv.Itoa(int(res.Code))
			res2.ErrorMessage = res.Logs[0].Log

		}

		rest.PostProcessResponse(w, cliCtx, res2)
	}
}

func broadcastCancelOrderRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res2 := placeCancelOrderResponse{
			res,
			req.Tx.Msgs[0].(*ordertype.MsgCancelOrders).OrderIDs[0],
			"",
			true,
			"",
			"",
		}
		if res.Code != 0 {
			res2.Result = false
			res2.ErrorCode = strconv.Itoa(int(res.Code))
			res2.ErrorMessage = res.Logs[0].Log

		}

		rest.PostProcessResponse(w, cliCtx, res2)
	}
}
