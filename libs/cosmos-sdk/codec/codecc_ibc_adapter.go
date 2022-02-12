package codec

import (
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
)

type MarshalProxy struct {
	protocMarshal Marshaler
	cdc           *Codec
}

func NewMarshalProxy(protocMarshal Marshaler, cdc *Codec) *MarshalProxy {
	return &MarshalProxy{protocMarshal: protocMarshal, cdc: cdc}
}

func (this *MarshalProxy) InterfaceRegistry() types.InterfaceRegistry {
	return this.protocMarshal.(ProtoCodecMarshaler).InterfaceRegistry()
}

func (this *MarshalProxy) Marshal(v proto.Message) ([]byte, error) {
	return this.protocMarshal.MarshalInterfaceJSON(v)
}
func (this *MarshalProxy) GetCdc() *Codec {
	return this.cdc
}
func (this *MarshalProxy) GetProtocMarshal() Marshaler {
	return this.protocMarshal
}

func (this *MarshalProxy) MustMarshal(v proto.Message) []byte {
	var ret []byte
	var err error
	ret, err = this.protocMarshal.MarshalInterface(v)
	if nil != err {
		panic(err)
	}
	return ret
}

func (this *MarshalProxy) UnMarshal(data []byte, ptr interface{}) error {
	return this.protocMarshal.UnmarshalInterface(data, ptr)
}

func (this *MarshalProxy) MustUnMarshal(data []byte, ptr interface{}) {
	err := this.UnMarshal(data, ptr)
	if nil != err {
		err = this.cdc.UnmarshalBinaryBare(data, ptr)
		if nil != err {
			panic(err)
		}
	}
}
