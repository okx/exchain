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
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrMsgSendToEvm(err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return ErrMsgSendToEvm(err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return ErrMsgSendToEvm(err.Error())
	}

	if msg.Amount.IsNegative() {
		return ErrMsgSendToEvm(fmt.Sprintf("negative coin amount: %v", msg.Amount))
	}
	return nil
}

func (msg MsgSendToEvm) GetSignBytes() []byte {
	panic(fmt.Errorf("MsgSendToEvm can not be sign beacuse it can not exist in tx. It only exist in wasm call"))
}

func (msg MsgSendToEvm) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err)
	}
	return []sdk.AccAddress{senderAddr}
}

func (msg MsgCallToEvm) Route() string {
	return RouterKey
}

func (msg MsgCallToEvm) Type() string {
	return WasmEvent2EvmMsgName
}

func (msg MsgCallToEvm) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrMsgSendToEvm(err.Error())
	}
	if !sdk.IsWasmAddress(sender) {
		return ErrIsNotWasmAddr
	}

	contract, err := sdk.AccAddressFromBech32(msg.Evmaddr)
	if err != nil {
		return ErrMsgSendToEvm(err.Error())
	}
	if sdk.IsWasmAddress(contract) {
		return ErrIsNotEvmAddr
	}

	if msg.Value.IsNegative() {
		return ErrMsgSendToEvm(fmt.Sprintf("negative value %v", msg.Value))
	}
	return nil
}

func (msg MsgCallToEvm) GetSignBytes() []byte {
	panic(fmt.Errorf("WasmEvent2Evm can not be sign beacuse it can not exist in tx. It only exist in wasm call"))
}

func (msg MsgCallToEvm) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err)
	}
	return []sdk.AccAddress{senderAddr}
}
