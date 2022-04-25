package lcd

import (
	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"io"
)

var (
	_ runtime.Marshaler = (*JSONMarshalAdapter)(nil)
)

type JSONMarshalAdapter struct {
	jsonPb *gateway.JSONPb
	codec  *codec.CodecProxy
}

func NewJSONMarshalAdapter(jsonPb *gateway.JSONPb, codec *codec.CodecProxy) *JSONMarshalAdapter {
	return &JSONMarshalAdapter{jsonPb: jsonPb, codec: codec}
}

func (m *JSONMarshalAdapter) Marshal(v interface{}) ([]byte, error) {
	if resp, ok := v.(context.CodecSensitive); ok {
		ret, err := resp.MarshalSensitive(m.codec)
		if nil == err {
			return ret, err
		}
	}
	return m.jsonPb.Marshal(v)
}

func (m *JSONMarshalAdapter) Unmarshal(data []byte, v interface{}) error {
	return m.jsonPb.Unmarshal(data, v)
}

func (m *JSONMarshalAdapter) NewDecoder(r io.Reader) runtime.Decoder {
	return m.jsonPb.NewDecoder(r)
}

func (m *JSONMarshalAdapter) NewEncoder(w io.Writer) runtime.Encoder {
	return m.jsonPb.NewEncoder(w)
}

func (m *JSONMarshalAdapter) ContentType() string {
	return m.jsonPb.ContentType()
}
