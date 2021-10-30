package keeper

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"strings"

	"github.com/okex/exchain/x/order/types"
)

const minFee = "0.00000001"

// GetFeeKeeper is an interface for calculating handling fees
type GetFeeKeeper interface {
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
}

// GetOrderNewFee is used to calculate the handling fee that needs to be locked when placing an order
func GetOrderNewFee(order *types.Order) sdk.SysCoins {
	orderExpireBlocks := sdk.NewDec(order.OrderExpireBlocks)
	amount := order.FeePerBlock.Amount.Mul(orderExpireBlocks)
	return sdk.SysCoins{sdk.NewDecCoinFromDec(order.FeePerBlock.Denom, amount)}
}

// GetOrderCostFee is used to calculate the handling fee when quiting an order
func GetOrderCostFee(order *types.Order, ctx sdk.Context) sdk.SysCoins {
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
	return sdk.SysCoins{sdk.NewDecCoinFromDec(order.FeePerBlock.Denom, costFee)}

}

// GetZeroFee returns zeroFee
func GetZeroFee() sdk.SysCoins {
	return sdk.SysCoins{sdk.ZeroFee()}
}

// GetDealFee is used to calculate the handling fee when matching an order
func GetDealFee(order *types.Order, fillAmt sdk.Dec, ctx sdk.Context, keeper GetFeeKeeper,
	feeParams *types.Params) sdk.SysCoins {
	symbols := strings.Split(order.Product, "_")
	symbol := symbols[0]
	quantity := fillAmt
	if order.Side == types.SellOrder {
		symbol = symbols[1]
		quantity = fillAmt.Mul(keeper.GetLastPrice(ctx, order.Product))
	}

	minFeeDec := sdk.MustNewDecFromStr(minFee)
	feeAmt := quantity.Mul(feeParams.TradeFeeRate)
	if feeAmt.GT(minFeeDec) {
		return sdk.SysCoins{sdk.NewDecCoinFromDec(symbol, feeAmt)}
	}
	return sdk.SysCoins{sdk.NewDecCoinFromDec(symbol, minFeeDec)}
}
