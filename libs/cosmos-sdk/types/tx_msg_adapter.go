package types

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

//import (
//	"github.com/gogo/protobuf/proto"
//	"github.com/okex/exchain/libs/cosmos-sdk/codec"
//	mempl "github.com/okex/exchain/libs/tendermint/mempool"
//	"math/big"
//)

const (
	IBCROUTER = "ibc"
)

var (
	//	_   Msg = (*RelayMsg)(nil)
	cdcRelayer = codec.New()
	//	_   Tx  = (*RelayTxMsg)(nil)
)

type MsgProtoAdapter interface {
	Msg
	codec.ProtoMarshaler
}
type MsgAdapter interface {
	Msg
	proto.Message
}

type RelayMsgWrapper struct {
	Msgs []*RelayMsg
}
type BytesWrapper struct {
	// json stdtx
	Data []byte
}

func NewBytesWrapper(data []byte) *BytesWrapper {
	return &BytesWrapper{Data: data}
}

func (this *BytesWrapper) Marshal() ([]byte, error) {
	return json.Marshal(this)
}
func (this *BytesWrapper) UnmarshalToTx(data []byte) ([]byte, error) {
	err := json.Unmarshal(data, this)
	if nil != err {
		return nil, err
	}
	return this.Data, nil
}

func (this *RelayMsgWrapper) UnMarshal(bs []byte) error {
	return json.Unmarshal(bs, this)
}

func (this *RelayMsgWrapper) Marshal() ([]byte, error) {
	return json.Marshal(this)
}

type RelayMsg struct {
	Bytes     []byte
	Signers   []AccAddress
	RouterStr string
	TypeStr   string
	MsgType   string
}
type RelayMsgOption func(m *RelayMsg)

var WithRouter = func(r string) RelayMsgOption {
	return func(m *RelayMsg) {
		m.RouterStr = r
	}
}
var WithType = func(t string) RelayMsgOption {
	return func(m *RelayMsg) {
		m.TypeStr = t
	}
}
var WithMsgDetailType = func(msgDetailType string) RelayMsgOption {
	return func(m *RelayMsg) {
		m.MsgType = msgDetailType
	}
}

func NewRelayMsg(data []byte, ss []AccAddress, ops ...RelayMsgOption) *RelayMsg {
	if len(ops) == 0 {
		ops = []RelayMsgOption{WithRouter(IBCROUTER), WithType(IBCROUTER)}
	}
	ret := &RelayMsg{}
	ret.Bytes = data
	ret.Signers = ss
	for _, o := range ops {
		o(ret)
	}
	return ret
}

func (r *RelayMsg) Route() string {
	if len(r.RouterStr) == 0 {
		return IBCROUTER
	}
	return r.RouterStr
}

func (r *RelayMsg) Type() string {
	if len(r.TypeStr) == 0 {
		return IBCROUTER
	}
	return r.TypeStr
}

func (r *RelayMsg) ValidateBasic() error {
	return nil
}

func (r *RelayMsg) GetSignBytes() []byte {
	ret, err := cdcRelayer.MarshalJSON(r)
	if nil != err {
		panic(err)
	}
	return ret
}

func (r *RelayMsg) GetSigners() []AccAddress {
	return r.Signers
}
