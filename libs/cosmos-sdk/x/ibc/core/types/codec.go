package types

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	clienttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
	commitmenttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/23-commitment/types"
	ibctmtypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
)

// RegisterInterfaces registers x/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	connectiontypes.RegisterInterfaces(registry)
	channeltypes.RegisterInterfaces(registry)
	//solomachinetypes.RegisterInterfaces(registry)
	ibctmtypes.RegisterInterfaces(registry)
	//localhosttypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
}
