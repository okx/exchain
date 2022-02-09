package keeper

import (
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/feemarket/types"
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
	Params types.Params `json:"param"`
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

// Params implements the Query/Params gRPC method
func (k Keeper) Params(ctx sdk.Context, _ *QueryParamsRequest) (*QueryParamsResponse, error) {
	params := k.GetParams(ctx)

	return &QueryParamsResponse{
		Params: params,
	}, nil
}

// BaseFee implements the Query/BaseFee gRPC method
func (k Keeper) BaseFee(ctx sdk.Context, _ *QueryBaseFeeRequest) (*QueryBaseFeeResponse, error) {
	res := &QueryBaseFeeResponse{}
	baseFee := k.GetBaseFee(ctx)

	if baseFee != nil {
		aux := sdk.NewIntFromBigInt(baseFee)
		res.BaseFee = aux.BigInt()
	}

	return res, nil
}

// BlockGas implements the Query/BlockGas gRPC method
func (k Keeper) BlockGas(ctx sdk.Context, _ *QueryBlockGasRequest) (*QueryBlockGasResponse, error) {
	gas := k.GetBlockGasUsed(ctx)

	return &QueryBlockGasResponse{
		Gas: int64(gas),
	}, nil
}
