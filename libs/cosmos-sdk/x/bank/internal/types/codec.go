package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)

	cdc.RegisterConcreteUnmarshaller("cosmos-sdk/MsgSend", func(c *amino.Codec, bytes []byte) (interface{}, int, error) {
		var msg MsgSend
		err := msg.UnmarshalFromAmino(bytes)
		if err != nil {
			return nil, 0, err
		}
		return msg, len(bytes), nil
	})
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
