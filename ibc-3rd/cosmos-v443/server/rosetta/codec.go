package rosetta

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	codectypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/codec/types"
	cryptocodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
	authcodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/types"
	bankcodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/bank/types"
)

// MakeCodec generates the codec required to interact
// with the cosmos APIs used by the rosetta gateway
func MakeCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)

	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptocodec.RegisterInterfaces(ir)

	return cdc, ir
}
