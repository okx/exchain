package farm

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
)

func handleMsgCreatePool(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePool) (*sdk.Result, error) {
	if _, found := k.GetFarmPool(ctx, msg.PoolName); found {
		return types.ErrPoolAlreadyExist(msg.PoolName).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.MinLockAmount.Denom); !ok {
		return types.ErrTokenNotExist(msg.MinLockAmount.Denom).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.YieldedSymbol); !ok {
		return types.ErrTokenNotExist(msg.YieldedSymbol).Result()
	}

	// fee
	params := k.GetParams(ctx)
	feeAmount := params.CreatePoolFee
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Owner, k.GetFeeCollector(), feeAmount.ToCoins(),
	); err != nil {
		return nil, common.ErrInsufficientCoins(DefaultParamspace, err.Error())
	}

	// deposit
	depositAmount := params.CreatePoolDeposit
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Owner, ModuleName, depositAmount.ToCoins(),
	); err != nil {
		return nil, common.ErrInsufficientCoins(DefaultParamspace, err.Error())
	}

	// create pool
	yieldedTokenInfo := types.NewYieldedTokenInfo(sdk.NewDecCoin(msg.YieldedSymbol, sdk.ZeroInt()),
		0, sdk.ZeroDec())
	pool := types.NewFarmPool(
		msg.Owner, msg.PoolName, msg.MinLockAmount, depositAmount, sdk.NewDecCoin(msg.MinLockAmount.Denom, sdk.ZeroInt()),
		[]types.YieldedTokenInfo{yieldedTokenInfo}, sdk.SysCoins{},
	)
	k.SetFarmPool(ctx, pool)

	// initial pool period
	poolHistoricalRewards := types.NewPoolHistoricalRewards(sdk.SysCoins{}, 1)
	k.SetPoolHistoricalRewards(ctx, msg.PoolName, 0, poolHistoricalRewards)
	poolCurrentRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), 1, sdk.SysCoins{})
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
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDestroyPool(ctx sdk.Context, k keeper.Keeper, msg types.MsgDestroyPool) (*sdk.Result, error) {
	evmKeeper := k.EvmKeeper()
	err := evmKeeper.ForEachStorage(ctx, msg.Contract, func(key, value ethcmn.Hash) bool {
		evmKeeper.DeleteStateDirectly(ctx, msg.Contract, key)
		return false // todo: need to add a judgement, in case of deleting too many keys in one transaction
	})
	if err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}
