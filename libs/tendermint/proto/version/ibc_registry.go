package version

import (
	protogogo "github.com/gogo/protobuf/proto"
)

func init() {
	protogogo.RegisterType((*Consensus)(nil), "tendermint.proto.version.Consensus")
}
