package margin

import (
	"fmt"

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
			k.SetAccount(ctx, borrowInfo.Address, account)

			return false
		})

	// force liquidation
}
