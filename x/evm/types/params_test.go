package types

import (
	"encoding/json"
	"strings"
	"testing"

	gojson "github.com/goccy/go-json"
	"github.com/json-iterator/go"
	"github.com/pquerna/ffjson/ffjson"
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
			NewParams(true, true, false, false, defaultMaxGasLimitPerTx, 2929, 1884, 1344),
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
	require.NoError(t, validateUint64(uint64(30000000)))
	require.Error(t, validateUint64("test"))
}

func TestParams_String(t *testing.T) {
	const expectedParamsStr = `enable_create: true
enable_call: true
extra_eips: []
enable_contract_deployment_whitelist: false
enable_contract_blocked_list: true
max_gas_limit_per_tx: 50000000
`
	require.True(t, strings.EqualFold(expectedParamsStr, DefaultParams().String()))
}

func BenchmarkParamsUnmarshal(b *testing.B) {
	s := `{"enable_create":true,"enable_call":false,"extra_eips":[1,1,1,1],"enable_contract_deployment_whitelist":true,"enable_contract_blocked_list":true,"max_gas_limit_per_tx":100}`
	bz := []byte(s)
	b.ResetTimer()
	b.Run("json", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var params Params
			_ = json.Unmarshal(bz, &params)
		}
	})

	b.Run("jsoniter", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var params Params
			_ = jsoniter.Unmarshal(bz, &params)
		}
	})
	b.Run("ffjson", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var params Params
			_ = ffjson.Unmarshal(bz, &params)
		}
	})
	b.Run("go-json", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var params Params
			_ = gojson.Unmarshal(bz, &params)
		}
	})
}

func BenchmarkParamsMarshal(b *testing.B) {
	params := NewParams(true, false, true, true, 100)
	b.ResetTimer()
	b.Run("json", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = json.Marshal(&params)

		}
	})

	b.Run("jsoniter", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = jsoniter.Marshal(&params)

		}
	})

	b.Run("ffjson", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = ffjson.Marshal(&params)
		}
	})

	b.Run("go-json", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = gojson.Marshal(&params)
		}
	})
	//NOTE: fastjson is just a parser, it does not provide "Marshal" method.
}
