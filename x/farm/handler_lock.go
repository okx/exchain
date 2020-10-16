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

	// 0.2 Get the lock info
	if _, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); !found {
		k.IncrementPoolPeriod(ctx, pool)
		lockInfo := types.NewLockInfo(
			msg.Address, pool.Name, sdk.NewDecCoinFromDec(pool.SymbolLocked, sdk.ZeroDec()),
			ctx.BlockHeight(), 0,
		)
		k.SetLockInfo(ctx, lockInfo)
	} else {
		_, err := k.WithdrawRewards(ctx, pool, msg.Address)
		if err != nil {
			return err.Result()
		}
	}

	// 2. Reinitialize the lock info
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount)

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
	k.UpdateLockInfo(ctx, msg.Address, msg.PoolName, msg.Amount.Amount.Neg())

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
