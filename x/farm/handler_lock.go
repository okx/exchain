package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock, logger log.Logger) sdk.Result {
	// 1. Get farm pool
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	if pool.SymbolLocked != msg.Amount.Denom {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}


	// 2. Calculate how many provided token & native token have been yielded between start_block_height and current_height
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)
	updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Add(yieldedTokens)

	// 3. Get lock info
	if _, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); found {
		// If it exists, withdraw money
		rewards, err := k.WithdrawRewards(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens, msg.Address)
		if err != nil {
			return err.Result()
		}
		if updatedPool.TotalAccumulatedRewards.IsAllLT(rewards) {
			panic("should not happen")
		}
		updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Sub(rewards)
	} else {
		// If it doesn't exist, only increase period
		k.IncrementPoolPeriod(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens)

		// Create new lock info
		lockInfo := types.NewLockInfo(
			msg.Address, pool.Name, sdk.NewDecCoinFromDec(pool.SymbolLocked, sdk.ZeroDec()),
			ctx.BlockHeight(), 0,
		)
		k.SetLockInfo(ctx, lockInfo)
		k.SetAddressInFarmPool(ctx, msg.PoolName, msg.Address)
	}

	// 4. Update lock info
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount)

	// 5. Send the locked-tokens from its own account to farm module account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, ModuleName, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

	// 6. Update farm pool
	updatedPool.TotalValueLocked = updatedPool.TotalValueLocked.Add(msg.Amount)
	k.SetFarmPool(ctx, updatedPool)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeLock,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgUnlock(ctx sdk.Context, k keeper.Keeper, msg types.MsgUnlock, logger log.Logger) sdk.Result {
	// 1.1 Check if there are enough tokens to unlock
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, msg.Address.String()).Result()
	}

	if lockInfo.Amount.IsLT(msg.Amount) {
		return types.ErrinsufficientAmount(DefaultCodespace, lockInfo.Amount.String(), msg.Amount.String()).Result()
	}

	// 1.2 Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	if pool.SymbolLocked != msg.Amount.Denom {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}

	// 2. Calculate how many provided token & native token have been yielded between start_block_height and current_height
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	// 3. Withdraw money
	rewards, err := k.WithdrawRewards(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens, msg.Address)
	if err != nil {
		return err.Result()
	}

	// 4. Update the lock info
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount.Neg())

	// 5. Send the locked-tokens from farm module account to its own account
	if err = k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Address, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

	// 6. Update farm pool
	updatedPool.TotalValueLocked = updatedPool.TotalValueLocked.Sub(msg.Amount)
	updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Add(yieldedTokens)
	if updatedPool.TotalAccumulatedRewards.IsAllLT(rewards) {
		panic("should not happen")
	}
	updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Sub(rewards)
	k.SetFarmPool(ctx, updatedPool)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnlock,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}
