package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	typeMsgDeposit  = "deposit"
	typeMsgWithdraw = "withdraw"
)

// MsgDeposit - struct for depositing to saving module
type MsgDeposit struct {
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.DecCoin    `json:"amount"`
}

// NewMsgDeposit creates a new MsgDeposit instance
func NewMsgDeposit(address sdk.AccAddress, amount sdk.DecCoin) MsgDeposit {
	return MsgDeposit{
		Address: address,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDeposit) Type() string { return typeMsgDeposit }

// ValidateBasic Implements Msg
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgWithdraw - struct for withdrawing from saving module
type MsgWithdraw struct {
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.DecCoin    `json:"amount"`
}

// NewMsgWithdraw creates a new MsgWithdraw instance
func NewMsgWithdraw(address sdk.AccAddress, amount sdk.DecCoin) MsgWithdraw {
	return MsgWithdraw{address, amount}
}

// Route Implements Msg
func (msg MsgWithdraw) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgWithdraw) Type() string { return typeMsgWithdraw }

// ValidateBasic Implements Msg
func (msg MsgWithdraw) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg
func (msg MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
