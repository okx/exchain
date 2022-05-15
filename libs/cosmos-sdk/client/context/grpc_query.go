package context

import (
	gocontext "context"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	grpctypes "github.com/okex/exchain/libs/cosmos-sdk/types/grpc"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strconv"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
)

var _ gogogrpc.ClientConn = CLIContext{}

var protoCodec = encoding.GetCodec(proto.Name)

// Invoke implements the grpc ClientConn.Invoke method
func (ctx CLIContext) Invoke(grpcCtx gocontext.Context, method string, req, reply interface{}, opts ...grpc.CallOption) (err error) {
	// Two things can happen here:
	// 1. either we're broadcasting a Tx, in which call we call Tendermint's broadcast endpoint directly,
	// 2. or we are querying for state, in which case we call ABCI's Query.

	// In both cases, we don't allow empty request args (it will panic unexpectedly).
	if reflect.ValueOf(req).IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "request cannot be nil")
	}

	// Case 1. Broadcasting a Tx.
	if reqProto, ok := req.(TxRequest); ok {
		broadcastRes, err := TxServiceBroadcast(grpcCtx, ctx, reqProto)
		if err != nil {
			return err
		}
		if resp, ok := reply.(TxResponse); ok {
			reply = resp.HandleResponse(ctx.CodecProy, broadcastRes)
		} else {
			reply = broadcastRes
		}

		return err
	}

	// Case 2. Querying state.
	reqBz, err := protoCodec.Marshal(req)
	if err != nil {
		return err
	}

	// parse height header
	md, _ := metadata.FromOutgoingContext(grpcCtx)
	if heights := md.Get(grpctypes.GRPCBlockHeightHeader); len(heights) > 0 {
		height, err := strconv.ParseInt(heights[0], 10, 64)
		if err != nil {
			return err
		}
		if height < 0 {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"client.Context.Invoke: height (%d) from %q must be >= 0", height, grpctypes.GRPCBlockHeightHeader)
		}

		ctx = ctx.WithHeight(height)
	}

	abciReq := abci.RequestQuery{
		Path:   method,
		Data:   reqBz,
		Height: ctx.Height,
	}

	res, err := ctx.QueryABCI(abciReq)
	if err != nil {
		return err
	}

	err = protoCodec.Unmarshal(res.Value, reply)
	if err != nil {
		return err
	}

	// Create header metadata. For now the headers contain:
	// - block height
	// We then parse all the call options, if the call option is a
	// HeaderCallOption, then we manually set the value of that header to the
	// metadata.
	md = metadata.Pairs(grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(res.Height, 10))
	for _, callOpt := range opts {
		header, ok := callOpt.(grpc.HeaderCallOption)
		if !ok {
			continue
		}

		*header.HeaderAddr = md
	}

	if ctx.InterfaceRegistry != nil {
		return types.UnpackInterfaces(reply, ctx.InterfaceRegistry)
	}

	return nil
}

// NewStream implements the grpc ClientConn.NewStream method
func (CLIContext) NewStream(gocontext.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("streaming rpc not supported")
}

func TxServiceBroadcast(grpcCtx context.Context, clientCtx CLIContext, req TxRequest) (interface{}, error) {
	if req == nil || req.GetData() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid empty tx")
	}

	clientCtx = clientCtx.WithBroadcastMode(normalizeBroadcastMode(req.GetModeDetail()))
	ret, err := clientCtx.BroadcastTx(req.GetData())
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// normalizeBroadcastMode converts a broadcast mode into a normalized string
// to be passed into the clientCtx.
func normalizeBroadcastMode(mode int32) string {
	switch mode {
	case 3:
		return "async"
	case 1:
		return "block"
	case 2:
		return "sync"
	default:
		return "unspecified"
	}
}
