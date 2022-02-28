package simulation

import (
	"bytes"
	"fmt"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/types/kv"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/evidence/exported"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/evidence/types"
)

type EvidenceUnmarshaler interface {
	UnmarshalEvidence([]byte) (exported.Evidence, error)
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding evidence type.
func NewDecodeStore(cdc EvidenceUnmarshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.KeyPrefixEvidence):
			evidenceA, err := cdc.UnmarshalEvidence(kvA.Value)
			if err != nil {
				panic(fmt.Sprintf("cannot unmarshal evidence: %s", err.Error()))
			}

			evidenceB, err := cdc.UnmarshalEvidence(kvB.Value)
			if err != nil {
				panic(fmt.Sprintf("cannot unmarshal evidence: %s", err.Error()))
			}

			return fmt.Sprintf("%v\n%v", evidenceA, evidenceB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
