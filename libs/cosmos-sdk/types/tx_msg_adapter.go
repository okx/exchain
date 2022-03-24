package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

const (
	IBCROUTER = "ibc"
)

type MsgProtoAdapter interface {
	Msg
	codec.ProtoMarshaler
}
type MsgAdapter interface {
	Msg
	proto.Message
}
