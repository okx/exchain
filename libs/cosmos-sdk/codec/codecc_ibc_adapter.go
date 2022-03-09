package codec

import (
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
)

type CodecProxy struct {
	protoCodec *ProtoCodec
	cdc        *Codec
}

func NewCodecProxy(protoCodec *ProtoCodec, cdc *Codec) *CodecProxy {
	return &CodecProxy{protoCodec: protoCodec, cdc: cdc}
}

func (mp *CodecProxy) GetProtocMarshal() *ProtoCodec {
	return mp.protoCodec
}

type MarshalProxy struct {
	protocMarshal Marshaler
	cdc           *Codec
}

func NewMarshalProxy(protocMarshal Marshaler, cdc *Codec) *MarshalProxy {
	return &MarshalProxy{protocMarshal: protocMarshal, cdc: cdc}
}

func (mp *MarshalProxy) InterfaceRegistry() types.InterfaceRegistry {
	return mp.protocMarshal.(ProtoCodecMarshaler).InterfaceRegistry()
}

func (mp *MarshalProxy) Marshal(v proto.Message) ([]byte, error) {
	return mp.protocMarshal.MarshalInterfaceJSON(v)
}
func (mp *MarshalProxy) GetCdc() *Codec {
	return mp.cdc
}
func (mp *MarshalProxy) GetProtocMarshal() Marshaler {
	return mp.protocMarshal
}

func (mp *MarshalProxy) MustMarshal(v proto.Message) []byte {
	var ret []byte
	var err error
	ret, err = mp.protocMarshal.MarshalInterface(v)
	if nil != err {
		panic(err)
	}
	return ret
}

func (mp *MarshalProxy) UnMarshal(data []byte, ptr interface{}) error {
	return mp.protocMarshal.UnmarshalInterface(data, ptr)
}

func (mp *MarshalProxy) MustUnMarshal(data []byte, ptr interface{}) {
	err := mp.UnMarshal(data, ptr)
	if nil != err {
		err = mp.cdc.UnmarshalBinaryBare(data, ptr)
		if nil != err {
			panic(err)
		}
	}
}
