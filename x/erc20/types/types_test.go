package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsValidIBCDenom(t *testing.T) {
	tests := []struct {
		name    string
		denom   string
		success bool
	}{
		{"wrong length", "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD", false},
		{"invalid denom", "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865", false},
		{"correct IBC denom", IbcDenomDefaultValue, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.success, IsValidIBCDenom(tt.denom))
		})
	}
}
