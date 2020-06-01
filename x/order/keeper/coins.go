package keeper

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/types"
)

// LockCoins locks coins from the specified address,
func (k Keeper) LockCoins(ctx sdk.Context, order *types.Order, coins sdk.DecCoins, lockCoinsType int) error {
	if coins.IsZero() {
		return nil
	}
	switch order.Type {
	case types.OrdinaryOrder:
		if err := k.tokenKeeper.LockCoins(ctx, order.Sender, coins, lockCoinsType); err != nil {
			return err
		}
	case types.MarginOrder:
	default:
		return fmt.Errorf("unrecognized order type:%d", order.Type)
	}
	return nil
}

// nolint
func (k Keeper) UnlockCoins(ctx sdk.Context, order *types.Order, coins sdk.DecCoins, lockCoinsType int) {
	if coins.IsZero() {
		return
	}
	switch order.Type {
	case types.OrdinaryOrder:
		if err := k.tokenKeeper.UnlockCoins(ctx, order.Sender, coins, lockCoinsType); err != nil {
			log.Printf("User(%s) unlock coins(%s) failed\n", order.Sender.String(), coins.String())
		}
	case types.MarginOrder:
	}
}

// AddCollectedFees adds fee to the feePool
func (k Keeper) AddCollectedFees(ctx sdk.Context, order *types.Order, coins sdk.DecCoins,
	feeType string, hasFeeDetail bool) error {
	if coins.IsZero() {
		return nil
	}
	if hasFeeDetail {
		k.tokenKeeper.AddFeeDetail(ctx, order.Sender.String(), coins, feeType)
	}

	baseCoins := coins
	switch order.Type {
	case types.OrdinaryOrder:
		return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, order.Sender, k.feeCollectorName, baseCoins)
	default:
		return fmt.Errorf("unrecognized order type:%d", order.Type)
	}
	return nil
}

// SendFeesToProductOwner sends fees from the specified address to productOwner
func (k Keeper) SendFeesToProductOwner(ctx sdk.Context, order *types.Order, coins sdk.DecCoins,
	feeType string, product string) error {
	if coins.IsZero() {
		return nil
	}
	to := k.GetProductOwner(ctx, product)

	switch order.Type {
	case types.OrdinaryOrder:
		k.tokenKeeper.AddFeeDetail(ctx, order.Sender.String(), coins, feeType)
		if err := k.tokenKeeper.SendCoinsFromAccountToAccount(ctx, order.Sender, to, coins); err != nil {
			log.Printf("Send fee(%s) to address(%s) failed\n", coins.String(), to.String())
			return err
		}
	case types.MarginOrder:
	default:
		return fmt.Errorf("unrecognized order type:%d", order.Type)
	}
	return nil
}

// BalanceAccount burns the specified coin and obtains another coin
func (k Keeper) BalanceAccount(ctx sdk.Context, order *types.Order,
	outputCoins sdk.DecCoins, inputCoins sdk.DecCoins) {
	switch order.Type {
	case types.OrdinaryOrder:
		if err := k.tokenKeeper.BalanceAccount(ctx, order.Sender, outputCoins, inputCoins); err != nil {
			log.Printf("User(%s) burn locked coins(%s) failed\n", order.Sender.String(), outputCoins.String())
		}
	case types.MarginOrder:
	}

}
