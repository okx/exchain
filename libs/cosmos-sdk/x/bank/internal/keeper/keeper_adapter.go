package keeper
//
//import (
//	"context"
//	"github.com/cosmos/cosmos-sdk/x/bank/types"
//	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
//	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
//	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
//	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/typesadapter"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)
//
//var ()
//
//var _ typesadapter.QueryServer = BaseKeeper{}
//
//// Balance implements the Query/Balance gRPC method
//func (k BaseKeeper) Balance(ctx context.Context, req *typesadapter.QueryBalanceRequest) (*typesadapter.QueryBalanceResponse, error) {
//
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.Address == "" {
//		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
//	}
//
//	if req.Denom == "" {
//		return nil, status.Error(codes.InvalidArgument, "invalid denom")
//	}
//
//	address, err := sdk.AccAddressFromBech32(req.Address)
//	if err != nil {
//		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
//	}
//
//	sdkCtx := sdk.UnwrapSDKContext(ctx)
//	coins := k.GetCoins(sdkCtx, address)
//	if coins == nil {
//		coins = sdk.NewCoins()
//	}
//	ret:=sdk.CoinAdapter{}
//	for _,c:=range coins{
//		if c.Denom==req.Denom{
//			ret.Denom=req.Denom
//			ret.Amount=sdk.NewIntFromBigInt(c.Amount.Int)
//			break
//		}
//	}
//
//	return &typesadapter.QueryBalanceResponse{Balance: &ret}, nil
//}
//
//// AllBalances implements the Query/AllBalances gRPC method
//func (k BaseKeeper) AllBalances(ctx context.Context, req *typesadapter.QueryAllBalancesRequest) (*typesadapter.QueryAllBalancesResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.Address == "" {
//		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
//	}
//
//	addr, err := sdk.AccAddressFromBech32(req.Address)
//	if err != nil {
//		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
//	}
//
//	sdkCtx := sdk.UnwrapSDKContext(ctx)
//
//	balances := sdk.NewCoins()
//	accountStore := k.getAccountStore(sdkCtx, addr)
//
//	pageRes, err := query.Paginate(accountStore, req.Pagination, func(_, value []byte) error {
//		var result sdk.Coin
//		err := k.cdc.Unmarshal(value, &result)
//		if err != nil {
//			return err
//		}
//		balances = append(balances, result)
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
//	}
//
//	return &typesadapter.QueryAllBalancesResponse{Balances: balances, Pagination: pageRes}, nil
//}
//
//// TotalSupply implements the Query/TotalSupply gRPC method
//func (k BaseKeeper) TotalSupply(ctx context.Context, req *typesadapter.QueryTotalSupplyRequest) (*typesadapter.QueryTotalSupplyResponse, error) {
//	sdkCtx := sdk.UnwrapSDKContext(ctx)
//	totalSupply, pageRes, err := k.GetPaginatedTotalSupply(sdkCtx, req.Pagination)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryTotalSupplyResponse{Supply: totalSupply, Pagination: pageRes}, nil
//}
//
//// SupplyOf implements the Query/SupplyOf gRPC method
//func (k BaseKeeper) SupplyOf(c context.Context, req *typesadapter.QuerySupplyOfRequest) (*typesadapter.QuerySupplyOfResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.Denom == "" {
//		return nil, status.Error(codes.InvalidArgument, "invalid denom")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//	supply := k.GetSupply(ctx, req.Denom)
//
//	return &typesadapter.QuerySupplyOfResponse{Amount: sdk.NewCoin(req.Denom, supply.Amount)}, nil
//}
//
//// Params implements the gRPC service handler for querying x/bank parameters.
//func (k BaseKeeper) Params(ctx context.Context, req *typesadapter.QueryParamsRequest) (*typesadapter.QueryParamsResponse, error) {
//	if req == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	sdkCtx := sdk.UnwrapSDKContext(ctx)
//	params := k.GetParams(sdkCtx)
//
//	return &typesadapter.QueryParamsResponse{Params: params}, nil
//}
//
//// DenomsMetadata implements Query/DenomsMetadata gRPC method.
//func (k BaseKeeper) DenomsMetadata(c context.Context, req *typesadapter.QueryDenomsMetadataRequest) (*typesadapter.QueryDenomsMetadataResponse, error) {
//	if req == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//	store := prefix.NewStore(ctx.KVStore(k.storeKey), typesadapter.DenomMetadataPrefix)
//
//	metadatas := []typesadapter.Metadata{}
//	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
//		var metadata typesadapter.Metadata
//		k.cdc.MustUnmarshal(value, &metadata)
//
//		metadatas = append(metadatas, metadata)
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDenomsMetadataResponse{
//		Metadatas:  metadatas,
//		Pagination: pageRes,
//	}, nil
//}
//
//// DenomMetadata implements Query/DenomMetadata gRPC method.
//func (k BaseKeeper) DenomMetadata(c context.Context, req *typesadapter.QueryDenomMetadataRequest) (*typesadapter.QueryDenomMetadataResponse, error) {
//	if req == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	if req.Denom == "" {
//		return nil, status.Error(codes.InvalidArgument, "invalid denom")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	metadata, found := k.GetDenomMetaData(ctx, req.Denom)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "client metadata for denom %s", req.Denom)
//	}
//
//	return &typesadapter.QueryDenomMetadataResponse{
//		Metadata: metadata,
//	}, nil
//}
