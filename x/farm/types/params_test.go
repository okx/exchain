package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	strExpected = `Params:
  Quote Symbol:								usdk
  Create Pool Fee:							0.00000000okt
  Create Pool Deposit:						10.00000000okt
  Yield Native Token Enabled:               false`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := DefaultParams()

	require.Equal(t, defaultState.Params, defaultParams)
	require.Equal(t, strExpected, defaultParams.String())
}
