package port

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
)

// Name returns the IBC port ICS name.
func Name() string {
	return types.SubModuleName
}

// GetQueryCmd returns the root query command for IBC ports.
//func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
//	return cli.GetQueryCmd(cdc, reg)
//}

// RegisterQueryService registers the gRPC query service for IBC ports.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
