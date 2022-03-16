package types

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_validateIbcDenomParam(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid type", args{sdk.OneInt()}, true},

		{"wrong length", args{"ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD"}, true},
		{"invalid denom", args{"aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865"}, true},
		{"correct IBC denom", args{IbcDenomDefaultValue}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, validateIbcDenom(tt.args.i) != nil)
		})
	}
}

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
		{"correct IBC timeout", args{IbcTimeoutDefaultValue}, false},
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
