package farm

import (
	"fmt"
	"strconv"

	"github.com/okex/exchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
)

// NewHandler creates an sdk.Handler for all the farm type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var handlerFun func() (*sdk.Result, error)
		var name string
		switch msg := msg.(type) {
		case types.MsgCreatePool:
			name = "handleMsgCreatePool"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgCreatePool(ctx, k, msg)
			}
		case types.MsgDestroyPool:
			name = "handleMsgDestroyPool"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgDestroyPool(ctx, k, msg)
			}
		case types.MsgProvide:
			name = "handleMsgProvide"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgProvide(ctx, k, msg)
			}
		case types.MsgLock:
			name = "handleMsgLock"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgLock(ctx, k, msg)
			}
		case types.MsgUnlock:
			name = "handleMsgUnlock"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgUnlock(ctx, k, msg)
			}
		case types.MsgClaim:
			name = "handleMsgClaim"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgClaim(ctx, k, msg)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return types.ErrUnknownFarmMsgType(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)

		res, err := handlerFun()
		common.SanityCheckHandler(res, err)
		return res, err
	}
}

func handleMsgProvide(ctx sdk.Context, k keeper.Keeper, msg types.MsgProvide) (*sdk.Result, error) {
	// 0. Check if the start_height_to_yield is more than current height
	if msg.StartHeightToYield <= ctx.BlockHeight() {
		return types.ErrInvalidStartHeight().Result()
	}

	// 1.1 Check farm pool
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(msg.PoolName).Result()
	}

	// 1.2 Check owner
	if !pool.Owner.Equals(msg.Address) {
		return types.ErrInvalidPoolOwner(msg.Address.String(), msg.PoolName).Result()
	}

	// 1.3 Check if the provided coin denom is the same as the locked coin name
	if len(pool.YieldedTokenInfos) != 1 {
		panic(fmt.Sprintf("The YieldedTokenInfos length is %d, which should be 1 in current code version",
			len(pool.YieldedTokenInfos)))
	}
	if pool.YieldedTokenInfos[0].RemainingAmount.Denom != msg.Amount.Denom {
		return types.ErrInvalidDenom(pool.YieldedTokenInfos[0].RemainingAmount.Denom, msg.Amount.Denom).Result()
	}

	// 2.1 Calculate how many provided token & native token could be yielded in current period
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	// 2.2 Check if remaining amount is already zero
	remainingAmount := updatedPool.YieldedTokenInfos[0].RemainingAmount
	if !remainingAmount.IsZero() {
		return types.ErrRemainingAmountNotZero(remainingAmount.String()).Result()
	}

	// 3. Terminate pool current period
	k.IncrementPoolPeriod(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens)

	// 4. Transfer coin to farm module account
	if err := k.SupplyKeeper().SendCoinsFromAccountToModule(
		ctx, msg.Address, YieldFarmingAccount, msg.Amount.ToCoins(),
	); err != nil {
		return types.ErrSendCoinsFromAccountToModuleFailed(err.Error()).Result()
	}

	// 5. init a new yielded_token_info struct, then set it into store
	updatedPool.YieldedTokenInfos[0] = types.NewYieldedTokenInfo(
		msg.Amount, msg.StartHeightToYield, msg.AmountYieldedPerBlock,
	)
	k.SetFarmPool(ctx, updatedPool)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeProvide,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyStartHeightToYield, strconv.FormatInt(msg.StartHeightToYield, 10)),
		sdk.NewAttribute(types.AttributeKeyAmountYieldPerBlock, msg.AmountYieldedPerBlock.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgClaim(ctx sdk.Context, k keeper.Keeper, msg types.MsgClaim) (*sdk.Result, error) {
	// 1. Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(msg.PoolName).Result()
	}

	// 2. Calculate how many provided token & native token could be yielded in current period
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	// 3. Withdraw rewards
	rewards, err := k.WithdrawRewards(ctx, pool.Name, pool.TotalValueLocked, yieldedTokens, msg.Address)
	if err != nil {
		return nil, err
	}

	// 4. Update the lock_info data
	k.UpdateLockInfo(ctx, msg.Address, pool.Name, sdk.ZeroDec())

	// 5. Update farm pool
	if updatedPool.TotalAccumulatedRewards.IsAllLT(rewards) {
		panic("should not happen")
	}
	updatedPool.TotalAccumulatedRewards = updatedPool.TotalAccumulatedRewards.Sub(rewards)
	k.SetFarmPool(ctx, updatedPool)

	// 6. notify backend
	k.OnClaim(ctx, msg.Address, pool.Name, rewards)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeClaim,
		sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
