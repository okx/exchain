package order

import (
	"encoding/json"
	"fmt"
	"math"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/order/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/willf/bitset"
)

func CalculateGas(msg sdk.Msg, params *types.Params) (gas uint64) {
	switch msg := msg.(type) {
	case types.MsgNewOrders:
		gas = msg.CalculateGas(params.NewOrderMsgGasUnit)
	case types.MsgCancelOrders:
		gas = msg.CalculateGas(params.CancelOrderMsgGasUnit)
	default:
		gas = math.MaxUint64
	}

	return gas
}

// NewOrderHandler returns the handler with version 0.
func NewOrderHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// disable order tx handler
		if sdk.HigherThanMercury(ctx.BlockHeight()) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "order message has been disabled")
		}

		gas := CalculateGas(msg, keeper.GetParams(ctx))

		// consume gas that msg required, it will panic if gas is insufficient
		ctx.GasMeter().ConsumeGas(gas, storetypes.GasWriteCostFlatDesc)

		if ctx.IsCheckTx() {
			return &sdk.Result{}, nil
		} else {
			// set an infinite gas meter and recovery it when return
			gasMeter := ctx.GasMeter()
			ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			defer func() { ctx = ctx.WithGasMeter(gasMeter) }()
		}

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() (*sdk.Result, error)
		var name string
		logger := ctx.Logger().With("module", "order")
		switch msg := msg.(type) {
		case types.MsgNewOrders:
			name = "handleMsgNewOrders"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgNewOrders(ctx, keeper, msg, logger)
			}
		case types.MsgCancelOrders:
			name = "handleMsgCancelOrders"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgCancelOrders(ctx, keeper, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("Invalid msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)

		res, err := handlerFun()
		common.SanityCheckHandler(res, err)
		return res, err
	}
}

// checkOrderNewMsg: check msg product, price & quantity fields
func checkOrderNewMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgNewOrder) error {
	tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, msg.Product)
	if tokenPair == nil {
		return types.ErrTokenPairNotExist(msg.Product)
	}

	// check if the order is involved with the tokenpair in dex Delist
	isDelisting, err := keeper.GetDexKeeper().CheckTokenPairUnderDexDelist(ctx, msg.Product)
	if err != nil {
		return err
	}
	if isDelisting {
		return types.ErrTradingPairIsDelisting(msg.Product)
	}

	priceDigit := tokenPair.MaxPriceDigit
	quantityDigit := tokenPair.MaxQuantityDigit
	roundedPrice := msg.Price.RoundDecimal(priceDigit)
	roundedQuantity := msg.Quantity.RoundDecimal(quantityDigit)
	if !roundedPrice.Equal(msg.Price) {
		return types.ErrPriceOverAccuracy(msg.Price, priceDigit)
	}
	if !roundedQuantity.Equal(msg.Quantity) {
		return types.ErrQuantityOverAccuracy(msg.Quantity, quantityDigit)
	}

	if msg.Quantity.LT(tokenPair.MinQuantity) {
		return types.ErrMsgQuantityLessThan(tokenPair.MinQuantity.String())
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
	err := checkOrderNewMsg(ctxItem, k, msg)

	if err == nil {
		if k.IsProductLocked(ctx, msg.Product) {
			err = types.ErrIsProductLocked(order.Product)
		} else {
			err = k.PlaceOrder(ctxItem, order)
		}
	}

	res := types.OrderResult{
		Error:   err,
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
	logger log.Logger) (*sdk.Result, error) {
	event := sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName))

	ratio := "1"
	if len(msg.OrderItems) > 1 {
		ratio = "0.8"
	}

	rs := make([]types.OrderResult, 0, len(msg.OrderItems))
	var handlerResult bitset.BitSet
	for idx, item := range msg.OrderItems {
		res, cacheItem, err := handleNewOrder(ctx, k, msg.Sender, item, ratio, logger)
		if err == nil {
			cacheItem.Write()
			handlerResult.Set(uint(idx))
		}
		rs = append(rs, res)
	}
	rss, err := json.Marshal(&rs)
	if err != nil {
		rss = []byte(fmt.Sprintf("failed to marshal result to JSON: %s", err))
	}
	event = event.AppendAttributes(sdk.NewAttribute("orders", string(rss)))
	ctx.EventManager().EmitEvent(event)

	if handlerResult.None() {
		return types.ErrAllOrderFailedToExecute().Result()
	}

	k.AddTxHandlerMsgResult(handlerResult)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// ValidateMsgNewOrders validates whether the msg of newOrders is valid.
func ValidateMsgNewOrders(ctx sdk.Context, k keeper.Keeper, msg types.MsgNewOrders) (*sdk.Result, error) {
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
			return nil, err
		}
		if k.IsProductLocked(ctx, msg.Product) {
			return types.ErrIsProductLocked(msg.Product).Result()
		}

		order := getOrderFromMsg(ctx, k, msg, ratio)
		_, err = k.TryPlaceOrder(ctx, order)
		if err != nil {
			return common.ErrInsufficientCoins(DefaultParamspace, err.Error()).Result()
		}
	}

	return &sdk.Result{}, nil

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
	err := validateCancelOrder(ctx, k, msg)
	var message string

	if err == nil {
		// cancel order
		order := k.GetOrder(ctx, orderID)
		fee := k.CancelOrder(ctx, order, logger)
		message = fee.String()
	}

	cancelRes := types.OrderResult{
		Error:   err,
		Message: message,
		OrderID: orderID,
	}

	return cancelRes, cacheItem
}

func handleMsgCancelOrders(ctx sdk.Context, k Keeper, msg types.MsgCancelOrders, logger log.Logger) (*sdk.Result, error) {
	cancelRes := []types.OrderResult{}
	var handlerResult bitset.BitSet
	for idx, orderID := range msg.OrderIDs {

		res, cacheItem := handleCancelOrder(ctx, k, msg.Sender, orderID, logger)
		cancelRes = append(cancelRes, res)
		cacheItem.Write()
		if res.Error == nil {
			handlerResult.Set(uint(idx))
		}

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

	if handlerResult.None() {
		return types.ErrNoOrdersIsCanceled().Result()
	}

	k.AddTxHandlerMsgResult(handlerResult)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func validateCancelOrder(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCancelOrder) error {
	order := keeper.GetOrder(ctx, msg.OrderID)

	// Check order
	if order == nil {
		return types.ErrOrderIsNotExistOrClosed(msg.OrderID)
	}
	if order.Status != types.OrderStatusOpen {
		return types.ErrOrderStatusIsNotOpen()
	}
	if !order.Sender.Equals(msg.Sender) {
		return types.ErrNotOrderOwner(msg.OrderID)
	}
	if keeper.IsProductLocked(ctx, order.Product) {
		return types.ErrIsProductLocked(order.Product)
	}
	return nil
}

// ValidateMsgCancelOrders validates whether the msg of cancelOrders is valid.
func ValidateMsgCancelOrders(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCancelOrders) error {
	for _, orderID := range msg.OrderIDs {
		msg := MsgCancelOrder{
			Sender:  msg.Sender,
			OrderID: orderID,
		}
		err := validateCancelOrder(ctx, keeper, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
