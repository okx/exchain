package types

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	strExpected = `Params:
  Quote Symbol:								usdk
  Create Pool Fee:							0.000000000000000000` + sdk.DefaultBondDenom + `
  Create Pool Deposit:						10.000000000000000000` + sdk.DefaultBondDenom + `
  Yield Native Token Enabled:               false`
)

func TestParams(t *testing.T) {
	defaultState := DefaultGenesisState()
	defaultParams := DefaultParams()

	require.Equal(t, defaultState.Params, defaultParams)
	require.Equal(t, strExpected, defaultParams.String())
}
