package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (msg MsgSendToEvm) Route() string {
	return RouterKey
}

func (msg MsgSendToEvm) Type() string {
	return SendToEvmSubMsgName
}

func (msg MsgSendToEvm) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	if !sdk.IsWasmAddress(sender) {
		return ErrIsNotWasmAddr
	}

	contract, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}
	if sdk.IsWasmAddress(contract) {
		return ErrIsNotEvmAddr
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return err
	}
	if sdk.IsWasmAddress(recipient) {
		return ErrIsNotEvmAddr
	}

	if msg.Amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", msg.Amount)
	}
	return nil
}

func (msg MsgSendToEvm) GetSignBytes() []byte {
	panic("MsgSendToEvm can not be sign beacuse it can not exist in tx. It only exist in wasm call")
}

func (msg MsgSendToEvm) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
