package consensus

import (
	"github.com/okex/exchain/libs/tendermint/consensus/types"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"
)

func TestTimeoutInfoAmino(t *testing.T) {
	testCases := []timeoutInfo{
		{},
		{
			Duration: 1 * time.Hour,
			Height:   2,
			Round:    123,
			Step:     types.RoundStepNewHeight,
		},
		{
			Duration: math.MaxInt64,
			Height:   math.MaxInt64,
			Round:    math.MaxInt,
			Step:     types.RoundStepPrecommit,
		},
		{
			Duration: math.MinInt64,
			Height:   math.MinInt64,
			Round:    math.MinInt,
			Step:     types.RoundStepNewRound,
		},
	}
	for _, tc := range testCases {
		expectData := cdc.MustMarshalBinaryBare(&tc)
		actualData, err := cdc.MarshalBinaryWithSizer(&tc, false)
		require.NoError(t, err)
		require.Equal(t, expectData, actualData)
		require.Equal(t, len(expectData), tc.AminoSize(cdc)+4)
	}
}

type testMsg struct{}

func (testMsg) ValidateBasic() error { return nil }

func TestMsgInfoAmino(t *testing.T) {
	cdc.RegisterConcrete(testMsg{}, "consensus/testMsg", nil)
	testCases := []msgInfo{
		{},
		{
			Msg:    &ProposalPOLMessage{},
			PeerID: "test",
		},
		{
			Msg: testMsg{},
		},
	}
	for _, tc := range testCases {
		expectData := cdc.MustMarshalBinaryBare(&tc)
		actualData, err := cdc.MarshalBinaryWithSizer(&tc, false)
		require.NoError(t, err)
		require.Equal(t, expectData, actualData)
		require.Equal(t, len(expectData), 4+tc.AminoSize(cdc))
	}
}
