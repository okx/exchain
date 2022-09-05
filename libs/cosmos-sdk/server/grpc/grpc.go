package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"

	"github.com/okex/exchain/libs/cosmos-sdk/server/grpc/gogoreflection"

	"google.golang.org/grpc"

	"github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/config"
)

// ServerStartTime defines the time duration that the server need to stay running after startup
// for the startup be considered successful
const ServerStartTime = 5 * time.Second

// StartGRPCServer starts a gRPC server on the given address.
func StartGRPCServer(cdc *codec.CodecProxy, interfaceReg jsonpb.AnyResolver, app types.ApplicationAdapter, cfg config.GRPCConfig) (*grpc.Server, error) {
	//clientCtx cliContext.CLIContext,
	maxSendMsgSize := cfg.MaxSendMsgSize
	if maxSendMsgSize == 0 {
		maxSendMsgSize = config.DefaultGRPCMaxSendMsgSize
	}

	maxRecvMsgSize := cfg.MaxRecvMsgSize
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = config.DefaultGRPCMaxRecvMsgSize
	}

	grpcSrv := grpc.NewServer(
		grpc.MaxSendMsgSize(maxSendMsgSize),
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
	)

	app.RegisterGRPCServer(grpcSrv)

	// Reflection allows consumers to build dynamic clients that can write to any
	// Cosmos SDK application without relying on application packages at compile
	// time.
	//err := reflection.Register(grpcSrv, reflection.Config{
	//	SigningModes: func() map[string]int32 {
	//		modes := make(map[string]int32, len(clientCtx.TxConfig.SignModeHandler().Modes()))
	//		for _, m := range clientCtx.TxConfig.SignModeHandler().Modes() {
	//			modes[m.String()] = (int32)(m)
	//		}
	//		return modes
	//	}(),
	//	ChainID:           clientCtx.ChainID,
	//	SdkConfig:         sdk.GetConfig(),
	//	InterfaceRegistry: clientCtx.InterfaceRegistry,
	//})
	//if err != nil {
	//	return nil, err
	//}

	// Reflection allows external clients to see what services and methods
	// the gRPC server exposes.
	gogoreflection.Register(grpcSrv)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)
	go func() {
		err = grpcSrv.Serve(listener)
		if err != nil {
			errCh <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return nil, err

	case <-time.After(ServerStartTime):
		// assume server started successfully
		return grpcSrv, nil
	}
}
