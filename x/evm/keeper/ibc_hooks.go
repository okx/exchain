package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
)

var (
	_ types.TransferHooks = IBCTransferHooks{}
)

type IBCTransferHooks struct {
	Keeper
}

func NewIBCTransferHooks(k Keeper) IBCTransferHooks {
	return IBCTransferHooks{k}
}

func (iths IBCTransferHooks) AfterSendTransfer(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	token sdk.SysCoin,
	sender sdk.AccAddress,
	receiver string,
	isSource bool) {
	return
}

func (iths IBCTransferHooks) AfterRecvTransfer(
	ctx sdk.Context,
	destPort, destChannel string,
	token sdk.SysCoin,
	receiver string,
	isSource bool) {
	// only after minting vouchers on this chain
	// the native coin come from other chain with ibc
	if !isSource {
		iths.Keeper.OnMintVouchers(ctx, sdk.NewCoins(token), receiver)
	}
}

func (iths IBCTransferHooks) AfterRefundTransfer(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	token sdk.SysCoin,
	sender string,
	isSource bool) {
	// only after minting vouchers on this chain
	// the native coin come from other chain with ibc
	if !isSource {
		iths.Keeper.OnMintVouchers(ctx, sdk.NewCoins(token), sender)
	}
}
