package types

import (
	"testing"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	strExpected = `Distribution Params:
  Community Tax:          0.020000000000000000
  Withdraw Addr Enabled:  true`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := defaultState.Params
	require.Equal(t, sdk.NewDecWithPrec(2, 2), defaultParams.CommunityTax)
	require.Equal(t, true, defaultParams.WithdrawAddrEnabled)

	require.Equal(t, strExpected, defaultParams.String())
	yamlStr, err := defaultParams.MarshalYAML()
	require.NoError(t, err)
	require.Equal(t, strExpected, yamlStr)
}

func Test_validateAuxFuncs(t *testing.T) {
	type args struct {
		i interface{}
	}

	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"wrong type", args{10.5}, true},
		{"nil Int pointer", args{sdk.Dec{}}, true},
		{"negative", args{sdk.NewDec(-1)}, true},
		{"one dec", args{sdk.NewDec(1)}, false},
		{"two dec", args{sdk.NewDec(2)}, true},
	}

	for _, tc := range testCases {
		stc := tc

		t.Run(stc.name, func(t *testing.T) {
			require.Equal(t, stc.wantErr, validateCommunityTax(stc.args.i) != nil)
		})
	}
}