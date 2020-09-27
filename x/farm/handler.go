package farm

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okexchain/x/common/perf"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

// NewHandler creates an sdk.Handler for all the farm type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", types.ModuleName)

		var handlerFun func() sdk.Result
		var name string
		switch msg := msg.(type) {
		case types.MsgCreatePool:
			name = "handleMsgCreatePool"
			handlerFun = func() sdk.Result {
				return handleMsgCreatePool(ctx, k, msg, logger)
			}
		case types.MsgDestroyPool:
			name = "handleMsgDestroyPool"
			handlerFun = func() sdk.Result {
				return handleMsgDestroyPool(ctx, k, msg, logger)
			}
		case types.MsgProvide:
			name = "handleMsgProvide"
			handlerFun = func() sdk.Result {
				return handleMsgProvide(ctx, k, msg, logger)
			}
		case types.MsgLock:
			name = "handleMsgLock"
			handlerFun = func() sdk.Result {
				return handleMsgLock(ctx, k, msg, logger)
			}
		case types.MsgUnlock:
			name = "handleMsgUnlock"
			handlerFun = func() sdk.Result {
				return handleMsgUnlock(ctx, k, msg, logger)
			}
		case types.MsgClaim:
			name = "handleMsgClaim"
			handlerFun = func() sdk.Result {
				return handleMsgClaim(ctx, k, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgCreatePool(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePool, logger log.Logger) sdk.Result {
	if _, found := k.GetFarmPool(ctx, msg.PoolName); found {
		return types.ErrPoolAlreadyExist(DefaultCodespace, msg.PoolName).Result()
	}

	if ok := k.TokenKeeper().TokenExist(ctx, msg.LockToken); !ok {
		return types.ErrTokenNotExist(DefaultCodespace, msg.LockToken).Result()
	}

	yieldTokenInfo := k.TokenKeeper().GetTokenInfo(ctx, msg.YieldToken)
	if !yieldTokenInfo.Owner.Equals(msg.Owner) {
		return types.ErrInvalidTokenOwner(DefaultCodespace, msg.Owner.String(), msg.YieldToken).Result()
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
	yieldedTokenInfo := types.NewYieldedTokenInfo(sdk.NewDecCoin(msg.LockToken, sdk.ZeroInt()), 0, sdk.ZeroDec())
	pool := types.FarmPool{
		Owner:             msg.Owner,
		Name:              msg.PoolName,
		SymbolLocked:      msg.LockToken,
		YieldedTokenInfos: []types.YieldedTokenInfo{yieldedTokenInfo},
		DepositAmount:     depositAmount,
	}
	k.SetFarmPool(ctx, pool)

	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(types.AttributeKeyLockToken, msg.LockToken),
			sdk.NewAttribute(types.AttributeKeyYieldToken, msg.YieldToken),
			sdk.NewAttribute(sdk.AttributeKeyFee, feeAmount.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, depositAmount.String()),
		),
	}}
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
	if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Owner, withdrawAmount.ToCoins()); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			withdrawAmount.String())).Result()
	}

	// delete pool
	k.DeleteFarmPool(ctx, msg.PoolName)

	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(types.AttributeKeyWithdraw, withdrawAmount.String()),
		),
	}}
}

func handleMsgProvide(ctx sdk.Context, k keeper.Keeper, msg types.MsgProvide, logger log.Logger) sdk.Result {
	// 0.1 Check if the start_height_to_yield is more than current height
	if msg.StartHeightToYield <= ctx.BlockHeight() {
		return types.ErrInvalidStartHeight(DefaultCodespace).Result()
	}

	// 0.2 Get the pool info
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}

	// 0.3 Check if the provided coin denom is the same as the locked coin name
	if len(pool.YieldedTokenInfos) != 1 { // TODO: use the panic temporarily
		panic(fmt.Sprintf("The YieldedTokenInfos length is %d, which should be 1 in current code version",
			len(pool.YieldedTokenInfos)))
	}
	if strings.Compare(pool.YieldedTokenInfos[0].RemainingAmount.Denom, msg.Amount.Denom) != 0 {
		return types.ErrInvalidDenom(
			DefaultCodespace, pool.YieldedTokenInfos[0].RemainingAmount.Denom, msg.Amount.Denom).Result()
	}

	// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
	updatedPool := liquidateYieldTokenInfo(ctx.BlockHeight(), pool)
	updatedPool.LastClaimedBlockHeight = ctx.BlockHeight()
	// Check if remaining amount is zero already
	if updatedPool.YieldedTokenInfos[0].RemainingAmount.IsZero() {
		// 2. refresh the yielding_coin if remaining amount is zero
		updatedPool.YieldedTokenInfos[0] = types.NewYieldedTokenInfo(msg.Amount, msg.StartHeightToYield, msg.AmountYieldedPerBlock)

		// 3. Transfer coin to farm module account
		if err := k.SupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Address, ModuleName, msg.Amount.ToCoins()); err != nil {
			return err.Result()
		}
	}
	k.SetFarmPool(ctx, updatedPool)

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeProvide,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyStartHeightToYield, strconv.FormatInt(msg.StartHeightToYield, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountYieldPerBlock, msg.AmountYieldedPerBlock.String()),
		),
	}}
}

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock, logger log.Logger) sdk.Result {
	// 0. Get the pool info
	pool, poolFound := k.GetFarmPool(ctx, msg.PoolName)
	if !poolFound {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}
	if strings.Compare(pool.SymbolLocked, msg.Amount.Denom) != 0 {
		return types.ErrInvalidDenom(DefaultCodespace, pool.SymbolLocked, msg.Amount.Denom).Result()
	}

	if lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); !found {
		// 1. If lock info doesn't exist, only initialize the LockInfo structure
		lockInfo = types.NewLockInfo(msg.Address, msg.PoolName, msg.Amount, ctx.BlockHeight())
		k.SetLockInfo(ctx, lockInfo)
	} else {
		// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
		updatedPool := liquidateYieldTokenInfo(ctx.BlockHeight(), pool)

		// 2. Claim
		err := claim(ctx, k, updatedPool, msg.Address, msg.Amount.Amount)
		if err != nil {
			return err.Result()
		}
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
	if lockInfo, found := k.GetLockInfo(ctx, msg.Address, msg.PoolName); !found {
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
	err := claim(ctx, k, updatedPool, msg.Address, sdk.ZeroDec())
	if err != nil {
		return err.Result()
	}

	// 3. Send the locked-tokens from farm module account to its own account
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

func handleMsgClaim(ctx sdk.Context, k keeper.Keeper, msg types.MsgClaim, logger log.Logger) sdk.Result {
	// 0 Get the pool info
	pool, found := k.GetFarmPool(ctx, msg.PoolName)
	if !found {
		return types.ErrNoFarmPoolFound(DefaultCodespace, msg.PoolName).Result()
	}

	// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
	updatedPool := liquidateYieldTokenInfo(ctx.BlockHeight(), pool)

	// 2. Claim
	err := claim(ctx, k, updatedPool, msg.Address, sdk.ZeroDec())
	if err != nil {
		return err.Result()
	}

	// Emit events
	return sdk.Result{Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPool, msg.PoolName),
		),
	}}
}

// liquidateYieldTokenInfo is used for calculating how many tokens haven been yielding from LastClaimedBlockHeight to CurrentHeight
// Then transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
func liquidateYieldTokenInfo(height int64, pool types.FarmPool) types.FarmPool {
	if height <= pool.LastClaimedBlockHeight { // TODO: is there any neccessary to make a height comparison?
		return pool
	}

	// TODO: there are too many operations about MulTruncate, check the amount carefully, and write checking codes in invariants.go !!!
	// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		if height > startBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if startBlockHeightToYield > pool.LastClaimedBlockHeight {
				blockInterval = sdk.NewDec(height - startBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(height - pool.LastClaimedBlockHeight)
			}

			// calculate how many coin have been yielded till the current block
			amountYielded := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
			remainingAmount := pool.YieldedTokenInfos[i].RemainingAmount
			if amountYielded.LT(remainingAmount.Amount) {
				// add yielded amount
				pool.AmountYielded = pool.AmountYielded.Add(sdk.NewDecCoinsFromDec(remainingAmount.Denom, amountYielded))
				// subtract yielded_coin amount
				pool.YieldedTokenInfos[i].RemainingAmount.Amount = remainingAmount.Amount.Sub(amountYielded)
			} else {
				// add yielded amount
				pool.AmountYielded = pool.AmountYielded.Add(sdk.NewCoins(remainingAmount))

				// initialize yieldedTokenInfo
				pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remainingAmount.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())

				// TODO: remove the YieldedTokenInfo when its amount become zero
				// Currently, we support only one token of yield farming at the same time,
				// so, it is unnecessary to remove the element in slice
			}
		}
	}

	return pool
}

func claim(ctx sdk.Context, k keeper.Keeper, pool types.FarmPool, address sdk.AccAddress, changedAmount sdk.Dec) sdk.Error {
	// 0. Get lock_info
	lockInfo, found := k.GetLockInfo(ctx, address, pool.Name)
	if !found {
		return types.ErrNoLockInfoFound(DefaultCodespace, address.String())
	}

	height := ctx.BlockHeight()
	currentHeight := sdk.NewDec(height)
	/* 1.1 Calculate its own weight during these blocks
	   (curHeight - Height1) * Amount1
	*/
	// TODO: is there any possibility that lockInfo.StartBlockHeight is more than ctx.BlockHeight()?
	oldWeight := sdk.NewDec(lockInfo.StartBlockHeight).MulTruncate(lockInfo.Amount.Amount)
	currentWeight := currentHeight.MulTruncate(lockInfo.Amount.Amount)
	numerator := currentWeight.Sub(oldWeight)

	/* 1.2 Calculate all weight during these blocks
	    (curHeight - Height1) * Amount1 + (curHeight - Height2) * Amount2 + (curHeight - Height3) * Amount3
												||
	                                            \/
	   curHeight * (Amount1 + Amount2 + Amount3) - (Height1*Amount1 + Height2*Amount2 + Height3*Amount3)
												||
	                                            \/
	ctx.BlockHeight()  *  pool.TotalValueLocked.Amount  -  ( pool.TotalLockedWeight )
	*/
	denominator := currentHeight.MulTruncate(pool.TotalValueLocked.Amount).Sub(pool.TotalLockedWeight)

	// 1.3 Calculate how many yielded token it could get
	selfAmountYielded := pool.AmountYielded.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	// 2. Transfer yielded tokens to personal account
	if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, address, selfAmountYielded); err != nil {
		return err
	}

	// 3. Update the pool data
	pool.AmountYielded = pool.AmountYielded.Sub(selfAmountYielded)
	if !changedAmount.IsZero() {
		pool.TotalValueLocked.Amount = pool.TotalValueLocked.Amount.Add(changedAmount)
		currentWeight = currentWeight.Add(currentHeight.MulTruncate(changedAmount))
	}
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(currentWeight)
	pool.LastClaimedBlockHeight = height
	// Set the updated pool into store
	k.SetFarmPool(ctx, pool)

	// 4. Update the lock_info data
	if !changedAmount.IsZero() {
		lockInfo.Amount.Amount = lockInfo.Amount.Amount.Add(changedAmount)
		if lockInfo.Amount.IsZero() { // If amount become zero, delete the lock_info
			k.DeleteLockInfo(ctx, address, pool.Name)
		} else { // Otherwise, update the lock_info
			lockInfo.StartBlockHeight = height
			k.SetLockInfo(ctx, lockInfo)
		}
	}

	return nil
}
