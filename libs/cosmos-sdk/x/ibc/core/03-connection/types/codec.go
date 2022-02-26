package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/msgservice"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.core.connection.v1.ConnectionI",
		(*exported.ConnectionI)(nil),
		&ConnectionEnd{},
	)
	registry.RegisterInterface(
		"ibc.core.connection.v1.CounterpartyConnectionI",
		(*exported.CounterpartyConnectionI)(nil),
		&Counterparty{},
	)
	registry.RegisterInterface(
		"ibc.core.connection.v1.Version",
		(*exported.Version)(nil),
		&Version{},
	)
	registry.RegisterImplementations(
		(*sdk.MsgAdapter)(nil),
		&MsgConnectionOpenInit{},
		&MsgConnectionOpenTry{},
		&MsgConnectionOpenAck{},
		&MsgConnectionOpenConfirm{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
func RegistCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&ConnectionEnd{}, "end", nil)
}

//var (
//	// SubModuleCdc references the global x/ibc/core/03-connection module codec. Note, the codec should
//	// ONLY be used in certain instances of tests and for JSON encoding.
//	//
//	// The actual codec used for serialization should be provided to x/ibc/core/03-connection and
//	// defined at the application level.
//	SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
//)
//
