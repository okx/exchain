package connection

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/client/cli"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	"github.com/spf13/cobra"
)

// Name returns the IBC connection ICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for the IBC connections.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the IBC connections.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for IBC connections.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
