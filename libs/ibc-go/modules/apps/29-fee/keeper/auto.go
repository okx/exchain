package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// SetParams sets the total set of ibc-transfer parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) EscrowPacketFeeFromFeeCollector(ctx sdk.Context, fee sdk.Coins) error {
	if len(k.packets) == 0 {
		return nil
	}

	if !k.AllowAutoDispatch(ctx) {
		return nil
	}

	avail := make(map[string]exported.PacketI)
	for key, p := range k.packets {
		// NOTE: we only consume the recv-packet
		if !k.IsFeeEnabled(ctx, p.GetSourcePort(), p.GetSourceChannel()) {
			continue
		}
		avail[key] = p
	}

	percent := k.GetFeePercent(ctx)
	feePercent := sdk.NewDecWithPrec(int64(percent), 2)
	fee = fee.MulDecTruncate(feePercent)

	bkk := k.bk
	skk := k.supplyK
	_, err := bkk.AddCoins(ctx, skk.GetModuleAccount(ctx, types.ModuleName).GetAddress(), fee)
	if nil != err {
		return err
	}
	recPortion := sdk.NewDecWithPrec(int64(k.GetRecvFeePercent(ctx)), 2)
	ackPortion := sdk.NewDecWithPrec(int64(k.GetAckFeePercent(ctx)), 2)
	timeoutPortion := sdk.NewDecWithPrec(int64(k.GetTimeOutFeePercent(ctx)), 2)
	packetPercent := (1 / float64(len(avail))) * 100
	packetPerFee := fee.MulDecTruncate(sdk.NewDecWithPrec(int64(packetPercent), 2))
	for key, v := range avail {
		rev := packetPerFee.MulDecTruncate(recPortion)
		ack := packetPerFee.MulDecTruncate(ackPortion)
		tf := packetPerFee.MulDecTruncate(timeoutPortion)
		k.Logger(ctx).Info("register packet fee", "recvFee", rev.String(), "ack", ack.String(), "timeout", tf.String())
		packetF := types.Fee{
			RecvFee:    sdk.CoinsToCoinAdapters(rev),
			AckFee:     sdk.CoinsToCoinAdapters(ack),
			TimeoutFee: sdk.CoinsToCoinAdapters(tf),
		}
		k.Logger(ctx).Info("register packet fee", "packetKey", key, fee, packetF.String())
		if err = k.escrowPacketFee(ctx, channeltypes.PacketId{
			PortId:    v.GetSourcePort(),
			ChannelId: v.GetSourceChannel(),
			Sequence:  v.GetSequence(),
		}, types.PacketFee{
			Fee:           packetF,
			RefundAddress: skk.GetModuleAccount(ctx, types.ModuleName).GetAddress().String(),
		}); nil != err {
			return err
		}
	}

	return nil
}
