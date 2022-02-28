package legacytx

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil)
}
