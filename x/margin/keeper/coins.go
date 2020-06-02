package keeper

import (
	"fmt"

	"github.com/okex/okchain/x/margin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) LockCoins(ctx sdk.Context, address sdk.AccAddress, product string, coins sdk.DecCoins) error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return fmt.Errorf(fmt.Sprintf("margin account %s does not exist", address.String()))
	}
	_, hasNeg := account.Available.SafeSub(coins)
	if hasNeg {
		return fmt.Errorf("insufficient account funds; %s is less than %s", account.Available, coins)
	}
	account.Available = account.Available.Sub(coins)
	account.Locked = account.Locked.Add(coins)
	k.SetAccount(ctx, address, product, account)
	return nil
}

func (k Keeper) UnLockCoins(ctx sdk.Context, address sdk.AccAddress, product string, coins sdk.DecCoins) error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return fmt.Errorf(fmt.Sprintf("margin account %s does not exist", address.String()))
	}
	_, hasNeg := account.Locked.SafeSub(coins)
	if hasNeg {
		return fmt.Errorf("failed to unlock <%s>. Address <%s>, coins locked <%s>", coins, address, account.Locked)
	}
	account.Available = account.Available.Add(coins)
	account.Locked = account.Locked.Sub(coins)
	k.SetAccount(ctx, address, product, account)
	return nil
}

func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, product string, recipientModule string, coins sdk.DecCoins) error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return fmt.Errorf(fmt.Sprintf("margin account %s does not exist", address.String()))
	}
	_, hasNeg := account.Available.SafeSub(account.Interest.Add(account.Borrowed).Add(coins))
	if hasNeg {
		return fmt.Errorf("insufficient account funds; account available<%s> borrowed<%s> interest<%s> is less than %s", account.Available, account.Borrowed, account.Interest, coins)
	}
	account.Available = account.Available.Sub(coins)
	k.SetAccount(ctx, address, product, account)
	if err := k.tokenKeeper.SendCoinsFromModuleToModule(ctx, recipientModule, coins); err != nil {
		return fmt.Errorf("insufficient module account %s, needs %s", types.ModuleAccount, coins)
	}
	return nil
}

func (k Keeper) SendCoinsFromAccountToAccount(ctx sdk.Context, address sdk.AccAddress, product string, to sdk.AccAddress, coins sdk.DecCoins) error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return fmt.Errorf(fmt.Sprintf("margin account %s does not exist", address.String()))
	}
	_, hasNeg := account.Available.SafeSub(account.Interest.Add(account.Borrowed).Add(coins))
	if hasNeg {
		return fmt.Errorf("insufficient account funds; account available<%s> borrowed<%s> interest<%s> is less than %s", account.Available, account.Borrowed, account.Interest, coins)
	}
	account.Available = account.Available.Sub(coins)
	k.SetAccount(ctx, address, product, account)
	if err := k.tokenKeeper.SendCoinsFromModuleToAccount(ctx, to, coins); err != nil {
		return fmt.Errorf("insufficient module account %s, needs %s", types.ModuleAccount, coins)
	}
	return nil
}

func (k Keeper) BalanceAccount(ctx sdk.Context, address sdk.AccAddress, product string, outputCoins, inputCoins sdk.DecCoins) error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return fmt.Errorf(fmt.Sprintf("margin account %s does not exist", address.String()))
	}
	_, hasNeg := account.Locked.SafeSub(outputCoins)
	if hasNeg {
		return fmt.Errorf("failed to unlock <%s>. Address <%s>, coins locked <%s>", outputCoins, address, account.Locked)
	}
	account.Available = account.Available.Add(inputCoins)
	account.Locked = account.Locked.Sub(outputCoins)
	k.SetAccount(ctx, address, product, account)
	return nil
}
