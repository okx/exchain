package keys

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	cryptocodec "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
)

// TODO: remove this file https://github.com/okex/exchain/ibc-3rd/cosmos-v443/issues/8047

// KeysCdc defines codec to be used with key operations
var KeysCdc *codec.LegacyAmino

func init() {
	KeysCdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(KeysCdc)
	KeysCdc.Seal()
}

// marshal keys
func MarshalJSON(o interface{}) ([]byte, error) {
	return KeysCdc.MarshalJSON(o)
}

// unmarshal json
func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return KeysCdc.UnmarshalJSON(bz, ptr)
}
