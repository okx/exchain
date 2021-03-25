package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid",
			NewParams(true, true, 2929, 1884, 1344),
			false,
		},
		{
			"invalid eip",
			Params{
				ExtraEIPs: []int{1},
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestParamsValidatePriv(t *testing.T) {
	require.Error(t, validateBool(""))
	require.NoError(t, validateBool(true))
	require.Error(t, validateEIPs(""))
	require.NoError(t, validateEIPs([]int{1884}))
}

func TestParams_String(t *testing.T) {
	const expectedParamsStr = `enable_create: false
enable_call: false
extra_eips: []
`
	require.True(t, strings.EqualFold(expectedParamsStr, DefaultParams().String()))
}
