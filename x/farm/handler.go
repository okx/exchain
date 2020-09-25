package farm

import (
	"fmt"
	"strconv"

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
		case types.MsgProvide:
			name = "handleMsgProvide"
			handlerFun = func() sdk.Result {
				return handleMsgProvide(ctx, k, msg, logger)
			}
		case types.MsgLock:
			name = "handleMsgLock"
			handlerFun = func() sdk.Result {
				return handleMsgLock(ctx, k, msg, logger)
			}
		case types.MsgUnlock:
			name = "handleMsgUnlock"
			handlerFun = func() sdk.Result {
				return handleMsgUnlock(ctx, k, msg, logger)
			}
		case types.MsgClaim:
			name = "handleMsgClaim"
			handlerFun = func() sdk.Result {
				return handleMsgClaim(ctx, k, msg, logger)
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

func handleMsgProvide(ctx sdk.Context, k keeper.Keeper, msg types.MsgProvide, logger log.Logger) sdk.Result {
	// 0. Check if the start_height_to_yield is more than current height
	if msg.StartHeightToYield <= ctx.BlockHeight() {
		return types.ErrInvalidStartHeight(DefaultCodespace).Result()
	}
	// 0. Check if this address is the owner of the token
	tokenInfo := k.TokenKeeper().GetTokenInfo(ctx, msg.Amount.Denom)
	if !tokenInfo.Owner.Equals(msg.Address) {
		return types.ErrInvalidTokenOwner(DefaultCodespace, msg.Address.String(), msg.Amount.Denom).Result()
	}
	// 0. Get the pool info
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}

	// 1. append yielding_coin into pool
	yieldingCoin := types.YieldingCoin{
		Coin:                    msg.Amount,
		StartBlockHeightToYield: msg.StartHeightToYield,
		YieldAmountPerBlock:     msg.YieldPerBlock,
	}
	pool.YieldingCoins = append(pool.YieldingCoins, yieldingCoin)

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeProvide,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyStartHeightToYield, strconv.FormatInt(msg.StartHeightToYield, 10)),
			sdk.NewAttribute(types.AttributeKeyYiledPerBlock, msg.YieldPerBlock.String()),
		),
	}}
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

	// Send the token from its own account to farm account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, ModuleName, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

	// Get the pool info, then update the total coin & weight
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound("", msg.PoolName).Result()
	}
	pool.TotalLockedCoin = pool.TotalLockedCoin.Add(msg.Amount)
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(sdk.NewDec(ctx.BlockHeight()).MulTruncate(msg.Amount.Amount))
	k.SetFarmPool(ctx, pool)

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	}}
}

func handleMsgUnlock(ctx sdk.Context, k keeper.Keeper, msg types.MsgUnlock, logger log.Logger) sdk.Result {
	// Get the specific lock_info.
	// If it doesn't exist, just return.
	// Otherwise, calculate the previous liquidity mining reward at first
	// Then update the lock_info
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, msg.Address.String()).Result()
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

	// Send the token from farm account to its own account
	if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Address, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

	// Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	pool.TotalLockedCoin = pool.TotalLockedCoin.Sub(msg.Amount)
	pool.TotalLockedWeight = pool.TotalLockedWeight.Sub(sdk.NewDec(ctx.BlockHeight()).MulTruncate(msg.Amount.Amount))
	k.SetFarmPool(ctx, pool)

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	}}
}

func handleMsgClaim(ctx sdk.Context, k keeper.Keeper, msg types.MsgClaim, logger log.Logger) sdk.Result {
	err := claim(ctx, k, msg.PoolName, msg.Address)
	if err != nil {
		return err.Result()
	}

	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, msg.Address.String()).Result()
	}
	lockInfo.StartBlockHeight = ctx.BlockHeight()
	k.SetLockInfo(ctx, lockInfo)

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		),
	}}
}

func claim(ctx sdk.Context, k keeper.Keeper, poolName string, address sdk.AccAddress) sdk.Error {
	// 0. get the lock info
	lockInfo, found := k.GetLockInfo(ctx, address, poolName)
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, address.String())
	}

	// 0. get the pool info
	pool, poolFound := k.GetFarmPool(ctx, poolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, poolName)
	}

	height := ctx.BlockHeight()
	if height < pool.LastYieldedBlockHeight {
		return nil
	}

	// TODO there are too many operations about MulTruncate, check the amount carefully!!!
	// TODO rename parameters?
	// 1. Transfer yileding_coin -> yileded_coin
	yieldedCoins := sdk.DecCoins{}
	for i := 0; i < len(pool.YieldingCoins); i++ {
		if height >= pool.YieldingCoins[i].StartBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if pool.YieldingCoins[i].StartBlockHeightToYield > pool.LastYieldedBlockHeight {
				blockInterval = sdk.NewDec(height - pool.YieldingCoins[i].StartBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(height - pool.LastYieldedBlockHeight)
			}

			// calculate how many coin have been yileded till the current block
			yieldedAmount := blockInterval.MulTruncate(pool.YieldingCoins[i].YieldAmountPerBlock)
			yieldedCoins = yieldedCoins.Add(sdk.NewDecCoinsFromDec(pool.YieldingCoins[i].Coin.Denom, yieldedAmount))

			// subtract yileding_coin amount
			pool.YieldingCoins[i].Coin.Amount = pool.YieldingCoins[i].Coin.Amount.Sub(yieldedAmount)

			// TODO what if pool.YieldingCoins[i].Coin become zero, or less than YieldAmountPerBlock
		}
	}
	pool.YieldedCoins = pool.YieldedCoins.Add(yieldedCoins)
	pool.LastYieldedBlockHeight = height

	/* 2.1 Calculate its own weight during these blocks
	   (curHeight - Height1) * Amount1
	*/
	// TODO: is there any possibility that lockInfo.StartBlockHeight is more than ctx.BlockHeight()?
	currentWeight := sdk.NewDec(lockInfo.StartBlockHeight).MulTruncate(lockInfo.Amount.Amount)
	oldWeight := sdk.NewDec(height).MulTruncate(lockInfo.Amount.Amount)
	numerator := currentWeight.Sub(oldWeight)

	/* 2.2 Calculate all weight during these blocks
	    (curHeight - Height1) * Amount1 + (curHeight - Height2) * Amount2 + (curHeight - Height3) * Amount3
												||
	                                            \/
	   curHeight * (Amount1 + Amount2 + Amount3) - (Height1*Amount1 + Height2*Amount2 + Height3*Amount3)
												||
	                                            \/
	ctx.BlockHeight()  *  pool.TotalLockedCoin.Amount  -  ( pool.TotalLockedWeight )
	*/
	denominator := sdk.NewDec(height).MulTruncate(pool.TotalLockedCoin.Amount).Sub(pool.TotalLockedWeight)

	// 2.3 Calculate how many yielded token it could get
	selfYieldedCoins := pool.YieldedCoins.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	// 2.4 Transfer yielded tokens
	if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, address, selfYieldedCoins); err != nil {
		return err
	}

	// 3 Update pool data
	pool.YieldedCoins = pool.YieldedCoins.Sub(selfYieldedCoins)
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(numerator)
	// set pool into store
	k.SetFarmPool(ctx, pool)

	return nil
}
