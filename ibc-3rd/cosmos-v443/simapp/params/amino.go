//go:build test_amino
// +build test_amino

package params

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/legacy/legacytx"
)

// MakeTestEncodingConfig creates an EncodingConfig for an amino based test configuration.
// This function should be used only internally (in the SDK).
// App user should'nt create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() EncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewAminoCodec(cdc)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          legacytx.StdTxConfig{Cdc: cdc},
		Amino:             cdc,
	}
}
