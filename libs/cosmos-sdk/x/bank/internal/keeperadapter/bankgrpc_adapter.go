package keeperadapter

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/typesadapter"
)

type BankQueryServer struct {
	bankKeeper   BankKeeperAdapter
	supplyKeeper SupplyKeeper
}

func NewBankQueryServer(bankKeeper BankKeeperAdapter, supplyKeeper SupplyKeeper) *BankQueryServer {
	return &BankQueryServer{bankKeeper: bankKeeper, supplyKeeper: supplyKeeper}
}

var _ typesadapter.QueryServer = &BankQueryServer{}

// Balance implements the Query/Balance gRPC method
func (k BankQueryServer) Balance(ctx context.Context, req *typesadapter.QueryBalanceRequest) (*typesadapter.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	balance := k.bankKeeper.GetBalance(sdkCtx, address, req.Denom)
	dapter := sdk.CoinToCoinAdapter(balance)
	return &typesadapter.QueryBalanceResponse{Balance: &dapter}, nil
}

// AllBalances implements the Query/AllBalances gRPC method
func (k BankQueryServer) AllBalances(ctx context.Context, req *typesadapter.QueryAllBalancesRequest) (*typesadapter.QueryAllBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	balances := k.bankKeeper.GetAllBalances(sdkCtx, addr)
	adapters := sdk.CoinsToCoinAdapters(balances)

	pageRes := &query.PageResponse{NextKey: nil, Total: uint64(len(balances))}

	return &typesadapter.QueryAllBalancesResponse{Balances: adapters, Pagination: pageRes}, nil
}

// TotalSupply implements the Query/TotalSupply gRPC method
func (k BankQueryServer) TotalSupply(ctx context.Context, req *typesadapter.QueryTotalSupplyRequest) (*typesadapter.QueryTotalSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	supply := k.supplyKeeper.GetSupply(sdkCtx)
	total := supply.GetTotal()
	adapters := sdk.CoinsToCoinAdapters(total)
	pageRes := &query.PageResponse{NextKey: nil, Total: uint64(len(adapters))}
	return &typesadapter.QueryTotalSupplyResponse{Supply: adapters, Pagination: pageRes}, nil
}

// SupplyOf implements the Query/SupplyOf gRPC method
func (k BankQueryServer) SupplyOf(c context.Context, req *typesadapter.QuerySupplyOfRequest) (*typesadapter.QuerySupplyOfResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	sdkCtx := sdk.UnwrapSDKContext(c)
	supply := k.supplyKeeper.GetSupply(sdkCtx)
	total := supply.GetTotal()
	coin := sdk.Coin{
		req.Denom,
		total.AmountOf(req.Denom),
	}
	adapter := sdk.CoinToCoinAdapter(coin)
	return &typesadapter.QuerySupplyOfResponse{Amount: adapter}, nil
}

// Params implements the gRPC service handler for querying x/bank parameters.
func (k BankQueryServer) Params(ctx context.Context, req *typesadapter.QueryParamsRequest) (*typesadapter.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	//TODO params is part adapter
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sendEnable := k.bankKeeper.GetSendEnabled(sdkCtx)
	adapter := typesadapter.Params{
		SendEnabled:        nil, // maybe need init
		DefaultSendEnabled: sendEnable,
	}
	return &typesadapter.QueryParamsResponse{Params: adapter}, nil
}

// DenomsMetadata implements Query/DenomsMetadata gRPC method.
func (k BankQueryServer) DenomsMetadata(c context.Context, req *typesadapter.QueryDenomsMetadataRequest) (*typesadapter.QueryDenomsMetadataResponse, error) {
	return nil, types.ErrUnSupportQueryType("Query/DenomsMetadata")
}

// DenomMetadata implements Query/DenomMetadata gRPC method.
func (k BankQueryServer) DenomMetadata(c context.Context, req *typesadapter.QueryDenomMetadataRequest) (*typesadapter.QueryDenomMetadataResponse, error) {
	return nil, types.ErrUnSupportQueryType("Query/DenomMetadata")
}
