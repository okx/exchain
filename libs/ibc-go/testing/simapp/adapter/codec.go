package adapter

import "github.com/okex/exchain/libs/cosmos-sdk/codec"

var (
	//amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/ibc-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc transfer and
	// defined at the application level.
	//ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	ModuleCdc = codec.New()
	Marshal   *codec.CodecProxy
)

func init() {
	codec.RegisterCrypto(ModuleCdc)
}

func SetMarshal(m *codec.CodecProxy) {
	Marshal = m
}
