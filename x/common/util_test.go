package common

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

//------------------
// test sdk.DecCoin
func TestParseDecCoinByDecimal(t *testing.T) {

	decCoin, err := sdk.ParseDecCoin("1000.01" + NativeToken)
	require.Nil(t, err)
	require.Equal(t, "1000.01000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	require.Equal(t, uint64(100001000000), decCoin.Amount.Uint64())
	require.Equal(t, int64(100001000000), decCoin.Amount.Int64())
	require.Equal(t, false, decCoin.Amount.IsInteger())

	require.Equal(t, "1000.01000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "1000.01000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"1000.01000000\"", string(decCoinAmountJSON))
}

func TestParseDecCoinByInteger(t *testing.T) {

	decCoin, err := sdk.ParseDecCoin("1000" + NativeToken)
	require.Nil(t, err)

	require.Equal(t, "1000.00000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	require.Equal(t, uint64(100000000000), decCoin.Amount.Uint64())
	require.Equal(t, int64(100000000000), decCoin.Amount.Int64())
	require.Equal(t, true, decCoin.Amount.IsInteger())

	require.Equal(t, "1000.00000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "1000.00000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"1000.00000000\"", string(decCoinAmountJSON))
}

//--------------
// test sdk.Coin
func TestParseIntCoinByDecimal(t *testing.T) {
	_, err := sdk.ParseCoin("1000.1" + NativeToken)
	require.Nil(t, err)
}

//--------------------
// test sdk.NewCoin. Dangerous!
func TestSdkNewCoin(t *testing.T) {
	// dangerous to use!!!
	intCoin := sdk.NewCoin(NativeToken, sdk.NewInt(1000))
	require.Equal(t, "1000.00000000"+NativeToken, intCoin.String())
}

//--------------------
// test sdk.NewDecCoin
func TestSdkNewDecCoin(t *testing.T) {
	// safe to use
	intCoin := sdk.NewDecCoin(NativeToken, sdk.NewInt(1000))
	require.Equal(t, "1000.00000000"+NativeToken, intCoin.String())
}

func TestHasSufCoins(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("okexchain18mxjm0knqjpkaxk2zd2jr67pgrd8c0ct0tycvl")
	require.Nil(t, err)

	availDecCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		200000, "btc", 100000, NativeToken))
	require.Nil(t, err)
	availCoins := availDecCoins

	spendDecCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s",
		200000, NativeToken, 100000, "btc"))
	require.NoError(t, err)
	spendCoins := spendDecCoins

	err = HasSufficientCoins(addr, availCoins, spendCoins)
	require.Error(t, err)
	spendDecCoins, err = sdk.ParseDecCoins(fmt.Sprintf("%d%s",
		200000, "xmr"))
	require.Nil(t, err)
	spendCoins = spendDecCoins

	err = HasSufficientCoins(addr, availCoins, spendCoins)
	require.Error(t, err)

	spendDecCoins, err = sdk.ParseDecCoins(fmt.Sprintf("%d%s",
		100000, "btc"))
	require.Nil(t, err)
	spendCoins = spendDecCoins
	err = HasSufficientCoins(addr, availCoins, spendCoins)
	require.Nil(t, err)
}
