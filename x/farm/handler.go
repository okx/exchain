package farm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okexchain/x/common/perf"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

// NewHandler creates an sdk.Handler for all the farm type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", types.ModuleName)

		var handlerFun func() sdk.Result
		var name string
		switch msg := msg.(type) {
		case types.MsgCreatePool:
			name = "handleMsgList"
			handlerFun = func() sdk.Result {
				return handleMsgCreatePool(ctx, k, msg, logger)
			}
		case types.MsgLock:
			name = "handleMsgLockFarm"
			handlerFun = func() sdk.Result {
				return handleMsgLock(ctx, k, msg, logger)
			}
		case types.MsgUnlock:
			name = "handleMsgUnlockFarm"
			handlerFun = func() sdk.Result {
				return handleMsgUnlock(ctx, k, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgCreatePool(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePool, logger log.Logger) sdk.Result {
	return sdk.Result{}
}

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock, logger log.Logger) sdk.Result {
	// Get the specific lock info.
	// If it doesn't exist, only initialize the LockInfo structure.
	// Otherwise, calculate the previous liquidity mining reward at first
	// Then update the lock_info
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		lockInfo.Addr = msg.Address
	} else {
		// excute claim
		err := claim(ctx, k, msg.PoolName, msg.Address)
		if err != nil {
			return err.Result()
		}
	}
	lockInfo.Amount = lockInfo.Amount.Add(msg.Amount)
	lockInfo.StartBlockHeight = ctx.BlockHeight()
	k.SetLockInfo(ctx, lockInfo)

	// Subtract the token from its own account into farm account
	// TODO k.supplyKeeper.SendCoins(xxx)

	// Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound("", msg.PoolName).Result()
	}
	pool.TotalLockedCoin = pool.TotalLockedCoin.Add(msg.Amount)
	pool.TotalLockedWeight= pool.TotalLockedWeight.Add(sdk.NewDec(ctx.BlockHeight()).MulTruncate(msg.Amount.Amount))
	k.SetFarmPool(ctx, pool)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgUnlock(ctx sdk.Context, k keeper.Keeper, msg types.MsgUnlock, logger log.Logger) sdk.Result {
	// Get the specific lock_info.
	// If it doesn't exist, just return.
	// Otherwise, calculate the previous liquidity mining reward at first
	// Then update the lock_info
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		return types.ErrNoLockInfoFound("", msg.Address.String()).Result()
	} else {
		// excute claim
		err := claim(ctx, k, msg.PoolName, msg.Address)
		if err != nil {
			return err.Result()
		}
	}
	lockInfo.Amount = lockInfo.Amount.Sub(msg.Amount)
	lockInfo.StartBlockHeight = ctx.BlockHeight()
	k.SetLockInfo(ctx, lockInfo)

	// Add the token from farm account into its own account
	// TODO k.supplyKeeper.SendCoins(xxx)

	// Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound("", msg.PoolName).Result()
	}
	pool.TotalLockedCoin = pool.TotalLockedCoin.Sub(msg.Amount)
	pool.TotalLockedWeight= pool.TotalLockedWeight.Sub(sdk.NewDec(ctx.BlockHeight()).MulTruncate(msg.Amount.Amount))
	k.SetFarmPool(ctx, pool)

	return sdk.Result{}
}

func claim(ctx sdk.Context, k keeper.Keeper, poolName string, address sdk.AccAddress) sdk.Error {
	// get the lock info
	lockInfo, lockFound := k.GetLockInfo(ctx, address, poolName)
	if !lockFound {
		return types.ErrNoLockInfoFound("", address.String())
	}

	// get the pool info
	pool, poolFound := k.GetFarmPool(ctx, poolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound("", poolName)
	}

	height := ctx.BlockHeight()
	var intervalHeight sdk.Dec
	if height > pool.LastYieldedBlockHeight {
		intervalHeight = sdk.NewDec(height - pool.LastYieldedBlockHeight)
		// yileding_coin -> yileded_coin
		var yieldedCoins sdk.DecCoins
		for _, yieldingCoin := range pool.YieldingCoins {
			if height > yieldingCoin.StartBlockHeightToYield {
				// calculate how many coin have been yileding
				yieldAmt := intervalHeight.MulTruncate(yieldingCoin.YieldAmountPerBlock)
				yieldedCoin := sdk.NewDecCoinFromDec(yieldingCoin.Coin.Denom, yieldAmt)
				yieldedCoins = yieldedCoins.Add(sdk.DecCoins{yieldedCoin})

				// subtract the yielded_coin from yileding_coin
				yieldingCoin.Coin = yieldingCoin.Coin.Sub(yieldedCoin)
			}
		}

		pool.YieldedCoins = pool.YieldedCoins.Add(yieldedCoins)
		pool.LastYieldedBlockHeight = height
	}
	// TODO what if height < pool.LastYieldedBlockHeight
	// TODO there are too many operations about MulTruncate, check the amount carefully!!!
	// TODO rename parameters?

	/* Calculate its own weight during these blocks
	    (curHeight - Height1) * Amount1
	*/
	currentWeight := sdk.NewDec(lockInfo.StartBlockHeight).MulTruncate(lockInfo.Amount.Amount)
	oldWeight := sdk.NewDec(height).MulTruncate(lockInfo.Amount.Amount)
	numerator := currentWeight.Sub(oldWeight)

	/* Calculate all weight during these blocks
	    (curHeight - Height1) * Amount1 + (curHeight - Height2) * Amount2 + (curHeight - Height3) * Amount3
												||
	                                            \/
	   curHeight * (Amount1 + Amount2 + Amount3) - (Height1*Amount1 + Height2*Amount2 + Height3*Amount3)
												||
	                                            \/
	ctx.BlockHeight()  *  pool.TotalLockedCoin.Amount  -  ( pool.TotalLockedWeight )
	*/
	denominator := sdk.NewDec(height).MulTruncate(pool.TotalLockedCoin.Amount).Sub(pool.TotalLockedWeight)

	// Calculate how many yielded token it could get
	selfYieldedCoins := pool.YieldedCoins.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	// Transfer tokens
	// TODO transfer selfYieldedCoins into its own account
	// k.supplyKeeper.XXXX()

	// Update pool data
	pool.YieldedCoins = pool.YieldedCoins.Sub(selfYieldedCoins)
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(numerator)

	// set pool into store
	k.SetFarmPool(ctx, pool)

	return nil
}