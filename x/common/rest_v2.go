package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/tendermint/tendermint/libs/log"
)

// const
const (
	// common error
	ErrorMissingRequiredParam errorCodeV2 = 60001
	ErrorInvalidParam         errorCodeV2 = 60002
	ErrorServerException      errorCodeV2 = 60003
	ErrorDataNotExist         errorCodeV2 = 60004
	ErrorCodecFails           errorCodeV2 = 60005
	ErrorABCIQueryFails       errorCodeV2 = 60006
	ErrorArgsWithLimit        errorCodeV2 = 60007

	// account error
	ErrorInvalidAddress errorCodeV2 = 61001

	// order error
	ErrorOrderNotExist        errorCodeV2 = 62001
	ErrorInvalidCurrency      errorCodeV2 = 62002
	ErrorEmptyInstrumentID    errorCodeV2 = 62003
	ErrorInstrumentIDNotExist errorCodeV2 = 62004

	// staking error
	ErrorInvalidValidatorAddress errorCodeV2 = 63001
	ErrorInvalidDelegatorAddress errorCodeV2 = 63002

	// farm error
	ErrorInvalidAccountAddress errorCodeV2 = 64001
)

func defaultErrorMessageV2(code errorCodeV2) (message string) {
	switch code {
	case ErrorMissingRequiredParam:
		message = "missing required param"
	case ErrorInvalidParam:
		message = "invalid request param"
	case ErrorServerException:
		message = "internal server error"
	case ErrorDataNotExist:
		message = "data not exists"
	case ErrorCodecFails:
		message = "inner CODEC failed"
	case ErrorABCIQueryFails:
		message = "abci query failed"
	case ErrorArgsWithLimit:
		message = "failed to parse args with limit"

	case ErrorInvalidAddress:
		message = "invalid address"

	case ErrorOrderNotExist:
		message = "order not exists"
	case ErrorInvalidCurrency:
		message = "invalid currency"
	case ErrorEmptyInstrumentID:
		message = "instrument_id is empty"
	case ErrorInstrumentIDNotExist:
		message = "instrument_id not exists"

	// staking
	case ErrorInvalidValidatorAddress:
		message = "invalid validator address"
	case ErrorInvalidDelegatorAddress:
		message = "invalid delegator address"

	// farm
	case ErrorInvalidAccountAddress:
		message = "invalid account address"

	default:
		message = "unknown error"
	}
	return
}

type errorCodeV2 int

func (code errorCodeV2) Code() string {
	return strconv.Itoa(int(code))
}

func (code errorCodeV2) Message() string {
	return defaultErrorMessageV2(code)
}

type responseErrorV2 struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HandleErrorResponseV2 is the handler of error response with V2 standard
func HandleErrorResponseV2(w http.ResponseWriter, statusCode int, errCode errorCodeV2) {
	response, err := json.Marshal(responseErrorV2{
		Code:    errCode.Code(),
		Message: errCode.Message(),
	})
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if _, err = w.Write(response); err != nil {
			log.NewTMLogger(os.Stdout).Debug(fmt.Sprintf("error: %v", err.Error()))
		}
	}
}

// HandleSuccessResponseV2 is the handler of successful response with V2 standard
func HandleSuccessResponseV2(w http.ResponseWriter, data []byte) {
	logger := log.NewTMLogger(os.Stdout)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		logger.Debug(fmt.Sprintf("error: %v", err.Error()))
	}
}

// HandleResponseV2 handles the response of V2 standard
func HandleResponseV2(w http.ResponseWriter, data []byte, err error) {
	if err != nil {
		HandleErrorResponseV2(w, http.StatusInternalServerError, ErrorServerException)
		return
	}
	if len(data) == 0 {
		HandleErrorResponseV2(w, http.StatusBadRequest, ErrorDataNotExist)
	}

	HandleSuccessResponseV2(w, data)
}

// JSONMarshalV2 marshals info into JSON based on V2 standard
func JSONMarshalV2(v interface{}) ([]byte, error) {
	var jsonV2 = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "v2",
	}.Froze()

	return jsonV2.MarshalIndent(v, "", "  ")
}

// JSONUnmarshalV2 unmarshals JSON bytes based on V2 standard
func JSONUnmarshalV2(data []byte, v interface{}) error {
	var jsonV2 = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "v2",
	}.Froze()

	return jsonV2.Unmarshal(data, v)
}
