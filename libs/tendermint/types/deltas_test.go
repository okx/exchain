package types

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Height:       1,
		Version:      1,
		CompressType: 2,
	}

	marshaled, err := d.Marshal()
	require.NoError(t, err)

	unmarshaled := &Deltas{}
	err = unmarshaled.Unmarshal(marshaled)
	require.NoError(t, err)

	assert.True(t, bytes.Compare(unmarshaled.ABCIRsp(), d.ABCIRsp()) == 0)
	assert.True(t, bytes.Compare(unmarshaled.DeltasBytes(), d.DeltasBytes()) == 0)
	assert.True(t, bytes.Compare(unmarshaled.WatchBytes(), d.WatchBytes()) == 0)
	assert.True(t, unmarshaled.Height == d.Height)
	assert.True(t, unmarshaled.Version == d.Version)
	assert.True(t, unmarshaled.CompressType == d.CompressType)
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
		{"no watchData", fields{Height: 1, Version: DeltaVersion, Payload: noWD}, args{1}, false},
		{"wrong height", fields{Height: 1, Version: DeltaVersion, Payload: payload}, args{2}, false},
		{"low version", fields{Height: 1, Version: DeltaVersion - 1, Payload: payload}, args{1}, false},
		{"high version", fields{Height: 1, Version: DeltaVersion + 1, Payload: payload}, args{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dds := &Deltas{
				Height:          tt.fields.Height,
				Version:         tt.fields.Version,
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

func TestGetDeltaVersionNone(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Flags().Int(FlagDeltaVersion, DeltaVersion, "Specify delta version")
	viper.BindPFlag(FlagDeltaVersion, cmd.Flags().Lookup(FlagDeltaVersion))
	cmd.Execute()

	tests := []struct {
		name string
		want int
	}{
		{name: "1. user do not set --delta-version", want: DeltaVersion},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetDeltaVersion(), "GetDeltaVersion()")
		})
	}
}

func TestGetDeltaVersionUserSpecified(t *testing.T) {
	userSpecified := DeltaVersion + 1
	cmd := cobra.Command{}
	cmd.Flags().Int(FlagDeltaVersion, DeltaVersion, "Specify delta version")
	viper.BindPFlag(FlagDeltaVersion, cmd.Flags().Lookup(FlagDeltaVersion))
	cmd.ParseFlags([]string{fmt.Sprintf("--%v=%d", FlagDeltaVersion, userSpecified)})
	cmd.Execute()

	tests := []struct {
		name string
		want int
	}{
		{name: "1. user specify --delta-version", want: userSpecified},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetDeltaVersion(), "GetDeltaVersion()")
		})
	}
}
