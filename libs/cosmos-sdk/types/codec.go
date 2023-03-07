package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
)

const (
	// MsgInterfaceProtoName defines the protobuf name of the cosmos Msg interface
	MsgInterfaceProtoName = "cosmos.base.v1beta1.Msg"
)

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*MsgProtoAdapter)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}
