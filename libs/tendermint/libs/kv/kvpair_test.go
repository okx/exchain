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
		{Key: []byte("key1"), Value: []byte("value1")},
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
