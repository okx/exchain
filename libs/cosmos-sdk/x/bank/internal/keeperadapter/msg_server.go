package keeperadapter

import (
	"context"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/typesadapter"
)

type msgServer struct {
	MsgServerBankKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper MsgServerBankKeeper) typesadapter.MsgServer {
	return &msgServer{MsgServerBankKeeper: keeper}
}

var _ typesadapter.MsgServer = msgServer{}

func (k msgServer) Send(goCtx context.Context, msg *typesadapter.MsgSend) (*typesadapter.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.GetSendEnabled(ctx) {
		return nil, types.ErrSendDisabled
	}

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	if k.BlacklistedAddr(to) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}
	amt := sdk.CoinAdaptersToCoins(msg.Amount)
	err = k.SendCoins(ctx, from, to, amt)
	if err != nil {
		return nil, err
	}

	//For using 0x prefix address at wasm tx ,so we must be follow code
	ctx.EventManager().PopEvent()
	ctx.EventManager().PopEvent()
	ctx.EventManager().EmitEvents(sdk.Events{
		// This event should have all info (to, from, amount) without looking at other events
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, sdk.AccToAWasmddress(to).String()),
			sdk.NewAttribute(types.AttributeKeySender, sdk.AccToAWasmddress(from).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, sdk.AccToAWasmddress(from).String()),
		),
	})

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &typesadapter.MsgSendResponse{}, nil
}

func (k msgServer) MultiSend(goCtx context.Context, msg *typesadapter.MsgMultiSend) (*typesadapter.MsgMultiSendResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//
	//// NOTE: totalIn == totalOut should already have been checked
	//for _, in := range msg.Inputs {
	//	k.GetSendEnabled(ctx)
	//	if err := k.IsSendEnabledCoins(ctx, in.Coins...); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//for _, out := range msg.Outputs {
	//	accAddr, err := sdk.AccAddressFromBech32(out.Address)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if k.BlacklistedAddr(accAddr) {
	//		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", out.Address)
	//	}
	//}
	//
	//err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	//if err != nil {
	//	return nil, err
	//}
	//
	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		sdk.EventTypeMessage,
	//		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	//	),
	//)

	return nil, sdkerrors.Wrap(types.ErrSendDisabled, "MultiSend is not support")
}
