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

	if ok := k.TokenKeeper().TokenExist(ctx, msg.SymbolLocked); !ok {
		return types.ErrTokenNotExist(DefaultCodespace, msg.SymbolLocked).Result()
	}

	yieldTokenInfo := k.TokenKeeper().GetTokenInfo(ctx, msg.YieldSymbol)
	if !yieldTokenInfo.Owner.Equals(msg.Owner) {
		return types.ErrInvalidTokenOwner(DefaultCodespace, msg.Owner.String(), msg.YieldSymbol).Result()
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
	yieldedTokenInfo := types.NewYieldedTokenInfo(sdk.NewDecCoin(msg.YieldSymbol, sdk.ZeroInt()), 0, sdk.ZeroDec())
	pool := types.FarmPool{
		Owner:             msg.Owner,
		Name:              msg.PoolName,
		SymbolLocked:      msg.SymbolLocked,
		YieldedTokenInfos: []types.YieldedTokenInfo{yieldedTokenInfo},
		DepositAmount:     depositAmount,
		TotalValueLocked:  sdk.DecCoin{Denom: msg.SymbolLocked, Amount: sdk.ZeroDec()},
	}
	// set pool
	k.SetFarmPool(ctx, pool)

	// initial pool period
	poolHistoricalRewards := types.NewPoolHistoricalRewards(sdk.DecCoins{}, 1)
	k.SetPoolHistoricalRewards(ctx, msg.PoolName, 0, poolHistoricalRewards)
	poolCurrentPeriod := types.NewPoolCurrentPeriod(ctx.BlockHeight(), 1, sdk.DecCoin{})
	k.SetPoolCurrentPeriod(ctx, msg.PoolName, poolCurrentPeriod)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreatePool,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(types.AttributeKeyLockToken, msg.SymbolLocked),
		sdk.NewAttribute(types.AttributeKeyYieldToken, msg.YieldSymbol),
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

	// delete pool
	k.DeleteFarmPool(ctx, msg.PoolName)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDestroyPool,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(types.AttributeKeyWithdraw, withdrawAmount.String()),
	))
	return sdk.Result{Events: ctx.EventManager().Events()}
}
