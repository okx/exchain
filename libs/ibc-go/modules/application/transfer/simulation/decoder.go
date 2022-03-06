package simulation

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/ibc-go/modules/application/transfer/types"
	"github.com/okex/exchain/libs/tendermint/libs/kv"
)

// TransferUnmarshaler defines the expected encoding store functions.
type TransferUnmarshaler interface {
	MustUnmarshalDenomTrace([]byte) types.DenomTrace
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding DenomTrace type.
func NewDecodeStore(kCdc TransferUnmarshaler) func(cdc *codec.Codec,kvA, kvB kv.Pair) string {
	return func(cdc *codec.Codec,kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PortKey):
			return fmt.Sprintf("Port A: %s\nPort B: %s", string(kvA.Value), string(kvB.Value))

		case bytes.Equal(kvA.Key[:1], types.DenomTraceKey):
			denomTraceA := kCdc.MustUnmarshalDenomTrace(kvA.Value)
			denomTraceB := kCdc.MustUnmarshalDenomTrace(kvB.Value)
			return fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s", denomTraceA.IBCDenom(), denomTraceB.IBCDenom())

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
