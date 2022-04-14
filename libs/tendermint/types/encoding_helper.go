package types

import (
	//gogotypes "github.com/gogo/protobuf/types"
	//"github.com/okex/exchain/libs/tendermint/libs/bytes"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"
)

// cdcEncode returns nil if the input is nil, otherwise returns
// cdc.MustMarshalBinaryBare(item)
func cdcEncode(item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		return cdc.MustMarshalBinaryBare(item)
	}
	return nil
}

func ibccdcEncode(item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		switch item := item.(type) {
		case string:
			i := gogotypes.StringValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case int64:
			i := gogotypes.Int64Value{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case bytes.HexBytes:
			i := gogotypes.BytesValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		default:
			return nil
		}
	}

	return nil
}
