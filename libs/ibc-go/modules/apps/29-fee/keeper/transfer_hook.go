package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	transfertypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	types2 "github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ transfertypes.TransferHooks = (*FeeTransferHook)(nil)
)

type FeeTransferHook struct {
	k *Keeper
}

func NewFeeTransferHook(k *Keeper) *FeeTransferHook {
	return &FeeTransferHook{k: k}
}

func (f *FeeTransferHook) AfterSendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.SysCoin, sender sdk.AccAddress, receiver string, isSource bool, p types.Packet) error {
	if !types2.HigherThanVenus4(ctx.BlockHeight()) {
		return nil
	}
	f.k.AddPacket(types.NewSignerPacketWrapper(p, []sdk.AccAddress{sender}, ctx.GasMeter().GasConsumed()))
	return nil
}

func (f *FeeTransferHook) AfterRecvTransfer(ctx sdk.Context, destPort, destChannel string, token sdk.SysCoin, receiver string, isSource bool) error {
	return nil
}

func (f *FeeTransferHook) AfterRefundTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.SysCoin, sender string, isSource bool) error {
	return nil
}
