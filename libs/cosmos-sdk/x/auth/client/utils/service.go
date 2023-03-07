package utils

import (
	"context"
	"fmt"
	"strings"

	gogogrpc "github.com/gogo/protobuf/grpc"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cliContext "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	codectypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types"
	typeadapter "github.com/okx/okbchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/query"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/tx"
)

var _ tx.ServiceServer = txServer{}

const (
	eventFormat = "{eventType}.{eventAttribute}={value}"
)

type baseAppSimulateFn func(txBytes []byte) (types.GasInfo, *types.Result, error)

// txServer is the server for the protobuf Tx service.
type txServer struct {
	clientCtx         cliContext.CLIContext
	simulate          baseAppSimulateFn
	interfaceRegistry codectypes.InterfaceRegistry
}

// NewTxServer creates a new Tx service server.
func NewTxServer(clientCtx cliContext.CLIContext, simulate baseAppSimulateFn, interfaceRegistry codectypes.InterfaceRegistry) tx.ServiceServer {
	return txServer{
		clientCtx:         clientCtx,
		simulate:          simulate,
		interfaceRegistry: interfaceRegistry,
	}
}

func (t txServer) Simulate(ctx context.Context, req *tx.SimulateRequest) (*tx.SimulateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid empty tx")
	}

	txBytes := req.TxBytes
	if txBytes == nil && req.Tx != nil {
		// This block is for backwards-compatibility.
		// We used to support passing a `Tx` in req. But if we do that, sig
		// verification might not pass, because the .Marshal() below might not
		// be the same marshaling done by the client.
		var err error
		txBytes, err = proto.Marshal(req.Tx)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid tx; %v", err)
		}
	}

	if txBytes == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty txBytes is not allowed")
	}

	gasInfo, res, err := t.simulate(txBytes)
	if err != nil {
		return nil, err
	}

	return &tx.SimulateResponse{
		GasInfo: &typeadapter.GasInfo{
			GasWanted: gasInfo.GasWanted,
			GasUsed:   gasInfo.GasUsed,
		},
		Result: ConvCM39SimulateResultTCM40(res),
	}, nil
}

func (t txServer) GetTx(ctx context.Context, req *tx.GetTxRequest) (*tx.GetTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if len(req.Hash) == 0 {
		return nil, status.Error(codes.InvalidArgument, "tx hash cannot be empty")
	}

	result, err := Query40Tx(t.clientCtx, req.Hash)
	if nil != err {
		return nil, err
	}

	protoTx, ok := result.Tx.GetCachedValue().(*tx.Tx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "expected %T, got %T", tx.Tx{}, result.Tx.GetCachedValue())
	}

	return &tx.GetTxResponse{
		Tx:         protoTx,
		TxResponse: result,
	}, nil
}

func (t txServer) BroadcastTx(ctx context.Context, request *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error) {
	resp, err := cliContext.TxServiceBroadcast(ctx, t.clientCtx, request)
	if nil != err {
		return nil, err
	}
	ret := new(tx.BroadcastTxResponse)
	ret.HandleResponse(t.clientCtx.CodecProy, resp)
	return ret, nil
}

func (t txServer) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	page := 1
	// Tendermint node.TxSearch that is used for querying txs defines pages starting from 1,
	// so we default to 1 if not provided in the request.
	limit := query.DefaultLimit

	if len(req.Events) == 0 {
		return nil, status.Error(codes.InvalidArgument, "must declare at least one event to search")
	}

	for _, event := range req.Events {
		if !strings.Contains(event, "=") || strings.Count(event, "=") > 1 {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid event; event %s should be of the format: %s", event, eventFormat))
		}
	}

	result, err := Query40TxsByEvents(t.clientCtx, req.Events, page, limit)
	if err != nil {
		return nil, err
	}

	// Create a proto codec, we need it to unmarshal the tx bytes.
	txsList := make([]*tx.Tx, len(result.Txs))

	for i, txx := range result.Txs {
		protoTx, ok := txx.Tx.GetCachedValue().(*tx.Tx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "expected %T, got %T", tx.Tx{}, txx.Tx.GetCachedValue())
		}

		txsList[i] = protoTx
	}

	return &tx.GetTxsEventResponse{
		Txs:         txsList,
		TxResponses: result.Txs,
	}, nil
}

// RegisterTxService registers the tx service on the gRPC router.
func RegisterTxService(
	qrt gogogrpc.Server,
	clientCtx cliContext.CLIContext,
	simulateFn baseAppSimulateFn,
	interfaceRegistry codectypes.InterfaceRegistry,
) {
	tx.RegisterServiceServer(
		qrt,
		NewTxServer(clientCtx, simulateFn, interfaceRegistry),
	)
}
