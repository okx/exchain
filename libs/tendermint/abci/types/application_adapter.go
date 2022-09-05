package types

import (
	gogogrpc "github.com/gogo/protobuf/grpc"
)

type ApplicationAdapter interface {
	RegisterGRPCServer(gogogrpc.Server)
}
