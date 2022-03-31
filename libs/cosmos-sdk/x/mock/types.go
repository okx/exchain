package mock

import (
	"github.com/okex/exchain/libs/tendermint/crypto"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
)

// DummySupplyKeeper defines a supply keeper used only for testing to avoid
// circle dependencies
type DummySupplyKeeper struct {
	ak auth.AccountKeeper
}

func (sk DummySupplyKeeper) AddCoinsToFeeCollector(ctx sdk.Context, amt sdk.Coins) error {
	return nil
}

func (sk DummySupplyKeeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (sk DummySupplyKeeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (sk DummySupplyKeeper) AddFee(amt sdk.Coins) {}

func (sk DummySupplyKeeper) GetFee() sdk.Coins {
	return nil
}

func (sk DummySupplyKeeper) ResetFee() {}

// NewDummySupplyKeeper creates a DummySupplyKeeper instance
func NewDummySupplyKeeper(ak auth.AccountKeeper) DummySupplyKeeper {
	return DummySupplyKeeper{ak}
}

// SendCoinsFromAccountToModule for the dummy supply keeper
func (sk DummySupplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	fromAcc := sk.ak.GetAccount(ctx, fromAddr)
	moduleAcc := sk.GetModuleAccount(ctx, recipientModule)

	newFromCoins, hasNeg := fromAcc.GetCoins().SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, fromAcc.GetCoins().String())
	}

	newToCoins := moduleAcc.GetCoins().Add(amt...)

	if err := fromAcc.SetCoins(newFromCoins); err != nil {
		return err
	}

	if err := moduleAcc.SetCoins(newToCoins); err != nil {
		return err
	}

	sk.ak.SetAccount(ctx, fromAcc)
	sk.ak.SetAccount(ctx, moduleAcc)

	return nil
}

// GetModuleAccount for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI {
	addr := sk.GetModuleAddress(moduleName)

	acc := sk.ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(exported.ModuleAccountI)
		if ok {
			return macc
		}
	}

	moduleAddress := sk.GetModuleAddress(moduleName)
	baseAcc := auth.NewBaseAccountWithAddress(moduleAddress)

	// create a new module account
	macc := &supply.ModuleAccount{
		BaseAccount: &baseAcc,
		Name:        moduleName,
		Permissions: nil,
	}

	maccI := (sk.ak.NewAccount(ctx, macc)).(exported.ModuleAccountI)
	sk.ak.SetAccount(ctx, maccI)
	return maccI
}

// GetModuleAddress for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(moduleName)))
}

func (sk DummySupplyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	recipientAcc := sk.ak.GetAccount(ctx, recipientAddr)
	moduleAcc := sk.GetModuleAccount(ctx, senderModule)

	newFromCoins, hasNeg := moduleAcc.GetCoins().SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, moduleAcc.GetCoins().String())
	}

	newToCoins := recipientAcc.GetCoins().Add(amt...)

	if err := moduleAcc.SetCoins(newFromCoins); err != nil {
		return err
	}

	if err := recipientAcc.SetCoins(newToCoins); err != nil {
		return err
	}

	sk.ak.SetAccount(ctx, recipientAcc)
	sk.ak.SetAccount(ctx, moduleAcc)

	return nil
}
