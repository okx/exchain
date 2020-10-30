package farm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

func handleMsgCreatePool(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePool) sdk.Result {
	if _, found := k.GetFarmPool(ctx, msg.PoolName); found {
		return types.ErrPoolAlreadyExist(DefaultCodespace, msg.PoolName).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.MinLockAmount.Denom); !ok {
		return types.ErrTokenNotExist(DefaultCodespace, msg.MinLockAmount.Denom).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.YieldedSymbol); !ok {
		return types.ErrTokenNotExist(DefaultParamspace, msg.YieldedSymbol).Result()
	}

	// fee
	params := k.GetParams(ctx)
	feeAmount := params.CreatePoolFee
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Owner, k.GetFeeCollector(), feeAmount.ToCoins(),
	); err != nil {
		return sdk.ErrInsufficientFee(fmt.Sprintf("insufficient fee coins(need %s)",
			feeAmount.String())).Result()
	}

	// deposit
	depositAmount := params.CreatePoolDeposit
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Owner, ModuleName, depositAmount.ToCoins(),
	); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			depositAmount.String())).Result()
	}

	// create pool
	yieldedTokenInfo := types.NewYieldedTokenInfo(sdk.NewDecCoin(msg.YieldedSymbol, sdk.ZeroInt()),
		0, sdk.ZeroDec())
	pool := types.NewFarmPool(
		msg.Owner, msg.PoolName, msg.MinLockAmount, depositAmount, sdk.NewDecCoin(msg.MinLockAmount.Denom, sdk.ZeroInt()),
		[]types.YieldedTokenInfo{yieldedTokenInfo}, sdk.DecCoins{},
	)
	k.SetFarmPool(ctx, pool)

	// initial pool period
	poolHistoricalRewards := types.NewPoolHistoricalRewards(sdk.DecCoins{}, 1)
	k.SetPoolHistoricalRewards(ctx, msg.PoolName, 0, poolHistoricalRewards)
	poolCurrentRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), 1, sdk.DecCoins{})
	k.SetPoolCurrentRewards(ctx, msg.PoolName, poolCurrentRewards)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreatePool,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(types.AttributeKeyMinLockAmount, msg.MinLockAmount.String()),
		sdk.NewAttribute(types.AttributeKeyYieldToken, msg.YieldedSymbol),
		sdk.NewAttribute(sdk.AttributeKeyFee, feeAmount.String()),
		sdk.NewAttribute(types.AttributeKeyDeposit, depositAmount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDestroyPool(ctx sdk.Context, k keeper.Keeper, msg types.MsgDestroyPool) sdk.Result {
	// 0. check pool and owner
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}

	if !pool.Owner.Equals(msg.Owner) {
		return types.ErrInvalidPoolOwner(DefaultCodespace, msg.Owner.String(), msg.PoolName).Result()
	}

	// 1. calculate how many provided token & native token could be yielded in current period
	updatedPool, _ := k.CalculateAmountYieldedBetween(ctx, pool)

	// 2. check pool status
	if !updatedPool.Finished() {
		return types.ErrPoolNotFinished(DefaultCodespace, msg.PoolName).Result()
	}

	// 3. give remaining rewards to the owner of pool
	if !updatedPool.TotalAccumulatedRewards.IsZero() {
		err := k.SupplyKeeper().SendCoinsFromModuleToAccount(
			ctx, YieldFarmingAccount, msg.Owner, updatedPool.TotalAccumulatedRewards,
		)
		if err != nil {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient rewards coins(need %s)",
				updatedPool.TotalAccumulatedRewards.String())).Result()
		}
	}

	// 4. withdraw deposit
	withdrawAmount := pool.DepositAmount
	if withdrawAmount.IsPositive() {
		err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Owner, withdrawAmount.ToCoins())
		if err != nil {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
				withdrawAmount.String())).Result()
		}
	}

	// 5. delete pool and white list
	k.DeleteFarmPool(ctx, msg.PoolName)

	// 6. delete historical period rewards and current period rewards.
	k.IteratePoolHistoricalRewards(ctx, msg.PoolName,
		func(store sdk.KVStore, key []byte, value []byte) (stop bool) {
			store.Delete(key)
			return false
		},
	)
	k.DeletePoolCurrentRewards(ctx, msg.PoolName)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDestroyPool,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(types.AttributeKeyWithdraw, withdrawAmount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}
