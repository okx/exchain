package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	auth "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// ModuleCdc references the global interchain accounts module codec. Note, the codec
// should ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to interchain accounts and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func init() {
	RegisterCodec(auth.ModuleCdc)
}

// RegisterInterfaces registers the concrete InterchainAccount implementation against the associated
// x/auth AccountI and GenesisAccount interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&AminoInterchainAccount{}, "cosmos-sdk/InterchainAccount", nil)
}

// SerializeCosmosTx serializes a slice of sdk.Msg's using the CosmosTx type. The sdk.Msg's are
// packed into Any's and inserted into the Messages field of a CosmosTx. The proto marshaled CosmosTx
// bytes are returned. Only the ProtoCodec is supported for serializing messages.
func SerializeCosmosTx(cdc *codec.CodecProxy, msgs []sdk.MsgAdapter) (bz []byte, err error) {
	// only ProtoCodec is supported

	msgAnys := make([]*codectypes.Any, len(msgs))

	for i, msg := range msgs {
		msgAnys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
	}

	cosmosTx := &CosmosTx{
		Messages: msgAnys,
	}

	bz, err = cdc.GetProtocMarshal().Marshal(cosmosTx)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// DeserializeCosmosTx unmarshals and unpacks a slice of transaction bytes
// into a slice of sdk.Msg's. Only the ProtoCodec is supported for message
// deserialization.
func DeserializeCosmosTx(cdc *codec.CodecProxy, data []byte) ([]sdk.MsgAdapter, error) {
	// only ProtoCodec is supported
	var cosmosTx CosmosTx
	if err := cdc.GetProtocMarshal().Unmarshal(data, &cosmosTx); err != nil {
		return nil, err
	}

	msgs := make([]sdk.MsgAdapter, len(cosmosTx.Messages))
	for i, any := range cosmosTx.Messages {
		var msg sdk.MsgAdapter

		err := cdc.GetProtocMarshal().UnpackAny(any, &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}
