package simulation

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/tendermint/libs/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding capability type.
func NewDecodeStore() func(cdc *codec.Codec,kvA, kvB kv.Pair) string {
	return func(cdc *codec.Codec,kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyIndex):
			idxA := sdk.BigEndianToUint64(kvA.Value)
			idxB := sdk.BigEndianToUint64(kvB.Value)
			return fmt.Sprintf("Index A: %d\nIndex B: %d\n", idxA, idxB)

		case bytes.HasPrefix(kvA.Key, types.KeyPrefixIndexCapability):
			var capOwnersA, capOwnersB types.CapabilityOwners
			cdc.MustUnmarshalBinaryBare(kvA.Value, &capOwnersA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &capOwnersB)
			return fmt.Sprintf("CapabilityOwners A: %v\nCapabilityOwners B: %v\n", capOwnersA, capOwnersB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X (%s)", types.ModuleName, kvA.Key, string(kvA.Key)))
		}
	}
}
