package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	//"github.com/stretchr/testify/require"
	"testing"
)

func TestDecay(t *testing.T) {
	tokens := sdk.NewDec(1000)
	nowDec := calculateWeight(tokens)
	afterDec := calculateWeight(tokens)
	require.Equal(t, afterDec, nowDec)
	require.Equal(t, afterDec, tokens)
	require.Equal(t, nowDec, tokens)
}
