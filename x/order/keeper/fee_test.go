package keeper

import (
	"testing"

	"github.com/okex/exchain/x/common"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/order/types"
	"github.com/okex/exchain/dependence/tendermint/libs/cli/flags"
)

type MockGetFeeKeeper struct {
	coins    sdk.Coins
	priceMap map[string]sdk.Dec
}

func NewMockGetFeeKeeper() MockGetFeeKeeper {
	return MockGetFeeKeeper{sdk.NewCoins(), make(map[string]sdk.Dec)}
}

func (k MockGetFeeKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.coins
}

func (k MockGetFeeKeeper) GetLastPrice(ctx sdk.Context, product string) sdk.Dec {
	if price, ok := k.priceMap[product]; ok {
		return price
	}
	return sdk.ZeroDec()
}

func TestGetOrderNewFee(t *testing.T) {
	order := mockOrder("ID0000001970-1", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	orderExpireBlocks := sdk.NewDec(order.OrderExpireBlocks)
	exceptFee := sdk.SysCoins{sdk.NewDecCoinFromDec(common.NativeToken, order.FeePerBlock.Amount.Mul(orderExpireBlocks))}
	require.EqualValues(t, exceptFee, GetOrderNewFee(order))
}

func TestGetOrderCostFee(t *testing.T) {
	var orderHeight int64 = 1970
	var currentHeight int64 = 2970
	diffHeight := currentHeight - orderHeight
	orderID := types.FormatOrderID(orderHeight, 2)
	order := mockOrder(orderID, types.TestTokenPair, types.BuyOrder, "10.0", "1.0")

	testInput := CreateTestInput(t)
	log, err := flags.ParseLogLevel("*:error", testInput.Ctx.Logger(), "error")
	require.Nil(t, err)
	ctx := testInput.Ctx
	ctx = ctx.WithLogger(log)
	ctx = ctx.WithBlockHeight(currentHeight)
	exceptFee := sdk.SysCoins{sdk.NewDecCoinFromDec(common.NativeToken, order.FeePerBlock.Amount.Mul(sdk.NewDec(diffHeight)))}
	require.EqualValues(t, exceptFee, GetOrderCostFee(order, ctx))

	ctx = ctx.WithBlockHeight(currentHeight + types.DefaultOrderExpireBlocks)
	fee := GetOrderCostFee(order, ctx)
	exceptFee = sdk.SysCoins{sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("0.2592"))}
	require.EqualValues(t, exceptFee, fee)

	currentHeight = 0
	ctx = ctx.WithBlockHeight(currentHeight)
	exceptFee = GetZeroFee()
	require.EqualValues(t, exceptFee, GetOrderCostFee(order, ctx))

}

func TestOrderDealFee(t *testing.T) {
	ctx := sdk.Context{}
	keeper := NewMockGetFeeKeeper()
	feeParams := types.DefaultTestParams()

	// 1. xxb_okb
	order := &types.Order{
		Product:  types.TestTokenPair,
		Side:     types.BuyOrder,
		Price:    sdk.MustNewDecFromStr("11.0"),
		Quantity: sdk.MustNewDecFromStr("100.0"),
	}
	keeper.priceMap[types.TestTokenPair] = sdk.MustNewDecFromStr("10.0")
	feeOther := GetDealFee(order, sdk.MustNewDecFromStr("10.0"), ctx, keeper, &feeParams)
	// 10 * 0.001
	expectFee := sdk.SysCoins{sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("0.01"))}
	require.EqualValues(t, expectFee, feeOther)

	// xxb_yyb BUY
	order = &types.Order{
		Product:  "xxb_yyb",
		Side:     types.BuyOrder,
		Price:    sdk.MustNewDecFromStr("11.0"),
		Quantity: sdk.MustNewDecFromStr("100.0"),
	}
	keeper.priceMap["xxb_yyb"] = sdk.MustNewDecFromStr("20.0")
	keeper.priceMap["yyb_"+common.NativeToken] = sdk.MustNewDecFromStr("0.6")

	feeOther = GetDealFee(order, sdk.MustNewDecFromStr("100.0"), ctx, keeper, &feeParams)
	// 100 * 0.001
	expectFee = sdk.SysCoins{sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("0.1"))}
	require.EqualValues(t, expectFee, feeOther)

	// xxb_yyb SELL
	order = &types.Order{
		Product:  "xxb_yyb",
		Side:     types.SellOrder,
		Price:    sdk.MustNewDecFromStr("11.0"),
		Quantity: sdk.MustNewDecFromStr("100.0"),
	}
	feeOther = GetDealFee(order, sdk.MustNewDecFromStr("100.0"), ctx, keeper, &feeParams)
	// 100 * 20 * 0.001
	expectFee = sdk.SysCoins{sdk.NewDecCoinFromDec("yyb", sdk.MustNewDecFromStr("2.0"))}
	require.EqualValues(t, expectFee, feeOther)

	// xxb_yyb SELL
	order = &types.Order{
		Product:  "xxb_yyb",
		Side:     types.BuyOrder,
		Price:    sdk.MustNewDecFromStr("1.0"),
		Quantity: sdk.MustNewDecFromStr("0.00000001"),
	}
	feeOther = GetDealFee(order, sdk.MustNewDecFromStr("0.000000000000000001"), ctx, keeper, &feeParams)
	expectFee = sdk.SysCoins{sdk.NewDecCoinFromDec("xxb", sdk.MustNewDecFromStr(minFee))}
	require.EqualValues(t, expectFee, feeOther)
}
