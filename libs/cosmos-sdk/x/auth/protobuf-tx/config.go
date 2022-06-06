package protobuf_tx

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/types/pb-tx"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	signing "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
)

type config struct {
	handler signing.SignModeHandler
	decoder ibctx.PbTxDecoder
	encoder ibctx.PbTxEncoder
	//jsonDecoder ibctx.PbTxEncoder
	jsonEncoder ibctx.PbTxEncoder
	protoCodec  codec.ProtoCodecMarshaler
}

// NewTxConfig returns a new protobuf TxConfig using the provided ProtoCodec and sign modes. The
// first enabled sign mode will become the default sign mode.
func NewTxConfig(protoCodec codec.ProtoCodecMarshaler, enabledSignModes []signing2.SignMode) client.TxConfig {
	return &config{
		handler: makeSignModeHandler(enabledSignModes),
		decoder: ProtoBufTxDecoder(protoCodec),
		encoder: IbcTxEncoder(),
		//jsonDecoder: DefaultJSONTxDecoder(protoCodec),
		jsonEncoder: DefaultJSONTxEncoder(protoCodec),
		protoCodec:  protoCodec,
	}
}

func (g config) NewTxBuilder() client.TxBuilder {
	return newBuilder()
}

// WrapTxBuilder returns a builder from provided transaction
func (g config) WrapTxBuilder(newTx ibctx.Tx) (client.TxBuilder, error) {
	newBuilder, ok := newTx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, newTx)
	}

	return newBuilder, nil
}

func (g config) SignModeHandler() signing.SignModeHandler {
	return g.handler
}

func (g config) TxEncoder() ibctx.PbTxEncoder {
	return g.encoder
}

func (g config) TxDecoder() ibctx.PbTxDecoder {
	return g.decoder
}

//
func (g config) TxJSONEncoder() ibctx.PbTxEncoder {
	return g.jsonEncoder
}

//
//func (g config) TxJSONDecoder() sdk.TxDecoder {
//	return g.jsonDecoder
//}
