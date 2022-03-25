package common

import (
	"errors"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// TODO
func UnmarshalMsgAdapter(cdc *codec.CodecProxy, data []byte) (sdk.MsgProtoAdapter, error) {
	var ret sdk.MsgProtoAdapter
	err := cdc.GetProtocMarshal().UnmarshalInterface(data, &ret)
	return ret, err
}

// TODO
func UnmarshalGuessss(cdc *codec.CodecProxy, data []byte, guesss ...sdk.MsgProtoAdapter) (sdk.MsgProtoAdapter, error) {
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
