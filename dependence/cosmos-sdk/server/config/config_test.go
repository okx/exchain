package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.True(t, sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("0.0000000001")).IsEqual(cfg.GetMinGasPrices()))
}

func TestSetMinimumFees(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinGasPrices(sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)})
	require.Equal(t, "5.000000000000000000foo", cfg.MinGasPrices)
}
