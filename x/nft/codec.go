package nft

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*txmsg.Msg)(nil),
		&MsgSend{},
	)
	registry.RegisterImplementations((*types2.MsgProtoAdapter)(nil), &MsgSend{})
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
