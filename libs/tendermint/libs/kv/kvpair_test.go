package kv

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func TestKvPairAmino(t *testing.T) {
	var pairs = []Pair{
		{},
		{Key: []byte("key")},
		{Value: []byte("value")},
		{Key: []byte("key1"), Value: []byte("value1")},
		{Key: []byte("key1"), Value: []byte("value1"), XXX_NoUnkeyedLiteral: struct{}{}, XXX_sizecache: -10, XXX_unrecognized: []byte("unrecognized")},
		{Key: []byte{}, Value: []byte{}},
		{Key: []byte{}, Value: []byte{}, XXX_sizecache: 10},
	}

	for _, pair := range pairs {
		expect, err := cdc.MarshalBinaryBare(pair)
		require.NoError(t, err)

		actual, err := MarshalPairToAmino(pair)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)
	}
}

func BenchmarkKvPairAmino(b *testing.B) {
	var pair = Pair{
		Key:   []byte("key"),
		Value: []byte("value"),
	}
	b.ResetTimer()
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := cdc.MarshalBinaryBare(pair)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := MarshalPairToAmino(pair)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
