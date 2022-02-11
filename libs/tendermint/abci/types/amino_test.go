package types

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/kv"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

var eventTestcases = []Event{
	{},
	{
		Type: "test",
	},
	{
		Attributes: []kv.Pair{
			{Key: []byte("key"), Value: []byte("value")},
			{Key: []byte("key2"), Value: []byte("value2")},
		},
	},
	{
		Type: "test",
		Attributes: []kv.Pair{
			{Key: []byte("key"), Value: []byte("value")},
			{Key: []byte("key2"), Value: []byte("value2")},
			{Key: []byte("key3"), Value: []byte("value3")},
			{},
		},
	},
	{
		Attributes: []kv.Pair{},
	},
}

func TestEventAmino(t *testing.T) {
	for _, event := range eventTestcases {
		expect, err := cdc.MarshalBinaryBare(event)
		require.NoError(t, err)

		actual, err := event.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.Equal(t, len(expect), event.AminoSize(cdc))

		var value Event
		err = cdc.UnmarshalBinaryBare(expect, &value)
		require.NoError(t, err)

		var value2 Event
		err = value2.UnmarshalFromAmino(cdc, expect)
		require.NoError(t, err)

		require.EqualValues(t, value, value2)
	}
}

func BenchmarkEventAminoMarshal(b *testing.B) {
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, event := range eventTestcases {
				_, err := cdc.MarshalBinaryBare(event)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, event := range eventTestcases {
				_, err := event.MarshalToAmino(cdc)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEventAminoUnmarshal(b *testing.B) {
	testData := make([][]byte, len(eventTestcases))
	for i, event := range eventTestcases {
		data, err := cdc.MarshalBinaryBare(event)
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
				var event Event
				err := cdc.UnmarshalBinaryBare(data, &event)
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
				var event Event
				err := event.UnmarshalFromAmino(cdc, data)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func TestPubKeyAmino(t *testing.T) {
	var pubkeys = []PubKey{
		{},
		{Type: "type"},
		{Data: []byte("testdata")},
		{
			Type: "test",
			Data: []byte("data"),
		},
		{
			Data: []byte{},
		},
	}

	for _, pubkey := range pubkeys {
		expect, err := cdc.MarshalBinaryBare(pubkey)
		require.NoError(t, err)

		actual, err := pubkey.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), pubkey.AminoSize(cdc))
	}
}

func TestValidatorUpdateAmino(t *testing.T) {
	var validatorUpdates = []ValidatorUpdate{
		{},
		{
			PubKey: PubKey{
				Type: "test",
				Data: []byte{},
			},
		},
		{
			PubKey: PubKey{
				Type: "test",
				Data: []byte("data"),
			},
		},
		{
			Power: math.MaxInt64,
		},
		{
			Power: math.MinInt64,
		},
		{
			PubKey: PubKey{
				Type: "",
				Data: []byte("data"),
			},
			Power: 100,
		},
	}

	for _, validatorUpdate := range validatorUpdates {
		expect, err := cdc.MarshalBinaryBare(validatorUpdate)
		require.NoError(t, err)

		actual, err := validatorUpdate.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), validatorUpdate.AminoSize(cdc))
	}
}

func TestBlockParamsAmino(t *testing.T) {
	tests := []BlockParams{
		{},
		{
			MaxBytes: 100,
			MaxGas:   200,
		},
		{
			MaxBytes: -100,
			MaxGas:   -200,
		},
		{
			MaxBytes: math.MaxInt64,
			MaxGas:   math.MaxInt64,
		},
		{
			MaxBytes: math.MinInt64,
			MaxGas:   math.MinInt64,
		},
	}

	for _, test := range tests {
		expect, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

		actual, err := test.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), test.AminoSize(cdc))
	}
}

func TestEvidenceParamsAmino(t *testing.T) {
	tests := []EvidenceParams{
		{},
		{
			MaxAgeNumBlocks: 100,
			MaxAgeDuration:  1000 * time.Second,
		},
		{
			MaxAgeNumBlocks: -100,
			MaxAgeDuration:  time.Hour * 24 * 365,
		},
		{
			MaxAgeNumBlocks: math.MinInt64,
			MaxAgeDuration:  math.MinInt64,
		},
		{
			MaxAgeNumBlocks: math.MaxInt64,
			MaxAgeDuration:  math.MaxInt64,
		},
	}

	for _, test := range tests {
		expect, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

		actual, err := test.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), test.AminoSize(cdc))
	}
}

func TestValidatorParamsAmino(t *testing.T) {
	tests := []ValidatorParams{
		{},
		{
			PubKeyTypes: []string{},
		},
		{
			PubKeyTypes: []string{""},
		},
		{
			PubKeyTypes: []string{"ed25519"},
		},
		{
			PubKeyTypes: []string{"ed25519", "ed25519", "", "rsa"},
		},
	}

	for _, test := range tests {
		expect, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

		actual, err := test.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), test.AminoSize(cdc))
	}
}

func TestConsensusParamsAmino(t *testing.T) {
	tests := []ConsensusParams{
		{
			Block:     &BlockParams{},
			Evidence:  &EvidenceParams{},
			Validator: &ValidatorParams{},
		},
		{
			Block: &BlockParams{
				MaxBytes: 100,
			},
			Evidence: &EvidenceParams{
				MaxAgeDuration: 5 * time.Second,
			},
			Validator: &ValidatorParams{
				PubKeyTypes: nil,
			},
		},
		{
			Validator: &ValidatorParams{
				PubKeyTypes: []string{"ed25519", "rsa"},
			},
		},
		{
			Block: &BlockParams{
				MaxBytes: 100,
				MaxGas:   200,
			},
			Evidence: &EvidenceParams{
				MaxAgeNumBlocks: 500,
				MaxAgeDuration:  6 * time.Second,
			},
			Validator: &ValidatorParams{
				PubKeyTypes: []string{},
			},
		},
	}

	for _, test := range tests {
		expect, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

		actual, err := test.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), test.AminoSize(cdc))
	}
}

var responseDeliverTxTestCases = []*ResponseDeliverTx{
	{},
	{123, nil, "", "", 0, 0, nil, "", struct{}{}, nil, 0},
	{Code: 123, Data: []byte("this is data"), Log: "log123", Info: "123info", GasWanted: 1234445, GasUsed: 98, Events: []Event{}, Codespace: "sssdasf"},
	{Code: math.MaxUint32, GasWanted: math.MaxInt64, GasUsed: math.MaxInt64},
	{Code: 0, GasWanted: -1, GasUsed: -1},
	{Code: 0, GasWanted: math.MinInt64, GasUsed: math.MinInt64},
	{Events: []Event{{}, {Type: "Event"}, {Type: "Event2"}}, Data: []byte{}},
}

func TestResponseDeliverTxAmino(t *testing.T) {
	for i, resp := range responseDeliverTxTestCases {
		expect, err := cdc.MarshalBinaryBare(resp)
		require.NoError(t, err)

		actual, err := resp.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.Equal(t, len(expect), resp.AminoSize(cdc))

		var resp1 ResponseDeliverTx
		err = cdc.UnmarshalBinaryBare(expect, &resp1)
		require.NoError(t, err)

		var resp2 ResponseDeliverTx
		err = resp2.UnmarshalFromAmino(cdc, expect)
		require.NoError(t, err, fmt.Sprintf("error case index %d", i))

		require.EqualValues(t, resp1, resp2)
	}
}

func BenchmarkResponseDeliverTxAminoMarshal(b *testing.B) {
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, resp := range responseDeliverTxTestCases {
				_, err := cdc.MarshalBinaryBare(resp)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, resp := range responseDeliverTxTestCases {
				_, err := resp.MarshalToAmino(cdc)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkResponseDeliverTxAminoUnmarshal(b *testing.B) {
	testData := make([][]byte, len(responseDeliverTxTestCases))
	for i, resp := range responseDeliverTxTestCases {
		data, err := cdc.MarshalBinaryBare(resp)
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
				var resp ResponseDeliverTx
				err := cdc.UnmarshalBinaryBare(data, &resp)
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
				var resp ResponseDeliverTx
				err := resp.UnmarshalFromAmino(cdc, data)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func TestResponseBeginBlockAmino(t *testing.T) {
	var resps = []*ResponseBeginBlock{
		{},
		{
			Events: []Event{
				{
					Type: "test",
				},
				{
					Type: "test2",
				},
			},
		},
		{
			Events: []Event{},
		},
		{
			Events: []Event{{}, {}, {}, {}},
		},
	}
	for _, resp := range resps {
		expect, err := cdc.MarshalBinaryBare(resp)
		require.NoError(t, err)

		actual, err := resp.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.Equal(t, len(expect), resp.AminoSize(cdc))
	}
}

func TestResponseEndBlockAmino(t *testing.T) {
	var resps = []*ResponseEndBlock{
		{},
		{
			ValidatorUpdates: []ValidatorUpdate{
				{
					PubKey: PubKey{
						Type: "test",
					},
				},
				{
					PubKey: PubKey{
						Type: "test2",
					},
				},
				{},
			},
			ConsensusParamUpdates: &ConsensusParams{},
			Events:                []Event{{}, {Type: "Event"}, {Type: "Event2"}},
		},
		{
			ValidatorUpdates:      []ValidatorUpdate{},
			ConsensusParamUpdates: &ConsensusParams{},
			Events:                []Event{},
		},
		{
			ValidatorUpdates:      []ValidatorUpdate{{}},
			ConsensusParamUpdates: &ConsensusParams{Block: &BlockParams{}, Evidence: &EvidenceParams{}, Validator: &ValidatorParams{}},
			Events:                []Event{{}},
		},
	}
	for _, resp := range resps {
		expect, err := cdc.MarshalBinaryBare(resp)
		require.NoError(t, err)

		actual, err := resp.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), resp.AminoSize(cdc))
	}
}
