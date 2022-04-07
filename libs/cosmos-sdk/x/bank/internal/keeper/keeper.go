package keeper

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	vestexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/vesting/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper

	DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error

	GetInnerTxKeeper() innertx.InnerTxKeeper
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper

	ak         types.AccountKeeper
	paramSpace params.Subspace

	marshal *codec.CodecProxy
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(
	ak types.AccountKeeper, paramSpace params.Subspace, blacklistedAddrs map[string]bool,
) BaseKeeper {

	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return BaseKeeper{
		BaseSendKeeper: NewBaseSendKeeper(ak, ps, blacklistedAddrs),
		ak:             ak,
		paramSpace:     ps,
	}
}

func NewBaseKeeperWithMarshal(ak types.AccountKeeper, marshal *codec.CodecProxy, paramSpace params.Subspace, blacklistedAddrs map[string]bool,
) BaseKeeper {
	ret := NewBaseKeeper(ak, paramSpace, blacklistedAddrs)
	ret.marshal = marshal
	return ret
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins.
// The coins are then transferred from the delegator address to a ModuleAccount address.
// If any of the delegation amounts are negative, an error is returned.
func (keeper BaseKeeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) (err error) {
	defer func() {
		if !ctx.IsCheckTx() && keeper.ik != nil {
			keeper.ik.UpdateInnerTx(ctx.TxBytes(), innertx.CosmosDepth, delegatorAddr, moduleAccAddr, innertx.CosmosCallType, innertx.DelegateCallName, amt, err)
		}
	}()
	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", delegatorAddr)
	}

	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleAccAddr)
	}

	if !amt.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	oldCoins := delegatorAcc.GetCoins()

	_, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", oldCoins, amt,
		)
	}

	if err := trackDelegation(delegatorAcc, ctx.BlockHeader().Time, amt); err != nil {
		return sdkerrors.Wrap(err, "failed to track delegation")
	}

	keeper.ak.SetAccount(ctx, delegatorAcc, false)

	_, err = keeper.AddCoins(ctx, moduleAccAddr, amt)
	if err != nil {
		return err
	}

	return nil
}

// UndelegateCoins performs undelegation by crediting amt coins to an account with
// address addr. For vesting accounts, undelegation amounts are tracked for both
// vesting and vested coins.
// The coins are then transferred from a ModuleAccount address to the delegator address.
// If any of the undelegation amounts are negative, an error is returned.
func (keeper BaseKeeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) (err error) {
	defer func() {
		if !ctx.IsCheckTx() && keeper.ik != nil {
			keeper.ik.UpdateInnerTx(ctx.TxBytes(), innertx.CosmosDepth, moduleAccAddr, delegatorAddr, innertx.CosmosCallType, innertx.UndelegateCallName, amt, err)
		}
	}()

	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", delegatorAddr)
	}

	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleAccAddr)
	}

	if !amt.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	oldCoins := moduleAcc.GetCoins()

	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", oldCoins, amt,
		)
	}

	if err = keeper.SetCoins(ctx, moduleAccAddr, newCoins); err != nil {
		return err
	}

	if err := trackUndelegation(delegatorAcc, amt); err != nil {
		return sdkerrors.Wrap(err, "failed to track undelegation")
	}

	keeper.ak.SetAccount(ctx, delegatorAcc, false)
	return nil
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error

	GetSendEnabled(ctx sdk.Context) bool
	SetSendEnabled(ctx sdk.Context, enabled bool)

	BlacklistedAddr(addr sdk.AccAddress) bool
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak         types.AccountKeeper
	ask        authexported.SizerAccountKeeper
	paramSpace params.Subspace

	// list of addresses that are restricted from receiving transactions
	blacklistedAddrs map[string]bool

	ik innertx.InnerTxKeeper
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(
	ak types.AccountKeeper, paramSpace params.Subspace, blacklistedAddrs map[string]bool,
) BaseSendKeeper {

	bsk := BaseSendKeeper{
		BaseViewKeeper:   NewBaseViewKeeper(ak),
		ak:               ak,
		paramSpace:       paramSpace,
		blacklistedAddrs: blacklistedAddrs,
	}
	bsk.ask, _ = bsk.ak.(authexported.SizerAccountKeeper)
	return bsk
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseSendKeeper) InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) (err error) {
	defer func() {
		if !ctx.IsCheckTx() && keeper.ik != nil {
			for _, in := range inputs {
				keeper.ik.UpdateInnerTx(ctx.TxBytes(), innertx.CosmosDepth, in.Address, sdk.AccAddress{}, innertx.CosmosCallType, innertx.MultiCallName, in.Coins, err)
			}

			for _, out := range outputs {
				keeper.ik.UpdateInnerTx(ctx.TxBytes(), innertx.CosmosDepth, sdk.AccAddress{}, out.Address, innertx.CosmosCallType, innertx.MultiCallName, out.Coins, err)
			}
		}
	}()
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	for _, in := range inputs {
		_, err := keeper.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}

	for _, out := range outputs {
		_, err := keeper.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeyRecipient, out.Address.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, out.Coins.String()),
			),
		)

		// Create account if recipient does not exist.
		//
		// NOTE: This should ultimately be removed in favor a more flexible approach
		// such as delegated fee messages.
		acc := keeper.ak.GetAccount(ctx, out.Address)
		if acc == nil {
			keeper.ak.SetAccount(ctx, keeper.ak.NewAccountWithAddress(ctx, out.Address), false)
		}
	}

	return nil
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (err error) {
	defer func() {
		if !ctx.IsCheckTx() && keeper.ik != nil {
			keeper.ik.UpdateInnerTx(ctx.TxBytes(), innertx.CosmosDepth, fromAddr, toAddr, innertx.CosmosCallType, innertx.SendCallName, amt, err)
		}
	}()
	fromAddrStr := fromAddr.String()
	ctx.EventManager().EmitEvents(sdk.Events{
		// This event should have all info (to, from, amount) without looking at other events
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(types.AttributeKeySender, fromAddrStr),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddrStr),
		),
	})

	if !amt.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	fromAcc, _ := ctx.GetFromAccountCacheData().(authexported.Account)
	toAcc, _ := ctx.GetToAccountCacheData().(authexported.Account)
	fromAccGas, toAccGas := ctx.GetFromAccountCacheGas(), ctx.GetToAccountCacheGas()

	fromAcc, fromAccGas = keeper.getAccount(&ctx, fromAddr, fromAcc, fromAccGas)
	_, err = keeper.subtractCoins(ctx, fromAddr, fromAcc, fromAccGas, amt)
	if err != nil {
		return err
	}

	ctx.UpdateFromAccountCache(fromAcc, 0)

	toAcc, toAccGas = keeper.getAccount(&ctx, toAddr, toAcc, toAccGas)
	_, err = keeper.addCoins(ctx, toAddr, toAcc, toAccGas, amt)
	if err != nil {
		return err
	}

	ctx.UpdateToAccountCache(toAcc, 0)

	return nil
}

// SubtractCoins subtracts amt from the coins at the addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func (keeper BaseSendKeeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	if !amt.IsValid() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}
	acc, gasUsed := authexported.GetAccountAndGas(&ctx, keeper.ak, addr)
	return keeper.subtractCoins(ctx, addr, acc, gasUsed, amt)
}

func (keeper *BaseSendKeeper) subtractCoins(ctx sdk.Context, addr sdk.AccAddress, acc authexported.Account, accGas sdk.Gas, amt sdk.Coins) (sdk.Coins, error) {
	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockTime())
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", spendableCoins, amt,
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := keeper.setCoinsToAccount(ctx, addr, acc, accGas, newCoins)

	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseSendKeeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	if !amt.IsValid() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	// oldCoins := keeper.GetCoins(ctx, addr)

	acc, gasUsed := authexported.GetAccountAndGas(&ctx, keeper.ak, addr)
	return keeper.addCoins(ctx, addr, acc, gasUsed, amt)
}

func (keeper *BaseSendKeeper) addCoins(ctx sdk.Context, addr sdk.AccAddress, acc authexported.Account, accGas sdk.Gas, amt sdk.Coins) (sdk.Coins, error) {
	var oldCoins sdk.Coins
	if acc == nil {
		oldCoins = sdk.NewCoins()
	} else {
		oldCoins = acc.GetCoins()
	}

	newCoins := oldCoins.Add(amt...)

	if newCoins.IsAnyNegative() {
		return amt, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", oldCoins, amt,
		)
	}

	err := keeper.setCoinsToAccount(ctx, addr, acc, accGas, newCoins)

	return newCoins, err
}

// SetCoins sets the coins at the addr.
func (keeper BaseSendKeeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	if !amt.IsValid() {
		sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		panic(err)
	}
	keeper.ak.SetAccount(ctx, acc, false)
	return nil
}

func (keeper *BaseSendKeeper) getAccount(ctx *sdk.Context, addr sdk.AccAddress, acc authexported.Account, getgas sdk.Gas) (authexported.Account, sdk.Gas) {
	gasMeter := ctx.GasMeter()
	if acc != nil && bytes.Equal(acc.GetAddress(), addr) {
		if getgas > 0 {
			gasMeter.ConsumeGas(getgas, "get account")
			return acc, getgas
		}
		if ok, gasused := authexported.TryAddGetAccountGas(gasMeter, keeper.ask, acc); ok {
			return acc, gasused
		}
	}
	return authexported.GetAccountAndGas(ctx, keeper.ak, addr)
}

func (keeper *BaseSendKeeper) setCoinsToAccount(ctx sdk.Context, addr sdk.AccAddress, acc authexported.Account, accGas sdk.Gas, amt sdk.Coins) error {
	if !amt.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amt.String())
	}

	acc, _ = keeper.getAccount(&ctx, addr, acc, accGas)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		panic(err)
	}

	keeper.ak.SetAccount(ctx, acc, false)
	return nil
}

// GetSendEnabled returns the current SendEnabled
func (keeper BaseSendKeeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

// BlacklistedAddr checks if a given address is blacklisted (i.e restricted from
// receiving funds)
func (keeper BaseSendKeeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	return keeper.blacklistedAddrs[addr.String()]
}

// SetInnerTxKeeper set innerTxKeeper
func (k *BaseKeeper) SetInnerTxKeeper(keeper innertx.InnerTxKeeper) {
	k.BaseSendKeeper.SetInnerTxKeeper(keeper)
}

func (k *BaseSendKeeper) SetInnerTxKeeper(keeper innertx.InnerTxKeeper) {
	k.ik = keeper
}

func (k BaseSendKeeper) GetInnerTxKeeper() innertx.InnerTxKeeper {
	return k.ik
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak types.AccountKeeper
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak types.AccountKeeper) BaseViewKeeper {
	return BaseViewKeeper{ak: ak}
}

// Logger returns a module-specific logger.
func (keeper BaseViewKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

// CONTRACT: assumes that amt is valid.
func trackDelegation(acc authexported.Account, blockTime time.Time, amt sdk.Coins) error {
	vacc, ok := acc.(vestexported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackDelegation
		vacc.TrackDelegation(blockTime, amt)
	}

	return acc.SetCoins(acc.GetCoins().Sub(amt))
}

// CONTRACT: assumes that amt is valid.
func trackUndelegation(acc authexported.Account, amt sdk.Coins) error {
	vacc, ok := acc.(vestexported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackUndelegation
		vacc.TrackUndelegation(amt)
	}

	return acc.SetCoins(acc.GetCoins().Add(amt...))
}
