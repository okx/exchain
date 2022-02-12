package common

import "github.com/okex/exchain/libs/cosmos-sdk/codec"

func DefaultMarshal(c codec.Codec,data interface{})([]byte,error){
	return c.MarshalBinaryBare(data)
}

