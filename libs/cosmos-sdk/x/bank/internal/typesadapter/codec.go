package typesadapter

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

var (
	cdc *codec.Codec
)

func init() {
	cdc = codec.New()
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
}

func RegisterInterface(registry interfacetypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgSend{},
	)
}
