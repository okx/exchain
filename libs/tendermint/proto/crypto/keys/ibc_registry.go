package crypto

import (
	protogogo "github.com/gogo/protobuf/proto"
)

func init() {
	protogogo.RegisterType((*PublicKey)(nil), "tendermint.proto.crypto.keys.PublicKey")
}
