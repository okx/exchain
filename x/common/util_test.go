package common

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	apptypes "github.com/okex/okexchain/app/types"
)

func initConfig() {
	config := sdk.GetConfig()
	apptypes.SetBech32Prefixes(config)
	apptypes.SetBip44CoinType(config)
	config.Seal()
}

func TestHasSufCoins(t *testing.T) {
	initConfig()
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

func TestBlackHoleAddress(t *testing.T) {
	addr := BlackHoleAddress()
	a := addr.String()
	fmt.Println(a)
	require.Equal(t, addr.String(), "okexchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqupa6dx")
}

func TestGetFixedLengthRandomString(t *testing.T) {
	require.Equal(t, 100, len(GetFixedLengthRandomString(100)))
}

func TestForPanicTrace(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			PanicTrace(4)
			//os.Exit(1)
		}
	}()
	panic("just for test")
}
