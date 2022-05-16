package types

import "github.com/gogo/protobuf/proto"

func (x BlockIDFlag) String() string {
	return proto.EnumName(BlockIDFlag_name, int32(x))
}
