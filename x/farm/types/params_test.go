package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	strExpected = `Params:
  Quote Symbol:								usdk
  Create Pool Fee:							0.00000000okt
  Create Pool Deposit:						10.00000000okt
  Manage White List Max Deposit Period:		24h0m0s
  Manage White List Min Deposit:			100.00000000okt
  Manage White List Voting Period:			72h0m0s`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := DefaultParams()

	require.Equal(t, defaultState.Params, defaultParams)

	fmt.Println(defaultParams.String())
	require.Equal(t, strExpected, defaultParams.String())
}
