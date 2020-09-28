package farm

import (
	"strings"

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
	if strings.Compare(pool.SymbolLocked, msg.Amount.Denom) != 0 {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}

	// 0.2 Get the pool info
	if lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); !found {
		// 1. If lock info doesn't exist, only initialize the LockInfo structure
		lockInfo = types.NewLockInfo(msg.Address, msg.PoolName, msg.Amount, ctx.BlockHeight())
		k.SetLockInfo(ctx, lockInfo)
	} else {
		// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
		updatedPool := liquidateYieldTokenInfo(ctx.BlockHeight(), pool)

		// 2. Claim
		err := claim(ctx, k, updatedPool, lockInfo, msg.Address, msg.Amount.Amount)
		if err != nil {
			return err.Result()
		}

		// 3. Update the lock info
		lockInfo.Amount = lockInfo.Amount.Add(msg.Amount)
		lockInfo.StartBlockHeight = ctx.BlockHeight()
		k.SetLockInfo(ctx, lockInfo)
	}

	// 3. Send the locked-tokens from its own account to farm module account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, ModuleName, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

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
	// 0.1 Check if there are enough tokens to unlock
	lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName);
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, msg.Address.String()).Result()
	} else {
		if lockInfo.Amount.IsLT(msg.Amount) {
			return types.ErrinsufficientAmount(DefaultCodespace, lockInfo.Amount.String(), msg.Amount.String()).Result()
		}
	}

	// 0.2 Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	if strings.Compare(pool.SymbolLocked, msg.Amount.Denom) != 0 {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}

	// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
	updatedPool := liquidateYieldTokenInfo(ctx.BlockHeight(), pool)

	// 2. Claim
	err := claim(ctx, k, updatedPool, lockInfo, msg.Address, sdk.ZeroDec().Sub(msg.Amount.Amount))
	if err != nil {
		return err.Result()
	}

	// 3. Update the lock info
	lockInfo.Amount = lockInfo.Amount.Sub(msg.Amount)
	if lockInfo.Amount.IsZero() {
		k.DeleteLockInfo(ctx, lockInfo.Owner, lockInfo.PoolName)
	} else {
		lockInfo.StartBlockHeight = ctx.BlockHeight()
		k.SetLockInfo(ctx, lockInfo)
	}

	// 4. Send the locked-tokens from farm module account to its own account
	if err = k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Address, msg.Amount.ToCoins()); err != nil {
		return err.Result()
	}

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
