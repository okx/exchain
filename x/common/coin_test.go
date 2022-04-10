package common

import (
	"fmt"
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

//------------------
// test sdk.SysCoin
func TestParseDecCoinByDecimal(t *testing.T) {

	decCoin, err := sdk.ParseDecCoin("1000.01" + NativeToken)
	require.Nil(t, err)
	require.Equal(t, "1000.010000000000000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(100001000000), decCoin.Amount.Uint64())
	//require.Equal(t, int64(100001000000), decCoin.Amount.Int64())
	require.Equal(t, false, decCoin.Amount.IsInteger())

	require.Equal(t, "1000.010000000000000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "1000.010000000000000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"1000.010000000000000000\"", string(decCoinAmountJSON))
}

func TestDecCoin(t *testing.T) {

	decCoin, err := sdk.ParseDecCoin("0.000000000000000100" + NativeToken)
	require.Nil(t, err)
	require.Equal(t, "0.000000000000000100"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	require.Equal(t, uint64(100), decCoin.Amount.Uint64())
	require.Equal(t, int64(100), decCoin.Amount.Int64())

	fmt.Printf("%v\n", decCoin.Amount.Uint64())
	fmt.Printf("%v\n", decCoin.Amount.Int64())
	require.Equal(t, false, decCoin.Amount.IsInteger())

	require.Equal(t, "0.000000000000000100", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "0.000000000000000100", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"0.000000000000000100\"", string(decCoinAmountJSON))
}
func TestDecCoin3(t *testing.T) {

	decCoin := sdk.NewDecCoinFromDec(NativeToken, sdk.NewDec(9))
	//require.Nil(t, err)
	require.Equal(t, "9.000000000000000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(600000000000000000000), decCoin.Amount.Uint64())
	//require.Equal(t, int64(600000000000000000000), decCoin.Amount.Int64())

	fmt.Printf("%v\n", decCoin.Amount.Uint64())
	fmt.Printf("%v\n", decCoin.Amount.Int64())
	require.Equal(t, true, decCoin.Amount.IsInteger())

	require.Equal(t, "9.000000000000000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "9.000000000000000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"9.000000000000000000\"", string(decCoinAmountJSON))
}
func TestDecCoin5(t *testing.T) {

	num := "299999999999999999999999999999999999999999999999999"

	int1, ok := sdk.NewIntFromString(num)
	if !ok {
		panic("")
	}
	fmt.Printf("int1: %v\n", int1.String())

	dec1 := sdk.MustNewDecFromStr(num + ".000000000000000000")
	decCoin := sdk.NewDecCoinFromDec(NativeToken, dec1)
	//require.Nil(t, err)
	require.Equal(t, num+".000000000000000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(600000000000000000000), decCoin.Amount.Uint64())
	//require.Equal(t, int64(600000000000000000000), decCoin.Amount.Int64())

	fmt.Printf("%v\n", decCoin.Amount.Uint64())
	fmt.Printf("%v\n", decCoin.Amount.Int64())
	require.Equal(t, true, decCoin.Amount.IsInteger())

	require.Equal(t, num+".000000000000000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, num+".000000000000000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\""+num+".000000000000000000\"", string(decCoinAmountJSON))
}
func TestDecCoin4(t *testing.T) {

	decCoin := sdk.NewDecCoinFromDec(NativeToken, sdk.NewDec(9))
	//require.Nil(t, err)
	require.Equal(t, "9.000000000000000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(600000000000000000000), decCoin.Amount.Uint64())
	//require.Equal(t, int64(600000000000000000000), decCoin.Amount.Int64())

	fmt.Printf("%v\n", decCoin.Amount.Uint64())
	fmt.Printf("%v\n", decCoin.Amount.Int64())
	require.Equal(t, true, decCoin.Amount.IsInteger())

	require.Equal(t, "9.000000000000000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "9.000000000000000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"9.000000000000000000\"", string(decCoinAmountJSON))
}

func TestIntCoin(t *testing.T) {
	intCoin, err := sdk.ParseCoin("100001" + NativeToken)
	require.Nil(t, err)

	fmt.Printf("%v\n", intCoin.String())

	fmt.Printf("%v\n", intCoin.Amount.Uint64())
	fmt.Printf("%v\n", intCoin.Amount.Int64())

	require.Equal(t, "100001.000000000000000000"+NativeToken, intCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(100001), intCoin.Amount.Uint64())
	//require.Equal(t, int64(100001), intCoin.Amount.Int64())

	require.Equal(t, "100001.000000000000000000", intCoin.Amount.String())

	intCoinAmountYaml, err := intCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "100001.000000000000000000", intCoinAmountYaml)

	intCoinAmountJSON, err := intCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"100001.000000000000000000\"", string(intCoinAmountJSON))
}

func TestIntCoin2(t *testing.T) {
	intCoin, err := sdk.ParseCoin("100001" + NativeToken)
	require.Nil(t, err)

	fmt.Printf("%v\n", intCoin.String())

	fmt.Printf("%v\n", intCoin.Amount.Uint64())
	fmt.Printf("%v\n", intCoin.Amount.Int64())

	require.Equal(t, "100001.000000000000000000"+NativeToken, intCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(100001), intCoin.Amount.Uint64())
	//require.Equal(t, int64(100001), intCoin.Amount.Int64())

	require.Equal(t, "100001.000000000000000000", intCoin.Amount.String())

	intCoinAmountYaml, err := intCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "100001.000000000000000000", intCoinAmountYaml)

	intCoinAmountJSON, err := intCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"100001.000000000000000000\"", string(intCoinAmountJSON))
}

func TestParseDecCoinByInteger(t *testing.T) {

	decCoin, err := sdk.ParseDecCoin("1000" + NativeToken)
	require.Nil(t, err)

	require.Equal(t, "1000.000000000000000000"+NativeToken, decCoin.String())

	//----------------
	// test sdk.Dec
	//require.Equal(t, uint64(100000000000), decCoin.Amount.Uint64())
	//require.Equal(t, int64(100000000000), decCoin.Amount.Int64())
	require.Equal(t, true, decCoin.Amount.IsInteger())

	require.Equal(t, "1000.000000000000000000", decCoin.Amount.String())

	decCoinAmountYaml, err := decCoin.Amount.MarshalYAML()
	require.Nil(t, err)
	require.Equal(t, "1000.000000000000000000", decCoinAmountYaml)

	decCoinAmountJSON, err := decCoin.Amount.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"1000.000000000000000000\"", string(decCoinAmountJSON))
}

//--------------
// test sdk.Coin
func TestParseIntCoinByDecimal(t *testing.T) {
	ret, err := sdk.ParseCoin("1000.1" + NativeToken)
	require.Nil(t, err)
	fmt.Println(ret.String())
}

//--------------------
// test sdk.NewCoin. Dangerous!
func TestSdkNewCoin(t *testing.T) {
	// dangerous to use!!!
	intCoin := sdk.NewCoin(NativeToken, sdk.NewInt(1000))
	require.Equal(t, "1000.000000000000000000"+NativeToken, intCoin.String())
}

//--------------------
// test sdk.NewDecCoin
func TestSdkNewDecCoin(t *testing.T) {
	// safe to use
	intCoin := sdk.NewDecCoin(NativeToken, sdk.NewInt(1000))
	require.Equal(t, "1000.000000000000000000"+NativeToken, intCoin.String())
}
