package types

import (
	"math/big"
)

// Supported endpoints
const (
	// QueryParameters defines 	QueryParameters = "params" query route path
	QueryParameters = "params"
	QueryBaseFee    = "basefee"
	QueryBlockGas   = "blockgas"
)

// QueryParamsResponse is response type for Param query
type QueryParamsRequest struct {
}

// QueryParamsResponse is response type for Param query
type QueryParamsResponse struct {
	Params Params `json:"param"`
}

// QueryParamsResponse is response type for Param query
type QueryBaseFeeRequest struct {
}

// QueryParamsResponse is response type for Param query
type QueryBaseFeeResponse struct {
	BaseFee *big.Int `json:"baseFee"`
}

// QueryParamsResponse is response type for Param query
type QueryBlockGasRequest struct {
}

// QueryParamsResponse is response type for Param query
type QueryBlockGasResponse struct {
	Gas int64 `json:"gas"`
}
