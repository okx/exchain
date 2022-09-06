package tx

//import (
//	"context"
//	"fmt"
//	"strings"
//
//	cliContext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
//	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
//
//	gogogrpc "github.com/gogo/protobuf/grpc"
//	proto "github.com/gogo/protobuf/proto"
//	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
//	types "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)
//
//// baseAppSimulateFn is the signature of the Baseapp#Simulate function.
//type baseAppSimulateFn func(txBytes []byte) (types.GasInfo, *types.Result, error)
//
//// txServer is the server for the protobuf Tx service.
//type txServer struct {
//	clientCtx         cliContext.CLIContext
//	simulate          baseAppSimulateFn
//	interfaceRegistry codectypes.InterfaceRegistry
//}
//
//// NewTxServer creates a new Tx service server.
//func NewTxServer(clientCtx cliContext.CLIContext, simulate baseAppSimulateFn, interfaceRegistry codectypes.InterfaceRegistry) ServiceServer {
//	return txServer{
//		clientCtx:         clientCtx,
//		simulate:          simulate,
//		interfaceRegistry: interfaceRegistry,
//	}
//}
//
//var _ ServiceServer = txServer{}
//
//const (
//	eventFormat = "{eventType}.{eventAttribute}={value}"
//)
//
//// TxsByEvents implements the ServiceServer.TxsByEvents RPC method.
//func (s txServer) GetTxsEvent(ctx context.Context, req *GetTxsEventRequest) (*GetTxsEventResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
//	}
//
//	page, limit, err := query.ParsePagination(req.Pagination)
//	if err != nil {
//		return nil, err
//	}
//	orderBy := parseOrderBy(req.OrderBy)
//
//	if len(req.Events) == 0 {
//		return nil, status.Error(codes.InvalidArgument, "must declare at least one event to search")
//	}
//
//	for _, event := range req.Events {
//		if !strings.Contains(event, "=") || strings.Count(event, "=") > 1 {
//			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid event; event %s should be of the format: %s", event, eventFormat))
//		}
//	}
//
//	//result, err := utils.QueryTxsByEvents(s.clientCtx, req.Events, page, limit)
//	result, err := QueryTxsByEvents(s.clientCtx, req.Events, page, limit, orderBy)
//	if err != nil {
//		return nil, err
//	}
//
//	// Create a proto codec, we need it to unmarshal the tx bytes.
//	txsList := make([]*Tx, len(result.Txs))
//
//	for i, tx := range result.Txs {
//		protoTx, ok := tx.Tx.GetCachedValue().(*Tx)
//		if !ok {
//			return nil, status.Errorf(codes.Internal, "expected %T, got %T", Tx{}, tx.Tx.GetCachedValue())
//		}
//
//		txsList[i] = protoTx
//	}
//
//	return &GetTxsEventResponse{
//		Txs:         txsList,
//		TxResponses: result.Txs,
//		Pagination: &query.PageResponse{
//			Total: result.TotalCount,
//		},
//	}, nil
//}
//
//// Simulate implements the ServiceServer.Simulate RPC method.
//func (s txServer) Simulate(ctx context.Context, req *SimulateRequest) (*SimulateResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "invalid empty tx")
//	}
//
//	txBytes := req.TxBytes
//	if txBytes == nil && req.Tx != nil {
//		// This block is for backwards-compatibility.
//		// We used to support passing a `Tx` in req. But if we do that, sig
//		// verification might not pass, because the .Marshal() below might not
//		// be the same marshaling done by the client.
//		var err error
//		txBytes, err = proto.Marshal(req.Tx)
//		if err != nil {
//			return nil, status.Errorf(codes.InvalidArgument, "invalid tx; %v", err)
//		}
//	}
//
//	if txBytes == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty txBytes is not allowed")
//	}
//
//	gasInfo, result, err := s.simulate(txBytes)
//	if err != nil {
//		return nil, err
//	}
//
//	return &SimulateResponse{
//		GasInfo: &gasInfo,
//		Result:  result,
//	}, nil
//}
//
//// GetTx implements the ServiceServer.GetTx RPC method.
//func (s txServer) GetTx(ctx context.Context, req *GetTxRequest) (*GetTxResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
//	}
//
//	// TODO We should also check the proof flag in gRPC header.
//	// https://github.com/cosmos/cosmos-sdk/issues/7036.
//	result, err := QueryTx(s.clientCtx, req.Hash)
//	if err != nil {
//		return nil, err
//	}
//
//	protoTx, ok := result.Tx.GetCachedValue().(*Tx)
//	if !ok {
//		return nil, status.Errorf(codes.Internal, "expected %T, got %T", Tx{}, result.Tx.GetCachedValue())
//	}
//
//	return &GetTxResponse{
//		Tx:         protoTx,
//		TxResponse: result,
//	}, nil
//}
//
//func (s txServer) BroadcastTx(ctx context.Context, req *BroadcastTxRequest) (*BroadcastTxResponse, error) {
//	_, err := cliContext.TxServiceBroadcast(ctx, s.clientCtx, req)
//	if nil != err {
//		return nil, err
//	}
//	return nil, nil
//}
//
//// RegisterTxService registers the tx service on the gRPC router.
//func RegisterTxService(
//	qrt gogogrpc.Server,
//	clientCtx cliContext.CLIContext,
//	simulateFn baseAppSimulateFn,
//	interfaceRegistry codectypes.InterfaceRegistry,
//) {
//	RegisterServiceServer(
//		qrt,
//		NewTxServer(clientCtx, simulateFn, interfaceRegistry),
//	)
//}
//
//func parseOrderBy(orderBy OrderBy) string {
//	switch orderBy {
//	case OrderBy_ORDER_BY_ASC:
//		return "asc"
//	case OrderBy_ORDER_BY_DESC:
//		return "desc"
//	default:
//		return "" // Defaults to Tendermint's default, which is `asc` now.
//	}
//}
