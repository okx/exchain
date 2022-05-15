package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateIsUint64(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid type", args{"a"}, true},
		{"correct IBC timeout", args{DefaultIbcTimeout}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, validateUint64(tt.args.i) != nil)
		})
	}
}

func Test_validateIsBool(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid bool", args{"a"}, true},
		{"correct bool", args{true}, false},
		{"correct bool", args{false}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, validateBool(tt.args.i) != nil)
		})
	}
}
