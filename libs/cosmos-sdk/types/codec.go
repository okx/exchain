package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*MsgProtoAdapter)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}
