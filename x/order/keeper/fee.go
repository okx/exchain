package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"strings"

	"github.com/okex/okchain/x/order/types"
)

const MinFee = "0.00000001"

type GetFeeKeeper interface {
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
}

// Currently, placing order does not need any fee, so we only support charging okb if necessary
func GetOrderNewFee(order *types.Order) sdk.DecCoins {
	orderExpireBlocks := sdk.NewDec(order.OrderExpireBlocks)
	amount := order.FeePerBlock.Amount.Mul(orderExpireBlocks)
	return sdk.DecCoins{sdk.NewDecCoinFromDec(order.FeePerBlock.Denom, amount)}
}

func GetOrderCostFee(order *types.Order, ctx sdk.Context) sdk.DecCoins {
	currentHeight := ctx.BlockHeight()
	orderHeight := types.GetBlockHeightFromOrderID(order.OrderID)
	blockNum := currentHeight - orderHeight
	if blockNum < 0 {
		ctx.Logger().Error(fmt.Sprintf("currentHeight(%d) should not less than orderHeight(%d)", currentHeight, orderHeight))
		return GetZeroFee()
	} else if blockNum > order.OrderExpireBlocks {
		blockNum = order.OrderExpireBlocks
		ctx.Logger().Error(fmt.Sprintf("currentHeight(%d) - orderHeight(%d) > OrderExpireBlocks(%d)", currentHeight, orderHeight, order.OrderExpireBlocks))
	}
	costFee := order.FeePerBlock.Amount.Mul(sdk.NewDec(blockNum))
	return sdk.DecCoins{sdk.NewDecCoinFromDec(order.FeePerBlock.Denom, costFee)}

}

func GetZeroFee() sdk.DecCoins {
	return sdk.DecCoins{sdk.ZeroFee()}
}

func GetDealFee(order *types.Order, fillAmt sdk.Dec, ctx sdk.Context, keeper GetFeeKeeper,
	feeParams *types.Params) sdk.DecCoins {
	symbols := strings.Split(order.Product, "_")
	symbol := symbols[0]
	quantity := fillAmt
	if order.Side == types.SellOrder {
		symbol = symbols[1]
		quantity = fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	}

	feeAmt := quantity.Mul(feeParams.TradeFeeRate)
	if feeAmt.IsPositive() {
		return sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, feeAmt)}
	}
	return sdk.DecCoins{sdk.NewDecCoinFromDec(symbol, sdk.MustNewDecFromStr(MinFee))}
}
