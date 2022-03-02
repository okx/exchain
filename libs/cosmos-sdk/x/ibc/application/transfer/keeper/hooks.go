package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/types"
)

// SetHooks sets the hooks for the IBC transfer module
// It should be called only once during initialization, it panics if called more than once.
func (k *Keeper) SetHooks(hooks types.TransferHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}

	k.hooks = hooks

	return k
}

func (k Keeper) CallAfterSendTransferHooks(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	token sdk.SysCoin,
	sender sdk.AccAddress,
	receiver string,
	isSource bool) {
	if k.hooks != nil {
		k.hooks.AfterSendTransfer(ctx, sourcePort, sourceChannel, token, sender, receiver, isSource)
	}
}
func (k Keeper) CallAfterRecvTransferHooks(
	ctx sdk.Context,
	destPort, destChannel string,
	token sdk.SysCoin,
	receiver string,
	isSource bool) {
	if k.hooks != nil {
		k.hooks.AfterRecvTransfer(ctx, destPort, destChannel, token, receiver, isSource)
	}
}
func (k Keeper) CallAfterRefundTransferHooks(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	token sdk.SysCoin,
	sender string,
	isSource bool) {
	if k.hooks != nil {
		k.hooks.AfterRefundTransfer(ctx, sourcePort, sourceChannel, token, sender, isSource)
	}
}
