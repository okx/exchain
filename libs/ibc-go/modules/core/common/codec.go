package common

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
)


func MustMarshalChannel(cdc *codec.MarshalProxy, c *types.Channel) []byte {
	ret := cdc.GetProtocMarshal().MustMarshalBinaryBare(c)
	return ret
}


func MarshalChannel(cdc *codec.MarshalProxy, c *types.Channel) ([]byte,error) {
	return cdc.GetProtocMarshal().MarshalBinaryBare(c)
}

func MustUnmarshalChannel(cdc *codec.MarshalProxy, data []byte) *types.Channel {
	var ret types.Channel
	cdc.GetProtocMarshal().MustUnmarshalBinaryBare(data, &ret)
	return &ret
}

func MustUnmarshalConnection(cdc *codec.MarshalProxy, data []byte) *types2.ConnectionEnd {
	var ret types2.ConnectionEnd
	cdc.GetProtocMarshal().MustUnmarshalBinaryBare(data, &ret)
	return &ret
}
func UnmarshalConnection(cdc *codec.MarshalProxy, data []byte) (*types2.ConnectionEnd,error){
	var ret types2.ConnectionEnd
	err:=cdc.GetProtocMarshal().UnmarshalBinaryBare(data, &ret)
	return &ret,err
}

func MustMarshalConnection(cdc *codec.MarshalProxy,c *types2.ConnectionEnd)[]byte{
	ret:=cdc.GetProtocMarshal().MustMarshalBinaryBare(c)
	return ret
}
