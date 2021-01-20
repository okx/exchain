package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock) (*sdk.Result, error) {
	// 1.1 Get farm pool
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(msg.PoolName).Result()
	}
	if pool.MinLockAmount.Denom != msg.Amount.Denom {
		return types.ErrInvalidDenom(pool.MinLockAmount.Denom, msg.Amount.Denom).Result()
	}

	// 1.2. check min lock amount
	hasLocked := k.HasLockInfo(ctx, msg.Address, msg.PoolName)
	if !hasLocked && msg.Amount.Amount.LT(pool.MinLockAmount.Amount) {
		return types.ErrLockAmountBelowMinimum(pool.MinLockAmount.Amount, msg.Amount.Amount).Result()
	}

	// 2. Calculate how many provided token & native token could be yielded in current period
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	// 3. Lock info
	var rewards sdk.SysCoins
	if hasLocked {
		// If it exists, withdraw money
		rewards, err := k.WithdrawRewards(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens, msg.Address)
		if err != nil {
			return nil, err
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
			msg.Address, pool.Name, sdk.NewDecCoinFromDec(pool.MinLockAmount.Denom, sdk.ZeroDec()),
			ctx.BlockHeight(), 0,
		)
		k.SetLockInfo(ctx, lockInfo)
		k.SetAddressInFarmPool(ctx, msg.PoolName, msg.Address)
	}

	// 4. Update lock info
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount)

	// 5. Send the locked-tokens from its own account to farm module account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Address, ModuleName, msg.Amount.ToCoins(),
	); err != nil {
		return nil, types.ErrSendCoinsFromAccountToModuleFailed(err.Error())
	}

	// 6. Update farm pool
	updatedPool.TotalValueLocked = updatedPool.TotalValueLocked.Add(msg.Amount)
	k.SetFarmPool(ctx, updatedPool)

	// 7. notify backend
	if hasLocked {
		k.OnClaim(ctx, msg.Address, pool.Name, rewards)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeLock,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUnlock(ctx sdk.Context, k keeper.Keeper, msg types.MsgUnlock) (*sdk.Result, error) {
	// 1.1 Check if there are enough tokens to unlock
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName)
	if !found {
		return types.ErrNoLockInfoFound(msg.Address.String(), msg.PoolName).Result()
	}

	if lockInfo.Amount.Denom != msg.Amount.Denom {
		return types.ErrInvalidDenom(lockInfo.Amount.Denom, msg.Amount.Denom).Result()
	}

	if lockInfo.Amount.IsLT(msg.Amount) {
		return types.ErrInsufficientAmount(lockInfo.Amount.String(), msg.Amount.String()).Result()
	}

	// 1.2 Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(msg.PoolName).Result()
	}
	if pool.MinLockAmount.Denom != msg.Amount.Denom {
		return types.ErrInvalidDenom(pool.MinLockAmount.Denom, msg.Amount.Denom).Result()
	}
	remainAmount := lockInfo.Amount.Amount.Sub(msg.Amount.Amount)
	if !remainAmount.IsZero() && remainAmount.LT(pool.MinLockAmount.Amount) {
		return types.ErrLockAmountBelowMinimum(pool.MinLockAmount.Amount, remainAmount).Result()
	}

	// 2. Calculate how many provided token & native token could be yielded in current period
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	// 3. Withdraw money
	rewards, err := k.WithdrawRewards(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens, msg.Address)
	if err != nil {
		return nil, err
	}

	// 4. Update the lock info
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount.Neg())

	// 5. Send the locked-tokens from farm module account to its own account
	if err = k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Address, msg.Amount.ToCoins()); err != nil {
		return nil, types.ErrSendCoinsFromModuleToAccountFailed(err.Error())
	}

	// 6. Update farm pool
	updatedPool.TotalValueLocked = updatedPool.TotalValueLocked.Sub(msg.Amount)
	if updatedPool.TotalAccumulatedRewards.IsAllLT(rewards) {
		panic("should not happen")
	}
	updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Sub(rewards)
	k.SetFarmPool(ctx, updatedPool)

	// 7. notify backend
	k.OnClaim(ctx, msg.Address, pool.Name, rewards)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnlock,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
