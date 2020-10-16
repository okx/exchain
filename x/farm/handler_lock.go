package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock, logger log.Logger) sdk.Result {
	// 0.1 Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	if pool.SymbolLocked != msg.Amount.Denom {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}

	// 0.2 Get the current period
	currentPeriod := k.GetPoolCurrentRewards(ctx, msg.PoolName)

	// 0.3 Get the lock info
	if lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); !found {
		// 1. If lock info doesn't exist, only initialize the LockInfo structure
		// TODO how to init the new lock info, use currentPeriod.Period or currentPeriod.Period+1
		lockInfo = types.NewLockInfo(msg.Address, msg.PoolName, msg.Amount, ctx.BlockHeight(), currentPeriod.Period)
		k.SetLockInfo(ctx, lockInfo)

		// TODO update period?

		// 2. Update the pool info
		pool.TotalValueLocked = pool.TotalValueLocked.Add(msg.Amount)
		k.SetFarmPool(ctx, pool)
	} else {
		// 1. TODO
		_, err := k.WithdrawRewards(ctx, pool, msg.Address)
		if err != nil {
			return err.Result()
		}

		// 2. Reinitialize the lock info
		k.InitializeLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount)
	}

	// 3. Send the locked-tokens from its own account to farm module account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, ModuleName, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

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

	// 2. Claim
	_, err := k.WithdrawRewards(ctx, pool, msg.Address)
	if err != nil {
		return err.Result()
	}

	// 3. Reinitialize the lock info
	k.InitializeLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount.Neg())

	// 4. Send the locked-tokens from farm module account to its own account
	if err = k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Address, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnlock,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}
