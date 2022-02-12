package common

import (
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func DefaultMarshal(c *codec.MarshalProxy, data proto.Message) ([]byte, error) {
	return c.GetProtocMarshal().MarshalInterface(data)
}

func UnmarshalMsgAdapter(cdc *codec.MarshalProxy, data []byte) (sdk.MsgAdapter, error) {
	var ret sdk.MsgAdapter
	return ret, cdc.GetProtocMarshal().UnmarshalInterface(data, &ret)
}
