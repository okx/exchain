package farm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

func handleMsgCreatePool(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePool, logger log.Logger) sdk.Result {
	if _, found := k.GetFarmPool(ctx, msg.PoolName); found {
		return types.ErrPoolAlreadyExist(DefaultCodespace, msg.PoolName).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.LockedSymbol); !ok {
		return types.ErrTokenNotExist(DefaultCodespace, msg.LockedSymbol).Result()
	}

	yieldTokenInfo := k.TokenKeeper().GetTokenInfo(ctx, msg.YieldedSymbol)
	if !yieldTokenInfo.Owner.Equals(msg.Owner) {
		return types.ErrInvalidTokenOwner(DefaultCodespace, msg.Owner.String(), msg.YieldedSymbol).Result()
	}

	// fee
	feeAmount := k.GetParams(ctx).CreatePoolFee
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, k.GetFeeCollector(), feeAmount.ToCoins()); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeAmount.String())).Result()
	}

	// deposit
	depositAmount := k.GetParams(ctx).CreatePoolDeposit
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, ModuleName, depositAmount.ToCoins()); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			depositAmount.String())).Result()
	}

	// create pool
	yieldedTokenInfo := types.NewYieldedTokenInfo(sdk.NewDecCoin(msg.YieldedSymbol, sdk.ZeroInt()),
		0, sdk.ZeroDec())
	pool := types.NewFarmPool(
		msg.Owner, msg.PoolName, msg.LockedSymbol, depositAmount, sdk.NewDecCoin(msg.LockedSymbol, sdk.ZeroInt()),
		[]types.YieldedTokenInfo{yieldedTokenInfo}, sdk.DecCoins{},
	)
	// set pool
	k.SetFarmPool(ctx, pool)

	// initial pool period
	poolHistoricalRewards := types.NewPoolHistoricalRewards(sdk.DecCoins{}, 1)
	k.SetPoolHistoricalRewards(ctx, msg.PoolName, 0, poolHistoricalRewards)
	PoolCurrentRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), 1, sdk.DecCoins{})
	k.SetPoolCurrentRewards(ctx, msg.PoolName, PoolCurrentRewards)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreatePool,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(types.AttributeKeyLockToken, msg.LockedSymbol),
		sdk.NewAttribute(types.AttributeKeyYieldToken, msg.YieldedSymbol),
		sdk.NewAttribute(sdk.AttributeKeyFee, feeAmount.String()),
		sdk.NewAttribute(types.AttributeKeyDeposit, depositAmount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDestroyPool(ctx sdk.Context, k keeper.Keeper, msg types.MsgDestroyPool, logger log.Logger) sdk.Result {
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}

	if !pool.Owner.Equals(msg.Owner) {
		return types.ErrInvalidPoolOwner(DefaultCodespace, msg.Owner.String(), msg.PoolName).Result()
	}

	if !pool.Finished() {
		return types.ErrPoolNotFinished(DefaultCodespace, msg.PoolName).Result()
	}

	// withdraw
	withdrawAmount := pool.DepositAmount
	if withdrawAmount.IsPositive() {
		err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Owner, withdrawAmount.ToCoins())
		if err != nil {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
				withdrawAmount.String())).Result()
		}
	}

	if !pool.TotalAccumulatedRewards.IsZero() {
		err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, YieldFarmingAccount, msg.Owner, pool.TotalAccumulatedRewards)
		if err != nil {
			return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient rewards coins(need %s)",
				pool.TotalAccumulatedRewards.String())).Result()
		}
	}

	// delete pool
	k.DeleteFarmPool(ctx, msg.PoolName)

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
