package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	strExpected = `Distribution Params:
  Community Tax:          0.02000000
  Withdraw Addr Enabled:  true`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := NewParams(defaultState.CommunityTax, defaultState.WithdrawAddrEnabled)
	require.Equal(t, defaultState.CommunityTax, defaultParams.CommunityTax)
	require.Equal(t, defaultState.WithdrawAddrEnabled, defaultParams.WithdrawAddrEnabled)

	require.Equal(t, strExpected, defaultParams.String())
	yamlStr, err := defaultParams.MarshalYAML()
	require.NoError(t, err)
	require.Equal(t, strExpected, yamlStr)
}
