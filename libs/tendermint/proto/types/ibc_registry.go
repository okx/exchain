package types

import (
	protogogo "github.com/gogo/protobuf/proto"
)

func init() {
	protogogo.RegisterType((*SignedHeader)(nil), "tendermint.types.SignedHeader")

	protogogo.RegisterType((*PartSetHeader)(nil), "tendermint.proto.types.PartSetHeader")
	protogogo.RegisterType((*Part)(nil), "tendermint.proto.types.Part")
	protogogo.RegisterType((*BlockID)(nil), "tendermint.proto.types.BlockID")
	protogogo.RegisterType((*Header)(nil), "tendermint.proto.types.Header")
	protogogo.RegisterType((*Data)(nil), "tendermint.proto.types.Data")
	protogogo.RegisterType((*Vote)(nil), "tendermint.proto.types.Vote")
	protogogo.RegisterType((*Commit)(nil), "tendermint.proto.types.Commit")
	protogogo.RegisterType((*CommitSig)(nil), "tendermint.proto.types.CommitSig")
	protogogo.RegisterType((*Proposal)(nil), "tendermint.proto.types.Proposal")
	protogogo.RegisterType((*BlockMeta)(nil), "tendermint.proto.types.BlockMeta")

	protogogo.RegisterType((*ValidatorSet)(nil), "tendermint.proto.types.ValidatorSet")
	protogogo.RegisterType((*Validator)(nil), "tendermint.proto.types.Validator")
}
