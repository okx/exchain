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

func TestParamsString(t *testing.T) {
	param := DefaultParams()
	expectString := `Order Params:
  OrderExpireBlocks: 259200
  MaxDealsPerBlock: 1000
  FeePerBlock: 0.00000100okt
  TradeFeeRate: 0.00100000`
	require.EqualValues(t, expectString, param.String())
}
