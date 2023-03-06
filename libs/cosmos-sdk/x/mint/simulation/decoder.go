package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/okx/okbchain/libs/tendermint/libs/kv"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding mint type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key, types.MinterKey):
		var minterA, minterB types.Minter
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &minterA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &minterB)
		return fmt.Sprintf("%v\n%v", minterA, minterB)
	default:
		panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
	}
}
