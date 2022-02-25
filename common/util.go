package common

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func DefaultMarshal(c *codec.MarshalProxy, data proto.Message) ([]byte, error) {
	return c.GetProtocMarshal().MarshalInterface(data)
}

func UnmarshalMsgAdapter(cdc *codec.MarshalProxy, data []byte) (sdk.MsgProtoAdapter, error) {
	var ret sdk.MsgProtoAdapter
	err := cdc.GetProtocMarshal().UnmarshalInterface(data, &ret)
	return ret, err
	//if nil == err {
	//	return ret, err
	//}
	//return ret, cdc.GetProtocMarshal().MustUnmarshalBinaryBare(data, &ret)
}

func UnmarshalGuessss(cdc *codec.MarshalProxy, data []byte, guesss ...sdk.MsgProtoAdapter) (sdk.MsgProtoAdapter, error) {
	ret, err := UnmarshalMsgAdapter(cdc, data)
	if nil == err {
		return ret, nil
	}
	for _, v := range guesss {
		err := cdc.GetProtocMarshal().UnmarshalBinaryBare(data, v)
		if nil == err {
			return v, err
		}
	}
	return nil, errors.New("asd")
}
