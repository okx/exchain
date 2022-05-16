package kv

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

var testPairs = []Pair{
	{},
	{Key: []byte("key")},
	{Value: []byte("value")},
	{Key: []byte("key1"), Value: []byte("value1")},
	{Key: []byte("key1"), Value: []byte("value1"), XXX_NoUnkeyedLiteral: struct{}{}, XXX_sizecache: -10, XXX_unrecognized: []byte("unrecognized")},
	{Key: []byte{}, Value: []byte{}},
	{Key: []byte{}, Value: []byte{}, XXX_sizecache: 10},
}

func TestKvPairAmino(t *testing.T) {
	for _, pair := range testPairs {
		expect, err := cdc.MarshalBinaryBare(pair)
		require.NoError(t, err)

		actual, err := pair.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.Equal(t, len(expect), pair.AminoSize(cdc))

		var pair2 Pair
		err = cdc.UnmarshalBinaryBare(expect, &pair2)
		require.NoError(t, err)
		var pair3 Pair
		err = pair3.UnmarshalFromAmino(cdc, expect)
		require.NoError(t, err)

		require.EqualValues(t, pair2, pair3)
	}
}

func BenchmarkKvPairAminoMarshal(b *testing.B) {
	b.ResetTimer()
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, pair := range testPairs {
				_, err := cdc.MarshalBinaryBare(pair)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, pair := range testPairs {
				_, err := pair.MarshalToAmino(cdc)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkKvPairAminoUnmarshal(b *testing.B) {
	testData := make([][]byte, len(testPairs))
	for i, pair := range testPairs {
		data, err := cdc.MarshalBinaryBare(pair)
		if err != nil {
			b.Fatal(err)
		}
		testData[i] = data
	}
	b.ResetTimer()
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, data := range testData {
				var pair Pair
				err := cdc.UnmarshalBinaryBare(data, &pair)
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
				var pair Pair
				err := pair.UnmarshalFromAmino(cdc, data)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
