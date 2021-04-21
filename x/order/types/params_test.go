package types

import (
	"testing"

	"github.com/okex/exchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamSetPairs(t *testing.T) {
	tests := []Params{
		{
			OrderExpireBlocks:     1000,
			MaxDealsPerBlock:      10000,
			FeePerBlock:           sdk.NewDecCoinFromDec(DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr("0.000001")),
			TradeFeeRate:          sdk.MustNewDecFromStr("0.001"),
			NewOrderMsgGasUnit:    123,
			CancelOrderMsgGasUnit: 456,
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
				if !v.Value.(*sdk.SysCoin).IsEqual(test.FeePerBlock) {
					t.Errorf("key(%s) -> %x, want %x", v.Key, test.FeePerBlock, v.Value)
				}
			case string(KeyTradeFeeRate):
				if !v.Value.(*sdk.Dec).Equal(test.TradeFeeRate) {
					t.Errorf("key(%s) -> %x, want %x", v.Key, test.TradeFeeRate, v.Value)
				}
			case string(KeyNewOrderMsgGasUnit):
				require.EqualValues(t, test.NewOrderMsgGasUnit, *(v.Value.(*uint64)))
			case string(KeyCancelOrderMsgGasUnit):
				require.EqualValues(t, test.CancelOrderMsgGasUnit, *(v.Value.(*uint64)))
			}
		}
	}
}

func TestParamsString(t *testing.T) {
	param := DefaultParams()
	expectString := `Order Params:
  OrderExpireBlocks: 259200
  MaxDealsPerBlock: 1000
  FeePerBlock: 0.000000000000000000` + common.NativeToken + `
  TradeFeeRate: 0.001000000000000000
  NewOrderMsgGasUnit: 40000
  CancelOrderMsgGasUnit: 30000`
	require.EqualValues(t, expectString, param.String())
}
