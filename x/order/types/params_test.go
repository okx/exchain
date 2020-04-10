package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamSetPairs(t *testing.T) {
	tests := []Params{
		{
			OrderExpireBlocks: 1000,
			MaxDealsPerBlock:  10000,
			FeePerBlock:       sdk.NewDecCoinFromDec(DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr("0.000001")),
			TradeFeeRate:      sdk.MustNewDecFromStr("0.001"),
		},
	}

	for _, test := range tests {
		p := test.ParamSetPairs()
		for _, v := range p {
			switch string(v.Key) {
			case string(KeyOrderExpireBlocks):
				require.EqualValues(t, test.OrderExpireBlocks, *(v.Value.(*int64)))
			case string(KeyMaxDealsPerBlock):
				require.EqualValues(t, test.MaxDealsPerBlock, *(v.Value.(*int64)))
			case string(KeyFeePerBlock):
				if !v.Value.(*sdk.DecCoin).IsEqual(test.FeePerBlock) {
					t.Errorf("key(%s) -> %x, want %x", v.Key, test.FeePerBlock, v.Value)
				}
			case string(KeyTradeFeeRate):
				if !v.Value.(*sdk.Dec).Equal(test.TradeFeeRate) {
					t.Errorf("key(%s) -> %x, want %x", v.Key, test.TradeFeeRate, v.Value)
				}
			}

		}
	}
}

//func TestValidateKV(t *testing.T) {
//	tests := []struct {
//		key   string
//		value string
//		valid bool
//	}{
//		{string(KeyOrderExpireBlocks), "1000", true},
//		{string(KeyOrderExpireBlocks), "something", false},
//		{string(KeyMaxDealsPerBlock), "10000", true},
//		{string(KeyMaxDealsPerBlock), "something", false},
//		{string(KeyNewOrder), "0", true},
//		{string(KeyNewOrder), "something", false},
//		{string(KeyCancel), "0.01", true},
//		{string(KeyCancel), "something", false},
//		{string(KeyCancelNative), "0.001", true},
//		{string(KeyCancelNative), "something", false},
//		{string(KeyExpire), "2.0", true},
//		{string(KeyExpire), "something", false},
//		{string(KeyExpireNative), "1.00", true},
//		{string(KeyExpireNative), "something", false},
//		{string(KeyTradeFeeRate), "0.004", true},
//		{string(KeyTradeFeeRate), "something", false},
//		{string(KeyTradeFeeRateNative), "0.001", true},
//		{string(KeyTradeFeeRateNative), "something", false},
//		{"default", "something", false},
//	}
//	p := Params{}
//	for _, test := range tests {
//		_, err := p.ValidateKV(test.key, test.value)
//		if err != nil && test.valid {
//			t.Errorf("key(%s) -> %s, error %s", test.key, test.value, err.Error())
//		}
//	}
//
//}

func TestParamsString(t *testing.T) {
	param := DefaultParams()
	expectString := "Params: \nOrderExpireBlocks: 259200\nMaxDealsPerBlock: 1000\nFeePerBlock: 0.00000100okt\nTradeFeeRate: 0.00100000\n"
	require.EqualValues(t, expectString, param.String())
}
