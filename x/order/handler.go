package order

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
)

// NewOrderHandler returns the handler with version 0.
func NewOrderHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() sdk.Result
		var name string
		logger := ctx.Logger().With("module", "order")
		switch msg := msg.(type) {
		case types.MsgNewOrders:
			name = "handleMsgNewOrders"
			handlerFun = func() sdk.Result {
				return handleMsgNewOrders(ctx, keeper, msg, logger)
			}
		case types.MsgCancelOrders:
			name = "handleMsgCancelOrders"
			handlerFun = func() sdk.Result {
				return handleMsgCancelOrders(ctx, keeper, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("Invalid msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		return handlerFun()
	}
}

// checkOrderNewMsg: check msg product, price & quantity fields
func checkOrderNewMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgNewOrder) error {
	tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, msg.Product)
	if tokenPair == nil {
		return fmt.Errorf("trading pair '%s' does not exist", msg.Product)
	}

	// check if the order is involved with the tokenpair in dex Delist
	isDelisting, err := keeper.GetDexKeeper().CheckTokenPairUnderDexDelist(ctx, msg.Product)
	if err != nil {
		return err
	}
	if isDelisting {
		return errors.Errorf("trading pair '%s' is delisting", msg.Product)
	}

	priceDigit := tokenPair.MaxPriceDigit
	quantityDigit := tokenPair.MaxQuantityDigit
	roundedPrice := msg.Price.RoundDecimal(priceDigit)
	roundedQuantity := msg.Quantity.RoundDecimal(quantityDigit)
	if !roundedPrice.Equal(msg.Price) {
		return fmt.Errorf("price(%v) over accuracy(%d)", msg.Price, priceDigit)
	}
	if !roundedQuantity.Equal(msg.Quantity) {
		return fmt.Errorf("quantity(%v) over accuracy(%d)", msg.Quantity, quantityDigit)
	}

	if msg.Quantity.LT(tokenPair.MinQuantity) {
		return fmt.Errorf("quantity should be greater than %s", tokenPair.MinQuantity)
	}
	var d int64 = 100000000
	baseQuantity := msg.Price.Mul(msg.Quantity)
	if !msg.Price.MulInt64(d).Mul(msg.Quantity).Equal(baseQuantity.MulInt64(d)) {
		return fmt.Errorf("price(%v) * quantity(%v) over accuracy(%d)", msg.Price, msg.Quantity, priceDigit)
	}
	return nil
}

func getOrderFromMsg(ctx sdk.Context, k keeper.Keeper, msg types.MsgNewOrder, ratio string) *types.Order {
	feeParams := k.GetParams(ctx)
	feePerBlockAmount := feeParams.FeePerBlock.Amount.Mul(sdk.MustNewDecFromStr(ratio))
	feePerBlock := sdk.NewDecCoinFromDec(feeParams.FeePerBlock.Denom, feePerBlockAmount)
	return types.NewOrder(
		fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())),
		msg.Sender,
		msg.Product,
		msg.Side,
		msg.Price,
		msg.Quantity,
		ctx.BlockHeader().Time.Unix(),
		feeParams.OrderExpireBlocks,
		feePerBlock,
	)
}

func handleNewOrder(ctx sdk.Context, k Keeper, sender sdk.AccAddress,
	item types.OrderItem, ratio string, logger log.Logger) (types.OrderResult, sdk.CacheMultiStore, error) {

	cacheItem := ctx.MultiStore().CacheMultiStore()
	ctxItem := ctx.WithMultiStore(cacheItem)
	msg := MsgNewOrder{
		Sender:   sender,
		Product:  item.Product,
		Side:     item.Side,
		Price:    item.Price,
		Quantity: item.Quantity,
	}
	order := getOrderFromMsg(ctxItem, k, msg, ratio)
	code := sdk.CodeOK
	err := checkOrderNewMsg(ctxItem, k, msg)

	if err != nil {
		code = sdk.CodeUnknownRequest
	} else {
		if k.IsProductLocked(msg.Product) {
			code = sdk.CodeInternal
			err = fmt.Errorf("the trading pair (%s) is locked, please retry later", order.Product)
		} else if err = k.PlaceOrder(ctxItem, order); err != nil {
			code = sdk.CodeInsufficientCoins
		}
	}

	res := types.OrderResult{
		Code:    code,
		OrderID: order.OrderID,
	}

	if err == nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"    msg<Product:%s,Sender:%s,Price:%s,Quantity:%s,Side:%s>\n"+
			"    TxHash<%s>, Status<%s>\n"+
			"    result<The User have created an order {ID:%s,RemainQuantity:%s,Status:%s} >\n",
			ctx.BlockHeight(), "handleMsgNewOrder",
			msg.Product, msg.Sender, msg.Price.String(), msg.Quantity.String(), msg.Side,
			order.TxHash, types.OrderStatus(types.OrderStatusOpen),
			order.OrderID, order.RemainQuantity.String(), types.OrderStatus(order.Status)))
	} else {
		res.Message = err.Error()
	}

	return res, cacheItem, err
}

func handleMsgNewOrders(ctx sdk.Context, k Keeper, msg types.MsgNewOrders,
	logger log.Logger) sdk.Result {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	ratio := "1"
	if len(msg.OrderItems) > 1 {
		ratio = "0.8"
	}

	rs := make([]types.OrderResult, 0, len(msg.OrderItems))
	for _, item := range msg.OrderItems {
		res, cacheItem, err := handleNewOrder(ctx, k, msg.Sender, item, ratio, logger)
		if err == nil {
			cacheItem.Write()
		}
		rs = append(rs, res)
	}
	rss, err := json.Marshal(&rs)
	if err != nil {
		rss = []byte(fmt.Sprintf("failed to marshal result to JSON: %s", err))
	}
	event = event.AppendAttributes(sdk.NewAttribute("orders", string(rss)))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// ValidateMsgNewOrders validates whether the msg of newOrders is valid.
func ValidateMsgNewOrders(ctx sdk.Context, k keeper.Keeper, msg types.MsgNewOrders) sdk.Result {
	ratio := "1"
	if len(msg.OrderItems) > 1 {
		ratio = "0.8"
	}

	for _, item := range msg.OrderItems {
		msg := MsgNewOrder{
			Sender:   msg.Sender,
			Product:  item.Product,
			Side:     item.Side,
			Price:    item.Price,
			Quantity: item.Quantity,
		}
		err := checkOrderNewMsg(ctx, k, msg)
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeUnknownRequest,
				Log:  err.Error(),
			}
		}
		if k.IsProductLocked(msg.Product) {
			return sdk.Result{
				Code: sdk.CodeInternal,
				Log:  fmt.Sprintf("the trading pair (%s) is locked, please retry later", msg.Product),
			}
		}

		order := getOrderFromMsg(ctx, k, msg, ratio)
		_, err = k.TryPlaceOrder(ctx, order)
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeInsufficientCoins,
				Log:  err.Error(),
			}
		}
	}

	return sdk.Result{}

}

func handleCancelOrder(context sdk.Context, k Keeper, sender sdk.AccAddress, orderID string, logger log.Logger) (
	types.OrderResult, sdk.CacheMultiStore) {

	cacheItem := context.MultiStore().CacheMultiStore()
	ctx := context.WithMultiStore(cacheItem)

	// Check order
	msg := MsgCancelOrder{
		Sender:  sender,
		OrderID: orderID,
	}
	validateResult := validateCancelOrder(ctx, k, msg)
	var message string

	if !validateResult.IsOK() {
		message = validateResult.Log
	} else {
		// cancel order
		order := k.GetOrder(ctx, orderID)
		fee := k.CancelOrder(ctx, order, logger)
		message = fee.String()
	}

	cancelRes := types.OrderResult{
		Code:    validateResult.Code,
		Message: message,
		OrderID: orderID,
	}

	return cancelRes, cacheItem
}

func handleMsgCancelOrders(ctx sdk.Context, k Keeper, msg types.MsgCancelOrders, logger log.Logger) sdk.Result {
	cancelRes := []types.OrderResult{}
	for _, orderID := range msg.OrderIDs {

		res, cacheItem := handleCancelOrder(ctx, k, msg.Sender, orderID, logger)
		cancelRes = append(cancelRes, res)
		cacheItem.Write()

		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"    msg<Sender:%s,ID:%s>\n"+
			"    result<The User have canceled an order {ID:%s} >\n",
			ctx.BlockHeight(), "handleMsgCancelOrder",
			msg.Sender, orderID, orderID))

	}
	rss, err := json.Marshal(&cancelRes)
	if err != nil {
		rss = []byte(fmt.Sprintf("failed to marshal result to JSON: %s", err))
	}

	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))
	event = event.AppendAttributes(sdk.NewAttribute("orders", string(rss)))
	ctx.EventManager().EmitEvent(event)
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func validateCancelOrder(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCancelOrder) sdk.Result {
	order := keeper.GetOrder(ctx, msg.OrderID)

	// Check order
	if order == nil {
		return sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  fmt.Sprintf("order(%s) does not exist or already closed", msg.OrderID),
		}
	}
	if order.Status != types.OrderStatusOpen {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("cannot cancel order with status(%d)", order.Status),
		}
	}
	if !order.Sender.Equals(msg.Sender) {
		return sdk.Result{
			Code: sdk.CodeUnauthorized,
			Log:  fmt.Sprintf("not the owner of order(%v)", msg.OrderID),
		}
	}
	if keeper.IsProductLocked(order.Product) {
		return sdk.Result{
			Code: sdk.CodeInternal,
			Log:  fmt.Sprintf("the trading pair (%s) is locked, please retry later", order.Product),
		}
	}
	return sdk.Result{}
}

// ValidateMsgCancelOrders validates whether the msg of cancelOrders is valid.
func ValidateMsgCancelOrders(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCancelOrders) sdk.Result {
	for _, orderID := range msg.OrderIDs {
		msg := MsgCancelOrder{
			Sender:  msg.Sender,
			OrderID: orderID,
		}
		res := validateCancelOrder(ctx, keeper, msg)
		if sdk.CodeOK != res.Code {
			return res
		}
	}

	return sdk.Result{}
}
