package version

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

var consensusTestCases = []Consensus{
	{},
	{1234, 5678},
	{math.MaxUint64, math.MaxUint64},
}

var cdc = amino.NewCodec()

func TestConsensusAmino(t *testing.T) {
	for _, cons := range consensusTestCases {
		expectData, err := cdc.MarshalBinaryBare(cons)
		require.NoError(t, err)

		var expectValue Consensus
		err = cdc.UnmarshalBinaryBare(expectData, &expectValue)
		require.NoError(t, err)

		var actualValue Consensus
		err = actualValue.UnmarshalFromAmino(expectData)
		require.NoError(t, err)

		require.EqualValues(t, expectValue, actualValue)
		require.EqualValues(t, len(expectData), cons.AminoSize())
	}
}

func BenchmarkConsensusAminoUnmarshal(b *testing.B) {
	var testData = make([][]byte, len(consensusTestCases))
	for i, cons := range consensusTestCases {
		data, err := cdc.MarshalBinaryBare(cons)
		require.NoError(b, err)
		testData[i] = data
	}
	b.ResetTimer()

	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, data := range testData {
				var value Consensus
				err := cdc.UnmarshalBinaryBare(data, &value)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
	b.Run("unmarshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, data := range testData {
				var value Consensus
				err := value.UnmarshalFromAmino(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
