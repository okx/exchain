package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/saving/types"
)

// Keeper of the saving store
type Keeper struct {
	storeKey     sdk.StoreKey
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper
	cdc          *codec.Codec
	paramspace   types.ParamSubspace
}

// NewKeeper creates a saving keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, tokenKeeper types.TokenKeeper, supplyKeeper types.SupplyKeeper, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		storeKey:     key,
		supplyKeeper: supplyKeeper,
		tokenKeeper:  tokenKeeper,
		cdc:          cdc,
		paramspace:   paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Deposit deposits amount of tokens to saving module
func (k Keeper) Deposit(ctx sdk.Context, from sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	tokenPair := k.GetTokenPair(ctx, product)
	if tokenPair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because non-exist product: %s", product))
	}

	if !tokenPair.Owner.Equals(from) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to deposit because %s is not the owner of product:%s", from.String(), product))
	}

	if amount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because deposits only support %s token", sdk.DefaultBondDenom))
	}

	depositCoins := amount.ToCoins()
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, depositCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", depositCoins.String()))
	}

	tokenPair.Deposits = tokenPair.Deposits.Add(amount)
	k.UpdateTokenPair(ctx, product, tokenPair)
	return nil
}

// Withdraw withdraws amount of tokens from saving module
func (k Keeper) Withdraw(ctx sdk.Context, to sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	tokenPair := k.GetTokenPair(ctx, product)
	if tokenPair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to withdraws because non-exist product: %s", product))
	}

	if !tokenPair.Owner.Equals(to) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to withdraws because %s is not the owner of product:%s", to.String(), product))
	}

	if amount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to withdraws because deposits only support %s token", sdk.DefaultBondDenom))
	}

	if tokenPair.Deposits.IsLT(amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to withdraws because deposits:%s is less than withdraw:%s", tokenPair.Deposits.String(), amount.String()))
	}

	completeTime := ctx.BlockHeader().Time.Add(k.GetParams(ctx).WithdrawPeriod)
	// add withdraw info to store
	withdrawInfo, ok := k.GetWithdrawInfo(ctx, to)
	if !ok {
		withdrawInfo = types.WithdrawInfo{
			Owner:        to,
			Deposits:     amount,
			CompleteTime: completeTime,
		}
	} else {
		k.DeleteWithdrawCompleteTimeAddress(ctx, withdrawInfo.CompleteTime, to)
		withdrawInfo.Deposits = withdrawInfo.Deposits.Add(amount)
		withdrawInfo.CompleteTime = completeTime
	}
	k.SetWithdrawInfo(ctx, withdrawInfo)
	k.SetWithdrawCompleteTimeAddress(ctx, completeTime, to)

	// update token pair
	tokenPair.Deposits = tokenPair.Deposits.Sub(amount)
	k.UpdateTokenPair(ctx, product, tokenPair)
	return nil
}
