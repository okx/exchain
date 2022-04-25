package context

import "github.com/okex/exchain/libs/cosmos-sdk/codec"

type TxRequest interface {
	GetData() []byte
	GetModeDetail() int32
}

type TxResponse interface {
	HandleResponse(codec *codec.CodecProxy, data interface{}) interface{}
}

type CodecSensitive interface {
	MarshalSensitive(proxy *codec.CodecProxy) ([]byte, error)
}
