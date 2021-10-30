package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestYieldedTokenInfo(t *testing.T) {
	yieldInfo1 := NewYieldedTokenInfo(
		sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)), 100, sdk.NewDec(10),
	)
	yieldInfo2 := NewYieldedTokenInfo(
		sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)), 50, sdk.NewDec(20),
	)
	yieldInfos := NewYieldedTokenInfos(yieldInfo1, yieldInfo2)

	require.Equal(t, yieldInfos.String(), yieldInfo1.String()+"\n"+yieldInfo2.String())
}
