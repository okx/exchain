package common

import (
	"fmt"
	"testing"

	apptypes "github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func initConfig() {
	config := sdk.GetConfig()
	apptypes.SetBech32Prefixes(config)
	apptypes.SetBip44CoinType(config)
	config.Seal()
}

func TestHasSufCoins(t *testing.T) {
	initConfig()
	addr, err := sdk.AccAddressFromBech32("ex1rf9wr069pt64e58f2w3mjs9w72g8vemzw26658")
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
	InitConfig()
	addr := BlackHoleAddress()
	a := addr.String()
	fmt.Println(a)
	require.Equal(t, "ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqm2k6w2", addr.String())
}

func TestGetFixedLengthRandomString(t *testing.T) {
	require.Equal(t, 100, len(GetFixedLengthRandomString(100)))
}

func mustAccAddressFromHex(addr string) sdk.AccAddress {
	ret, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	return ret
}

func TestCheckSignerAddress(t *testing.T) {
	testcases := []struct {
		signers    []sdk.AccAddress
		delegators []sdk.AccAddress
		result     bool
	}{
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			result: true,
		},
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			result: false,
		},
		{
			signers:    []sdk.AccAddress{},
			delegators: []sdk.AccAddress{},
			result:     false,
		},
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			delegators: []sdk.AccAddress{},
			result:     false,
		},
		{
			signers: []sdk.AccAddress{},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			result: false,
		},

		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			result: false,
		},
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			result: false,
		},
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("bbE4733d85bc2b90682147779DA49caB38C0aA1F"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			result: false,
		},
		{
			signers: []sdk.AccAddress{
				mustAccAddressFromHex("c1fB47342851da0F7a6FD13866Ab37a2A125bE36"),
				mustAccAddressFromHex("bbE4733d85bc2b90682147779DA49caB38C0aA1F"),
			},
			delegators: []sdk.AccAddress{
				mustAccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				mustAccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"),
			},
			result: false,
		},
	}

	for _, ts := range testcases {
		tr := CheckSignerAddress(ts.signers, ts.delegators)
		require.Equal(t, ts.result, tr)
	}
}
