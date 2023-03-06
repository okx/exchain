package types

import (
	"fmt"
	"testing"

	tmrand "github.com/okx/okbchain/libs/tendermint/libs/rand"

	"github.com/okx/okbchain/libs/tendermint/crypto/secp256k1"
	"github.com/stretchr/testify/require"
)

func TestValidatorProtoBuf(t *testing.T) {
	val, _ := RandValidator(true, 100)
	testCases := []struct {
		msg      string
		v1       *Validator
		expPass1 bool
		expPass2 bool
	}{
		{"success validator", val, true, true},
		{"failure empty", &Validator{}, false, false},
		{"failure nil", nil, false, false},
	}
	for _, tc := range testCases {
		protoVal, err := tc.v1.ToProto()

		if tc.expPass1 {
			require.NoError(t, err, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}

		val, err := ValidatorFromProto(protoVal)
		if tc.expPass2 {
			require.NoError(t, err, tc.msg)
			require.Equal(t, tc.v1, val, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func TestPubKeyFromProto(t *testing.T) {

	r := func(ed25519 bool) *Validator {
		privVal := MockPV{secp256k1.GenPrivKey(), false, false}
		if ed25519 {
			privVal = NewMockPV()
		}
		votePower := int64(1)
		votePower += int64(tmrand.Uint32())
		pubKey, err := privVal.GetPubKey()
		if err != nil {
			panic(fmt.Errorf("could not retrieve pubkey %w", err))
		}
		val := NewValidator(pubKey, votePower)
		return val
	}

	type testCase struct {
		desc    string
		success bool
		f       func()
		v       *Validator
	}
	cases := []testCase{
		{
			desc:    "ed25519 pubkey",
			success: true,
			f:       func() {},
			v:       r(true),
		},
		{
			desc:    "secp256k1",
			success: true,
			f:       func() {},
			v:       r(false),
		},
	}

	for _, c := range cases {
		c.f()
		protoVal, err := c.v.ToProto()
		require.NoError(t, err)
		_, err = ValidatorFromProto(protoVal)
		if c.success {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
