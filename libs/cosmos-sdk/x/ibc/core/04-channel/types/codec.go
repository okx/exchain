package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
)

var SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
