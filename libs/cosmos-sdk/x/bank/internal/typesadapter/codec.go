package typesadapter

import (
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

func RegisterInterface(registry interfacetypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgSend{},
	)
}
