package types

import (
	gogogrpc "github.com/gogo/protobuf/grpc"
)

// TODO,循环依赖,换个位置
type ApplicationAdapter interface {
	RegisterGRPCServer(gogogrpc.Server)
	//RegisterTxService(clientCtx cliContext.CLIContext)
}
