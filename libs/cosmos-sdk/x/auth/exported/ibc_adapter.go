package exported

import "github.com/gogo/protobuf/proto"

type AccountAdapter interface {
	Account
	proto.Message
}
