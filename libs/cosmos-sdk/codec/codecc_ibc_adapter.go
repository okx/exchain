package codec

import (
	"github.com/gogo/protobuf/proto"
)

type CodecProxy struct {
	protoCodec *ProtoCodec
	cdc        *Codec
}

func NewCodecProxy(protoCodec *ProtoCodec, cdc *Codec) *CodecProxy {
	return &CodecProxy{protoCodec: protoCodec, cdc: cdc}
}

func (mp *CodecProxy) GetCdc() *Codec {
	return mp.cdc
}

func (mp *CodecProxy) GetProtocMarshal() *ProtoCodec {
	return mp.protoCodec
}

func (mp *CodecProxy) MustUnMarshal(data []byte, ptr interface{}) {
	err := mp.UnMarshal(data, ptr)
	if nil != err {
		err = mp.cdc.UnmarshalBinaryBare(data, ptr)
		if nil != err {
			panic(err)
		}
	}
}

func (mp *CodecProxy) UnMarshal(data []byte, ptr interface{}) error {
	return mp.protoCodec.UnmarshalInterface(data, ptr)
}

func (mp *CodecProxy) MustMarshal(v proto.Message) []byte {
	var ret []byte
	var err error
	ret, err = mp.protoCodec.MarshalInterface(v)
	if nil != err {
		panic(err)
	}
	return ret
}

/////////////
