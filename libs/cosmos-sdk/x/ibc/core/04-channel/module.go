package channel

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/client/cli"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/spf13/cobra"
)

// Name returns the IBC channel ICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for IBC channels.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for IBC channels.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for IBC channels.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
