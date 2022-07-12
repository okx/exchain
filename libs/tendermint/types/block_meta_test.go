package types

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"math"
	"testing"
)

func TestBlockMetaValidateBasic(t *testing.T) {
	// TODO
}

func TestBlockMetaAmino(t *testing.T) {
	var testCases = []BlockMeta{
		{},
		{
			BlockID:   BlockID{Hash: []byte("hash"), PartsHeader: PartSetHeader{Total: 1, Hash: []byte("hash")}},
			BlockSize: 123,
			Header: Header{
				ChainID: "chainID",
			},
			NumTxs: -123,
		},
		{
			BlockSize: math.MinInt,
			NumTxs:    math.MinInt,
		},
		{
			BlockSize: math.MaxInt,
			NumTxs:    math.MaxInt,
		},
	}

	for _, tc := range testCases {
		expectBz, err := cdc.MarshalBinaryBare(tc)
		require.NoError(t, err)

		var expectValue BlockMeta
		err = cdc.UnmarshalBinaryBare(expectBz, &expectValue)
		require.NoError(t, err)

		var actualValue BlockMeta
		err = actualValue.UnmarshalFromAmino(cdc, expectBz)
		require.NoError(t, err)

		require.Equal(t, expectValue, actualValue)
	}

	for _, tc := range [][]byte{
		{4<<3 | byte(amino.Typ3_Varint), 0, 2<<3 | byte(amino.Typ3_Varint), 0},
		{4<<3 | byte(amino.Typ3_Varint), 0, 0},
		{4<<3 | byte(amino.Typ3_ByteLength), 0},
		{0},
	} {
		var expectValue BlockMeta
		err := cdc.UnmarshalBinaryBare(tc, &expectValue)
		require.Error(t, err)

		var actualValue BlockMeta
		err = actualValue.UnmarshalFromAmino(cdc, tc)
		require.Error(t, err)
	}

	{
		meta := BlockMeta{
			BlockID:   BlockID{Hash: []byte("hash"), PartsHeader: PartSetHeader{Total: 1, Hash: []byte("hash")}},
			BlockSize: 123,
			Header:    Header{},
			NumTxs:    -123,
		}
		meta2 := meta
		bz, _ := cdc.MarshalBinaryBare(BlockMeta{})
		err := cdc.UnmarshalBinaryBare(bz, &meta)
		require.NoError(t, err)

		err = meta2.UnmarshalFromAmino(cdc, bz)
		require.NoError(t, err)

		require.Equal(t, meta, meta2)
		require.Equal(t, BlockMeta{}, meta)
	}
}
