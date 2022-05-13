package types

import (
	"context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	bank "github.com/okex/exchain/libs/cosmos-sdk/x/bank"
)

type BankMsgServer struct {
	bankKeeper BankKeeper
}

func NewBankMsgServer(bankKeeper BankKeeper) *BankMsgServer {
	return &BankMsgServer{bankKeeper: bankKeeper}
}

func (bms BankMsgServer) Send(goCtx context.Context, msg *bank.MsgSendAdapter) (*bank.MsgSendResponseAdapter, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	toAddr, err := AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "to")
	}
	coins := sdk.CoinAdaptersToCoins(msg.Amount)
	if err := bms.bankKeeper.SendCoins(ctx, senderAddr, toAddr, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "send coins")
	}
	return &bank.MsgSendResponseAdapter{}, nil
}

// MultiSend defines a method for sending coins from some accounts to other accounts.
func (bms BankMsgServer) MultiSend(context.Context, *bank.MsgMultiSendAdapter) (*bank.MsgMultiSendResponseAdapter, error) {
	return nil, sdkerrors.Wrap(ErrInvalid, "MultiSend is not support")
}
