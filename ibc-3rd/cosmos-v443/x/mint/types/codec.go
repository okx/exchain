package types

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	cryptocodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
