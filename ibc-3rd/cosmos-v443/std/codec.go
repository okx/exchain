package std

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec/types"
	cryptocodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	txtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/tx"
)

// RegisterLegacyAminoCodec registers types with the Amino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers Interfaces from sdk/types, vesting, crypto, tx.
func RegisterInterfaces(interfaceRegistry types.InterfaceRegistry) {
	sdk.RegisterInterfaces(interfaceRegistry)
	txtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
}
