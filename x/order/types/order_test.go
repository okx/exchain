package types

import (
	"testing"

	"github.com/okex/okexchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	// test new order
	params := DefaultTestParams()
	order1 := NewOrder("hash1", nil, TestTokenPair, SellOrder, sdk.MustNewDecFromStr("1.1"),
		sdk.MustNewDecFromStr("10.0"), 123, params.OrderExpireBlocks, params.FeePerBlock)
	order2 := NewOrder("hash2", nil, TestTokenPair, BuyOrder, sdk.MustNewDecFromStr("1.1"),
		sdk.MustNewDecFromStr("10.0"), 123, params.OrderExpireBlocks, params.FeePerBlock)

	require.Equal(t, sdk.MustNewDecFromStr("10"), order1.RemainLocked)
	require.Equal(t, sdk.MustNewDecFromStr("11"), order2.RemainLocked)
	require.Equal(t, "Open", OrderStatus(order1.Status).String())

	// test order unlock
	order1.Unlock()
	require.Equal(t, sdk.ZeroDec().String(), order1.RemainLocked.String())

	// test order string
	expected := `{"txhash":"hash1","order_id":"","sender":"","product":"xxb_` + common.NativeToken + `","side":"SELL","price":"1.100000000000000000","quantity":"10.000000000000000000","status":0,"filled_avg_price":"0","remain_quantity":"10.000000000000000000","remain_locked":"0.000000000000000000","timestamp":123,"order_expire_blocks":259200,"fee_per_block":{"denom":"` + common.NativeToken + `","amount":"0.000001000000000000"},"extra_info":""}`

	require.Equal(t, expected, order1.String())

	order1.Status = 10
	require.Equal(t, "Unknown", OrderStatus(order1.Status).String())
}

func TestOrderUpdateExtraInfo(t *testing.T) {
	order := MockOrder("", "", SellOrder, "0.1", "10.0")
	order.setExtraInfoWithKeyValue(OrderExtraInfoKeyCancelFee, "0.002"+common.NativeToken)
	expectExtra := `{"cancelFee":"0.002` + common.NativeToken + `"}`
	require.EqualValues(t, expectExtra, order.ExtraInfo)
	require.EqualValues(t, "0.002"+common.NativeToken, order.GetExtraInfoWithKey(OrderExtraInfoKeyCancelFee))

	// Record deal fee
	fee := sdk.SysCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("0.01")}}
	order.RecordOrderDealFee(fee)
	require.EqualValues(t, fee.String(), order.GetExtraInfoWithKey(OrderExtraInfoKeyDealFee))

	order.RecordOrderDealFee(fee)
	require.EqualValues(t, fee.Add2(fee).String(),
		order.GetExtraInfoWithKey(OrderExtraInfoKeyDealFee))

	// Record new fee
	order.RecordOrderNewFee(fee)
	require.EqualValues(t, fee.String(), order.GetExtraInfoWithKey(OrderExtraInfoKeyNewFee))
	// Record cancel fee
	order.RecordOrderCancelFee(fee)
	require.EqualValues(t, fee.String(), order.GetExtraInfoWithKey(OrderExtraInfoKeyCancelFee))
	// Record expire fee
	order.recordOrderExpireFee(fee)
	require.EqualValues(t, fee.String(), order.GetExtraInfoWithKey(OrderExtraInfoKeyExpireFee))
}

func TestOrderFill(t *testing.T) {
	// Sell Order
	order := MockOrder("", "", SellOrder, "0.1", "10.0")
	order.Fill(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("5"))

	// Partial fill
	require.EqualValues(t, sdk.MustNewDecFromStr("5"), order.RemainQuantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.2"), order.FilledAvgPrice)
	require.EqualValues(t, sdk.MustNewDecFromStr("5.0"), order.RemainLocked)
	require.EqualValues(t, OrderStatusOpen, order.Status)

	// Full fill
	order.Fill(sdk.MustNewDecFromStr("0.1"), sdk.MustNewDecFromStr("5"))
	require.True(t, order.RemainQuantity.IsZero())
	require.EqualValues(t, sdk.ZeroDec().String(), order.RemainLocked.String())
	require.EqualValues(t, sdk.MustNewDecFromStr("0.15"), order.FilledAvgPrice)
	require.EqualValues(t, OrderStatusFilled, order.Status)
	require.EqualValues(t, "Filled", OrderStatus(order.Status).String())

	// Buy order
	order = MockOrder("", "", BuyOrder, "0.1", "10.0")
	order.Fill(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("5"))

	// Partial fill
	require.EqualValues(t, sdk.MustNewDecFromStr("5"), order.RemainQuantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.05"), order.FilledAvgPrice)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.75"), order.RemainLocked)
	require.EqualValues(t, OrderStatusOpen, order.Status)

	// Full fill
	order.Fill(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("5"))
	require.True(t, order.RemainQuantity.IsZero())
	require.EqualValues(t, sdk.MustNewDecFromStr("0.5"), order.RemainLocked)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.05"), order.FilledAvgPrice)
	require.EqualValues(t, OrderStatusFilled, order.Status)
}

func TestOrderCancel(t *testing.T) {
	// Full cancel
	order := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	order.Cancel()
	require.EqualValues(t, OrderStatusCancelled, order.Status)
	require.EqualValues(t, "Cancelled", OrderStatus(order.Status).String())

	// Partial cancel
	order2 := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	order2.Fill(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("5"))
	order2.Cancel()
	require.EqualValues(t, OrderStatusPartialFilledCancelled, order2.Status)
	require.EqualValues(t, "PartialFilledCancelled", OrderStatus(order2.Status).String())
}

func TestOrderExpire(t *testing.T) {
	// Full expire
	order := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	order.Expire()
	require.EqualValues(t, OrderStatusExpired, order.Status)
	require.EqualValues(t, "Expired", OrderStatus(order.Status).String())

	// Partial expire
	order2 := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	order2.Fill(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("5"))
	order2.Expire()
	require.EqualValues(t, OrderStatusPartialFilledExpired, order2.Status)
	require.EqualValues(t, "PartialFilledExpired", OrderStatus(order2.Status).String())
}

func TestOrderNeedLockCoins(t *testing.T) {
	order := MockOrder("", TestTokenPair, BuyOrder, "0.1", "10.0")
	decCoins := order.NeedLockCoins()
	require.EqualValues(t, "1.000000000000000000"+common.NativeToken, decCoins.String())

	order2 := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	decCoins = order2.NeedLockCoins()
	require.EqualValues(t, "10.000000000000000000xxb", decCoins.String())
}

func TestOrderNeedUnlockCoins(t *testing.T) {
	order := MockOrder("", TestTokenPair, BuyOrder, "0.1", "10.0")
	decCoins := order.NeedUnlockCoins()
	require.EqualValues(t, "1.000000000000000000"+common.NativeToken, decCoins.String())
	order.Fill(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("5"))
	decCoins = order.NeedUnlockCoins()
	require.EqualValues(t, "0.750000000000000000"+common.NativeToken, decCoins.String()) // 0.1 * 10 - 0.05 * 5

	order2 := MockOrder("", TestTokenPair, SellOrder, "0.1", "10.0")
	decCoins = order2.NeedUnlockCoins()
	require.EqualValues(t, "10.000000000000000000xxb", decCoins.String())
	order2.Fill(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("6"))
	decCoins = order2.NeedUnlockCoins()
	require.EqualValues(t, "4.000000000000000000xxb", decCoins.String())
}

func TestGetBlockHeightFromOrderID(t *testing.T) {
	var blockHeight int64 = 100
	var orderNum int64 = 2
	orderID := FormatOrderID(blockHeight, orderNum)
	num := GetBlockHeightFromOrderID(orderID)
	require.Equal(t, blockHeight, num)

	blockHeight = 99999999990
	orderID = FormatOrderID(blockHeight, orderNum)
	num = GetBlockHeightFromOrderID(orderID)
	require.Equal(t, blockHeight, num)
}
