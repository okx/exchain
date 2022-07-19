package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	trensferTypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/x/erc20/types"
)

var (
	_ trensferTypes.TransferHooks = IBCTransferHooks{}
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
	isSource bool) error {
	iths.Logger(ctx).Info(
		"trigger ibc transfer hook",
		"hook", "AfterSendTransfer",
		"sourcePort", sourcePort,
		"sourceChannel", sourceChannel,
		"token", token.String(),
		"sender", sender.String(),
		"receiver", receiver,
		"isSource", isSource)
	return nil
}

func (iths IBCTransferHooks) AfterRecvTransfer(
	ctx sdk.Context,
	destPort, destChannel string,
	token sdk.SysCoin,
	receiver string,
	isSource bool) error {
	iths.Logger(ctx).Info(
		"trigger ibc transfer hook",
		"hook", "AfterRecvTransfer",
		"destPort", destPort,
		"destChannel", destChannel,
		"token", token.String(),
		"receiver", receiver,
		"isSource", isSource)

	if !isSource {
		// only after minting vouchers on this chain
		// the native coin come from other chain with ibc
		if err := iths.Keeper.OnMintVouchers(ctx, sdk.NewCoins(token), receiver); err != types.ErrNoContractNotAuto {
			return err
		}
	} else if token.Denom != sdk.DefaultBondDenom {
		// the native coin come from this chain,
		return iths.Keeper.OnUnescrowNatives(ctx, sdk.NewCoins(token), receiver)
	}
	return nil
}

func (iths IBCTransferHooks) AfterRefundTransfer(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	token sdk.SysCoin,
	sender string,
	isSource bool) error {
	iths.Logger(ctx).Info(
		"trigger ibc transfer hook",
		"hook", "AfterRefundTransfer",
		"sourcePort", sourcePort,
		"sourceChannel", sourceChannel,
		"token", token.String(),
		"sender", sender,
		"isSource", isSource)
	// only after minting vouchers on this chain
	// the native coin come from other chain with ibc
	if !isSource {
		if err := iths.Keeper.OnMintVouchers(ctx, sdk.NewCoins(token), sender); err != types.ErrNoContractNotAuto {
			return err
		}
	} else if token.Denom != sdk.DefaultBondDenom {
		// the native coin come from this chain,
		return iths.Keeper.OnUnescrowNatives(ctx, sdk.NewCoins(token), sender)
	}
	return nil
}
