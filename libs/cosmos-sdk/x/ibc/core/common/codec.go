package common

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
)


func MustMarshalChannel(cdc *codec.MarshalProxy, c *types.Channel) []byte {
	ret := cdc.GetProtocMarshal().MustMarshalBinaryLengthPrefixed(c)
	return ret
}


func MarshalChannel(cdc *codec.MarshalProxy, c *types.Channel) ([]byte,error) {
	return cdc.GetProtocMarshal().MarshalBinaryLengthPrefixed(c)
}

func MustUnmarshalChannel(cdc *codec.MarshalProxy, data []byte) *types.Channel {
	var ret types.Channel
	cdc.GetProtocMarshal().MustUnmarshalBinaryLengthPrefixed(data, &ret)
	return &ret
}

func MustUnmarshalConnection(cdc *codec.MarshalProxy, data []byte) *types2.ConnectionEnd {
	var ret types2.ConnectionEnd
	cdc.GetProtocMarshal().MustUnmarshalBinaryLengthPrefixed(data, &ret)
	return &ret
}
func UnmarshalConnection(cdc *codec.MarshalProxy, data []byte) (*types2.ConnectionEnd,error){
	var ret types2.ConnectionEnd
	err:=cdc.GetProtocMarshal().UnmarshalBinaryLengthPrefixed(data, &ret)
	return &ret,err
}

func MustMarshalConnection(cdc *codec.MarshalProxy,c *types2.ConnectionEnd)[]byte{
	ret:=cdc.GetProtocMarshal().MustMarshalBinaryLengthPrefixed(c)
	return ret
}
