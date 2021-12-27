package types

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestDelta(t *testing.T) {
	d := &Deltas{
		Payload: DeltaPayload{
			[]byte("abci"),
			[]byte("detal"),
			[]byte("watch"),
		},
		Height: 1,
		Version: 1,
		CompressType: 2,
	}

	marshaled, err := d.Marshal()
	require.NoError(t, err)

	unmarshaled := &Deltas{}
	err = unmarshaled.Unmarshal(marshaled)
	require.NoError(t, err)

	assert.True(t,	bytes.Compare(unmarshaled.ABCIRsp(), d.ABCIRsp()) == 0)
	assert.True(t,	bytes.Compare(unmarshaled.DeltasBytes(), d.DeltasBytes()) == 0)
	assert.True(t,	bytes.Compare(unmarshaled.WatchBytes(), d.WatchBytes()) == 0)
	assert.True(t, 	unmarshaled.Height == d.Height)
	assert.True(t, 	unmarshaled.Version == d.Version)
	assert.True(t, 	unmarshaled.CompressType == d.CompressType)
}

func TestDeltas_Marshal(t *testing.T) {
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
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deltas{
				Height:          tt.fields.Height,
				Version:         tt.fields.Version,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}