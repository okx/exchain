package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/types/kv"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp"
	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	dec := simulation.NewDecodeStore(app.TransferKeeper)

	trace := types.DenomTrace{
		BaseDenom: "uatom",
		Path:      "transfer/channelToA",
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   types.PortKey,
				Value: []byte(types.PortID),
			},
			{
				Key:   types.DenomTraceKey,
				Value: app.TransferKeeper.MustMarshalDenomTrace(trace),
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"PortID", fmt.Sprintf("Port A: %s\nPort B: %s", types.PortID, types.PortID)},
		{"DenomTrace", fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s", trace.IBCDenom(), trace.IBCDenom())},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				//require.Panics(t, func() { dec(nil, kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
				kvA := tmkv.Pair{
					Key:   kvPairs.Pairs[i].GetKey(),
					Value: kvPairs.Pairs[i].GetValue(),
				}
				require.Panics(t, func() { dec(nil, kvA, kvA) }, tt.name)
			} else {
				kvA := tmkv.Pair{
					Key:   kvPairs.Pairs[i].GetKey(),
					Value: kvPairs.Pairs[i].GetValue(),
				}
				//require.Equal(t, tt.expectedLog, dec(nil, kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
				require.Equal(t, tt.expectedLog, dec(nil, kvA, kvA), tt.name)
			}
		})
	}
}
