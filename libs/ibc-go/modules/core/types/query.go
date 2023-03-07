package types

import (
	"github.com/gogo/protobuf/grpc"
	client "github.com/okx/okbchain/libs/ibc-go/modules/core/02-client"
	clienttypes "github.com/okx/okbchain/libs/ibc-go/modules/core/02-client/types"
	connection "github.com/okx/okbchain/libs/ibc-go/modules/core/03-connection"
	connectiontypes "github.com/okx/okbchain/libs/ibc-go/modules/core/03-connection/types"
	channel "github.com/okx/okbchain/libs/ibc-go/modules/core/04-channel"
	channeltypes "github.com/okx/okbchain/libs/ibc-go/modules/core/04-channel/types"
	port "github.com/okx/okbchain/libs/ibc-go/modules/core/05-port"
	porttypes "github.com/okx/okbchain/libs/ibc-go/modules/core/05-port/types"
)

// QueryService defines the IBC interfaces that the gRPC query server must implement
type QueryService interface {
	clienttypes.QueryServer
	connectiontypes.QueryServer
	channeltypes.QueryServer
	porttypes.QueryServer
	QueryServer
}

// RegisterQueryService registers each individual IBC submodule query service
func RegisterQueryService(server grpc.Server, queryService QueryService) {
	client.RegisterQueryService(server, queryService)
	connection.RegisterQueryService(server, queryService)
	channel.RegisterQueryService(server, queryService)
	port.RegisterQueryService(server, queryService)
	RegisterQueryServer(server, queryService)
}
