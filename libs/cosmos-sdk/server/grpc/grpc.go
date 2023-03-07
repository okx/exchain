package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/okx/okbchain/libs/tendermint/node"

	app2 "github.com/okx/okbchain/libs/cosmos-sdk/server/types"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/viper"

	reflection "github.com/okx/okbchain/libs/cosmos-sdk/server/grpc/reflection/v2alpha1"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"

	"github.com/okx/okbchain/libs/cosmos-sdk/server/grpc/gogoreflection"

	"google.golang.org/grpc"

	"github.com/okx/okbchain/libs/tendermint/config"
)

// ServerStartTime defines the time duration that the server need to stay running after startup
// for the startup be considered successful
const ServerStartTime = 5 * time.Second

// StartGRPCServer starts a gRPC server on the given address.
func StartGRPCServer(cdc *codec.CodecProxy, interfaceReg jsonpb.AnyResolver, app app2.ApplicationAdapter, cfg config.GRPCConfig, tmNode *node.Node) (*grpc.Server, error) {
	txCfg := utils.NewPbTxConfig(interfaceReg.(interfacetypes.InterfaceRegistry))

	cliCtx := context.NewCLIContext().WithProxy(cdc).WithInterfaceRegistry(interfaceReg.(interfacetypes.InterfaceRegistry)).WithTrustNode(true)
	if tmNode != nil {
		cliCtx = cliCtx.WithChainID(tmNode.ConsensusState().GetState().ChainID)
	} else {
		cliCtx = cliCtx.WithChainID(viper.GetString(flags.FlagChainID))
	}

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

	app.RegisterTxService(cliCtx)
	app.RegisterGRPCServer(grpcSrv)

	// Reflection allows consumers to build dynamic clients that can write to any
	// Cosmos SDK application without relying on application packages at compile
	// time.
	err := reflection.Register(grpcSrv, reflection.Config{
		SigningModes: func() map[string]int32 {
			modes := make(map[string]int32, len(txCfg.SignModeHandler().Modes()))
			for _, m := range txCfg.SignModeHandler().Modes() {
				modes[m.String()] = (int32)(m)
			}
			return modes
		}(),
		ChainID:           cliCtx.ChainID,
		SdkConfig:         sdk.GetConfig(),
		InterfaceRegistry: cliCtx.InterfaceRegistry,
	})
	if err != nil {
		return nil, err
	}

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
