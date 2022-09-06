package types

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelta(t *testing.T) {
	d := &Deltas{
		Payload: DeltaPayload{
			ABCIRsp:        []byte("abci"),
			DeltasBytes:    []byte("detal"),
			WatchBytes:     []byte("watch"),
			WasmWatchBytes: []byte("wasm watch"),
		},
		Height:       1,
		CompressType: 2,
	}

	marshaled, err := d.Marshal()
	require.NoError(t, err)

	unmarshaled := &Deltas{}
	err = unmarshaled.Unmarshal(marshaled)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(unmarshaled.ABCIRsp(), d.ABCIRsp()), "ABCIRsp not equal")
	assert.True(t, bytes.Equal(unmarshaled.DeltasBytes(), d.DeltasBytes()), "DeltasBytes not equal")
	assert.True(t, bytes.Equal(unmarshaled.WatchBytes(), d.WatchBytes()), "WatchBytes not equal")
	assert.True(t, bytes.Equal(unmarshaled.WasmWatchBytes(), d.WasmWatchBytes()), "WasmWatchBytes not equal")
	assert.True(t, unmarshaled.Height == d.Height, "Height not equal")
	assert.True(t, unmarshaled.CompressType == d.CompressType, "CompressType not equal")
}

func TestDeltas_MarshalUnMarshal(t *testing.T) {
	payload := DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: []byte("DeltasBytes"), WatchBytes: []byte("WatchBytes")}
	type fields struct {
		Height          int64
		Version         int
		Payload         DeltaPayload
		CompressType    int
		CompressFlag    int
		marshalElapsed  time.Duration
		compressElapsed time.Duration
		hashElapsed     time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"no compress", fields{Height: 1, Version: 1, Payload: payload}, false},
		{"compress", fields{Height: 1, Version: 1, Payload: payload, CompressType: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deltas{
				Height:          tt.fields.Height,
				Payload:         tt.fields.Payload,
				CompressType:    tt.fields.CompressType,
				CompressFlag:    tt.fields.CompressFlag,
				marshalElapsed:  tt.fields.marshalElapsed,
				compressElapsed: tt.fields.compressElapsed,
				hashElapsed:     tt.fields.hashElapsed,
			}
			got, err := d.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = d.Unmarshal(got)
			assert.Nil(t, err)
		})
	}
}

func TestDeltas_Validate(t *testing.T) {
	payload := DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: []byte("DeltasBytes"), WatchBytes: []byte("WatchBytes")}
	noABCIRsp := DeltaPayload{ABCIRsp: nil, DeltasBytes: []byte("DeltasBytes"), WatchBytes: []byte("WatchBytes")}
	noDeltaBytes := DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: nil, WatchBytes: []byte("WatchBytes")}
	noWD := DeltaPayload{ABCIRsp: []byte("ABCIRsp"), DeltasBytes: []byte("DeltasBytes"), WatchBytes: nil}

	type fields struct {
		Height          int64
		Version         int
		Payload         DeltaPayload
		CompressType    int
		CompressFlag    int
		marshalElapsed  time.Duration
		compressElapsed time.Duration
		hashElapsed     time.Duration
	}
	type args struct {
		height int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"normal case", fields{Height: 1, Version: DeltaVersion, Payload: payload}, args{1}, true},
		{"no ABCIRsp", fields{Height: 1, Version: DeltaVersion, Payload: noABCIRsp}, args{1}, false},
		{"no deltaBytes", fields{Height: 1, Version: DeltaVersion, Payload: noDeltaBytes}, args{1}, false},
		{"no watchData", fields{Height: 1, Version: DeltaVersion, Payload: noWD}, args{1}, !FastQuery},
		{"wrong height", fields{Height: 1, Version: DeltaVersion, Payload: payload}, args{2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dds := &Deltas{
				Height:          tt.fields.Height,
				Payload:         tt.fields.Payload,
				CompressType:    tt.fields.CompressType,
				CompressFlag:    tt.fields.CompressFlag,
				marshalElapsed:  tt.fields.marshalElapsed,
				compressElapsed: tt.fields.compressElapsed,
				hashElapsed:     tt.fields.hashElapsed,
			}
			if got := dds.Validate(tt.args.height); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeltaPayloadAmino(t *testing.T) {
	var testCases = []DeltaPayload{
		{},
		{
			ABCIRsp:        []byte("ABCIResp"),
			DeltasBytes:    []byte("DeltasBytes"),
			WatchBytes:     []byte("WatchBytes"),
			WasmWatchBytes: []byte("WasmWatchBytes"),
		},
		{
			ABCIRsp:        []byte{},
			DeltasBytes:    []byte{},
			WatchBytes:     []byte{},
			WasmWatchBytes: []byte{},
		},
	}
	for _, testCase := range testCases {
		expectBz, err := cdc.MarshalBinaryBare(testCase)
		require.NoError(t, err)

		actulaBz, err := testCase.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.Equal(t, expectBz, actulaBz)
		require.Equal(t, len(expectBz), testCase.AminoSize(cdc))

		var expectValue DeltaPayload
		err = cdc.UnmarshalBinaryBare(expectBz, &expectValue)
		require.NoError(t, err)

		var actualValue DeltaPayload
		err = actualValue.UnmarshalFromAmino(cdc, expectBz)
		require.NoError(t, err)

		require.Equal(t, expectValue, actualValue)
	}

	{
		bz := []byte{1<<3 | 2, 0, 2<<3 | 2, 0, 3<<3 | 2, 0, 4<<3 | 2, 0}
		var expectValue DeltaPayload
		err := cdc.UnmarshalBinaryBare(bz, &expectValue)
		require.NoError(t, err)

		var actualValue DeltaPayload
		err = actualValue.UnmarshalFromAmino(cdc, bz)
		require.NoError(t, err)

		require.Equal(t, expectValue, actualValue)
	}
}

func TestDeltasMessageAmino(t *testing.T) {
	var testCases = []DeltasMessage{
		{},
		{
			Metadata:     []byte("Metadata"),
			MetadataHash: []byte("MetadataHash"),
			Height:       12345,
			CompressType: 1234,
			From:         "from",
		},
		{
			Metadata:     []byte{},
			MetadataHash: []byte{},
		},
		{
			Height:       math.MaxInt64,
			CompressType: math.MaxInt,
		},
		{
			Height:       math.MinInt64,
			CompressType: math.MinInt,
		},
	}
	utInitValue := DeltasMessage{
		Metadata:     []byte("Metadata"),
		MetadataHash: []byte("MetadataHash"),
		Height:       12345,
		CompressType: 1234,
		From:         "from",
	}

	for _, testCase := range testCases {
		expectBz, err := cdc.MarshalBinaryBare(testCase)
		require.NoError(t, err)

		actulaBz, err := testCase.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.Equal(t, expectBz, actulaBz)
		require.Equal(t, len(expectBz), testCase.AminoSize(cdc))

		var expectValue = utInitValue
		err = cdc.UnmarshalBinaryBare(expectBz, &expectValue)
		require.NoError(t, err)

		var actualValue = utInitValue
		err = actualValue.UnmarshalFromAmino(cdc, expectBz)
		require.NoError(t, err)

		require.Equal(t, expectValue, actualValue)
	}

	{
		bz := []byte{1<<3 | 2, 0, 2<<3 | 2, 0, 3 << 3, 0, 4 << 3, 0, 5<<3 | 2, 0}
		var expectValue = utInitValue
		err := cdc.UnmarshalBinaryBare(bz, &expectValue)
		require.NoError(t, err)

		var actualValue = utInitValue
		err = actualValue.UnmarshalFromAmino(cdc, bz)
		require.NoError(t, err)

		require.Equal(t, expectValue, actualValue)
	}
}
