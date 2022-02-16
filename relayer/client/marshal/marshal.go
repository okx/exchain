package marshal

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
)

type OECMarshal struct {
	interfaceRegistry types.InterfaceRegistry
	modBasic          module.BasicManager
	codec             *codec.Codec
}

func NewOECMarshal(codec *codec.Codec, modules []module.AppModuleBasic) *OECMarshal {
	ret := &OECMarshal{
		interfaceRegistry: types.NewInterfaceRegistry(),
		modBasic:          module.NewBasicManager(modules...),
		codec:             codec,
	}
	return ret
}

func (d *OECMarshal) Init() {
	d.modBasic.RegisterCodec(d.codec)
	d.modBasic.RegisterInterfaces(d.interfaceRegistry)
}

func (d *OECMarshal) Marshal(data interface{}) ([]byte, error) {
	return d.codec.MarshalJSON(data)
}

func (d *OECMarshal) Unmarshal(data []byte, ptr interface{}) error {
	return d.codec.UnmarshalJSON(data, ptr)
}
