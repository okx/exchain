package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	strExpected = `Params:
  Quote Symbol:								usdk
  Create Pool Fee:							0.000000000000000000tokt
  Create Pool Deposit:						10.000000000000000000tokt
  Yield Native Token Enabled:               false`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := DefaultParams()

	require.Equal(t, defaultState.Params, defaultParams)
	require.Equal(t, strExpected, defaultParams.String())
}
