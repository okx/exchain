package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types"

	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

// RegisterLegacyAminoCodec registers the necessary x/ibc transfer interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgTransfer{}, "cosmos-sdk/MsgTransfer", nil)
}

// RegisterInterfaces register the ibc transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {

	registry.RegisterImplementations((*types.MsgProtoAdapter)(nil), &MsgTransfer{})
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgTransfer{},
	)

	registry.RegisterImplementations((*types.MsgProtoAdapter)(nil), &MsgTransfer{})
	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

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
	// AminoCdc is a amino codec created to support amino json compatible msgs.
	//AminoCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func SetMarshal(m *codec.CodecProxy) {
	Marshal = m
}
