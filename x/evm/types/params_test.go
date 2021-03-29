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
			NewParams("ara", true, true, false, false, DefaultMaxGasLimitPerTx, 2929, 1884, 1344),
			false,
		},
		{
			"empty",
			Params{},
			true,
		},
		{
			"invalid evm denom",
			Params{
				EvmDenom: "@!#!@$!@5^32",
			},
			true,
		},
		{
			"invalid eip",
			Params{
				EvmDenom:  "stake",
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
	require.Error(t, validateEVMDenom(false))
	require.NoError(t, validateEVMDenom("aphoton"))
	require.Error(t, validateBool(""))
	require.NoError(t, validateBool(true))
	require.Error(t, validateEIPs(""))
	require.NoError(t, validateEIPs([]int{1884}))
}

func TestParams_String(t *testing.T) {
	const expectedParamsStr = `evm_denom: okt
enable_create: false
enable_call: false
extra_eips: []
enable_contract_deployment_whitelist: false
enable_contract_blocked_list: false
max_gas_limit: 30000000
`
	require.True(t, strings.EqualFold(expectedParamsStr, DefaultParams().String()))
}
