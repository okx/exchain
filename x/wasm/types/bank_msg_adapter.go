package types

import (
	"context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

func (msg MsgSend) Route() string {
	return RouterKey
}

func (msg MsgSend) Type() string {
	return "bank"
}

func (msg MsgSend) ValidateBasic() error {

	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return sdkerrors.Wrap(err, "sender")
	}

	if _, err := sdk.AccAddressFromBech32(msg.ToAddress); err != nil {
		return sdkerrors.Wrap(err, "to")
	}
	return msg.Amount.Validate()
}

func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}

}

type BankMsgServer struct {
	bankKeeper BankKeeper
}

func NewBankMsgServer(bankKeeper BankKeeper) *BankMsgServer {
	return &BankMsgServer{bankKeeper: bankKeeper}
}

func (bms BankMsgServer) Send(goCtx context.Context, msg *MsgSend) (*MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "to")
	}
	coins := CoinAdaptersToCoins(msg.Amount)
	if err := bms.bankKeeper.SendCoins(ctx, senderAddr, toAddr, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "send coins")
	}
	return &MsgSendResponse{}, nil
}
