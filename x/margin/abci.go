package margin

import (
	"fmt"
	"sync"

	"github.com/okex/okchain/x/order"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/margin/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
	seq := perf.GetPerf().OnEndBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnEndBlockExit(ctx, types.ModuleName, seq)

	// complete dex withdraw
	currentTime := ctx.BlockHeader().Time
	k.IterateDexWithdrawAddress(ctx, currentTime,
		func(_ int64, key []byte) (stop bool) {
			oldTime, addr := types.SplitWithdrawTimeKey(key)
			err := k.CompleteWithdraw(ctx, addr)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("complete withdraw failed: %s, addr:%s", err.Error(), addr.String()))
			} else {
				ctx.Logger().Debug(fmt.Sprintf("complete withdraw successful, addr: %s", addr.String()))
				k.DeleteDexWithdrawCompleteTimeAddress(ctx, oldTime, addr)
			}
			return false
		})

	// calculate interest
	k.IterateCalculateInterest(ctx, currentTime,
		func(limit int64, key []byte) (stop bool) {
			// execute limit per block
			if limit >= types.CalculateInterestLimitPerBlock {
				return true
			}
			oldTime, borrowInfoKey := types.SplitCalculateInterestTimeKey(key)
			borrowInfo := k.GetBorrowInfoByKey(ctx, borrowInfoKey)
			if borrowInfo == nil {
				// delete repaid key
				k.DeleteCalculateInterestKey(ctx, oldTime, borrowInfoKey)
				return false
			}

			// calculate interest
			interest := sdk.DecCoins{}
			for _, borrowDecCoin := range borrowInfo.BorrowAmount {
				interest = interest.Add(sdk.NewCoins(borrowDecCoin).MulDec(borrowInfo.Rate))
			}

			// update next calculate time key
			k.SetCalculateInterestKey(ctx, oldTime.Add(k.GetParams(ctx).InterestPeriod), borrowInfo.Address,
				borrowInfo.Product, uint64(borrowInfo.BlockHeight))
			k.DeleteCalculateInterestKey(ctx, oldTime, borrowInfoKey)

			// update account to db
			account := k.GetAccount(ctx, borrowInfo.Address, borrowInfo.Product)
			if account == nil {
				k.Logger(ctx).Error(fmt.Sprintf("unexpected error margin account not exists:address=(%s),product=(%s)",
					borrowInfo.Address.String(), borrowInfo.Product))
				return false
			}
			account.Interest = account.Interest.Add(interest)
			k.SetAccount(ctx, borrowInfo.Address, borrowInfo.Product, account)

			return false
		})

	// liquidation
	var wg sync.WaitGroup
	tradePairs := k.GetAllTradePairs(ctx)
	for _, tradePair := range tradePairs {
		// get all address borrowed on tradePair
		addressList := k.GetBorrowedAddress(ctx, tradePair.Name)
		if len(addressList) > 0 {
			wg.Add(1)
			go doLiquidation(&wg, ctx, k, tradePair, addressList)
		}
	}
	wg.Wait()
}

func doLiquidation(wg *sync.WaitGroup, ctx sdk.Context, k Keeper, tradePair *types.TradePair, addressList []sdk.AccAddress) {
	defer wg.Done()
	for _, address := range addressList {
		account := k.GetAccount(ctx, address, tradePair.Name)
		latestPrice := k.GetOrderKeeper().GetLastPrice(ctx, tradePair.Name)
		baseSymbol := tradePair.BaseSymbol()
		quoteSymbol := tradePair.QuoteSymbol()
		marginRatio := account.MarginRatio(baseSymbol, quoteSymbol, latestPrice)
		// do not need liquidation
		if marginRatio.GT(tradePair.MaintenanceMarginRatio) {
			continue
		}

		// force liquidation
		if account.Borrowed.AmountOf(baseSymbol).IsPositive() && account.Available.AmountOf(baseSymbol).IsPositive() {
			k.Repay(ctx, account, address, tradePair, sdk.NewDecCoinFromDec(baseSymbol, account.Available.AmountOf(baseSymbol)))
			k.Logger(ctx).Debug(fmt.Sprintf("force liquidation, repay:address(%s) product(%s) amount(%s)",
				address, tradePair.Name, sdk.NewDecCoinFromDec(baseSymbol, account.Available.AmountOf(baseSymbol))))
		}
		if account.Borrowed.AmountOf(quoteSymbol).IsPositive() && account.Available.AmountOf(quoteSymbol).IsPositive() {
			k.Repay(ctx, account, address, tradePair, sdk.NewDecCoinFromDec(quoteSymbol, account.Available.AmountOf(quoteSymbol)))
			k.Logger(ctx).Debug(fmt.Sprintf("force liquidation, repay:address(%s) product(%s) amount(%s)",
				address, tradePair.Name, sdk.NewDecCoinFromDec(quoteSymbol, account.Available.AmountOf(quoteSymbol))))
		}

		// force placing order
		placeOrder(ctx, k, tradePair, latestPrice, address, account)
	}
}

func placeOrder(ctx sdk.Context, k Keeper, tradePair *types.TradePair, latestPrice sdk.Dec, address sdk.AccAddress, account *types.Account) {
	if k.GetOrderKeeper().IsProductLocked(ctx, tradePair.Name) {
		return
	}
	baseSymbol := tradePair.BaseSymbol()
	quoteSymbol := tradePair.QuoteSymbol()

	bestBid, bestAsk := k.GetOrderKeeper().GetBestBidAndAsk(ctx, tradePair.Name)
	if bestBid.IsZero() {
		bestBid = latestPrice
	}
	if bestAsk.IsZero() {
		bestAsk = latestPrice
	}
	bestBid = bestBid.Mul(sdk.MustNewDecFromStr("1.05"))
	bestAsk = bestAsk.Mul(sdk.MustNewDecFromStr("0.95"))

	baseAvailable := account.Available.AmountOf(baseSymbol)
	baseBorrowed := account.Borrowed.AmountOf(baseSymbol)
	quoteAvailable := account.Available.AmountOf(quoteSymbol)
	quoteBorrowed := account.Borrowed.AmountOf(quoteSymbol)

	// place buy order to repay base token
	if baseBorrowed.IsPositive() && quoteAvailable.IsPositive() {
		side := order.BuyOrder
		quantity := baseBorrowed
		maxBuyQuantity := quoteAvailable.Quo(bestBid)
		if quantity.GT(maxBuyQuantity) {
			quantity = maxBuyQuantity
		}
		msg := order.MsgNewOrder{
			Sender:   address,
			Product:  tradePair.Name,
			Side:     side,
			Price:    bestBid,
			Quantity: quantity,
			Type:     order.MarginOrder,
		}
		orderKeeper := k.GetOrderKeeper().(*order.Keeper)
		newOrder := order.GetOrderFromMsg(ctx, *orderKeeper, msg, order.DefaultNewOrderFeeRatio)
		k.GetOrderKeeper().PlaceOrder(ctx, newOrder)
		k.Logger(ctx).Debug(fmt.Sprintf("force liquidation, place order: %+v", msg))
	}

	// place sell order to repay quote token
	if quoteBorrowed.IsPositive() && baseAvailable.IsPositive() {
		side := order.SellOrder
		quantity := quoteBorrowed.Quo(bestAsk)
		maxSellQuantity := baseAvailable
		if quantity.GT(maxSellQuantity) {
			quantity = maxSellQuantity
		}
		msg := order.MsgNewOrder{
			Sender:   address,
			Product:  tradePair.Name,
			Side:     side,
			Price:    bestBid,
			Quantity: quantity,
			Type:     order.MarginOrder,
		}
		orderKeeper := k.GetOrderKeeper().(*order.Keeper)
		newOrder := order.GetOrderFromMsg(ctx, *orderKeeper, msg, order.DefaultNewOrderFeeRatio)
		k.GetOrderKeeper().PlaceOrder(ctx, newOrder)
		k.Logger(ctx).Debug(fmt.Sprintf("force liquidation, place order: %+v", msg))
	}
}
