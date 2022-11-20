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

// MsgTypeURL returns the TypeURL of a `sdk.Msg`.
func MsgTypeURL(msg proto.Message) string {
	return "/" + proto.MessageName(msg)
}
