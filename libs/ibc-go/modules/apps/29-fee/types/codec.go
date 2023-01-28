package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/ibc 29-fee interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgPayPacketFee{}, "cosmos-sdk/MsgPayPacketFee", nil)
	cdc.RegisterConcrete(&MsgPayPacketFeeAsync{}, "cosmos-sdk/MsgPayPacketFeeAsync", nil)
	cdc.RegisterConcrete(&MsgRegisterPayee{}, "cosmos-sdk/MsgRegisterPayee", nil)
	cdc.RegisterConcrete(&MsgRegisterCounterpartyPayee{}, "cosmos-sdk/MsgRegisterCounterpartyPayee", nil)
}

// RegisterInterfaces register the 29-fee module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgPayPacketFee{},
		&MsgPayPacketFeeAsync{},
		&MsgRegisterPayee{},
		&MsgRegisterCounterpartyPayee{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.New()
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}
