//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Verify interface at compile time
var _, _ sdk.Msg = &MsgSetWithdrawAddress{}, &MsgWithdrawValidatorCommission{}

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgSetWithdrawAddress struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	WithdrawAddress  sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
}

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		DelegatorAddress: delAddr,
		WithdrawAddress:  withdrawAddr,
	}
}

func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }
func (msg MsgSetWithdrawAddress) Type() string  { return "set_withdraw_address" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgSetWithdrawAddress) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.WithdrawAddress.Empty() {
		return ErrNilWithdrawAddr(DefaultCodespace)
	}
	return nil
}

// msg struct for validator withdraw
type MsgWithdrawValidatorCommission struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func NewMsgWithdrawValidatorCommission(valAddr sdk.ValAddress) MsgWithdrawValidatorCommission {
	return MsgWithdrawValidatorCommission{
		ValidatorAddress: valAddr,
	}
}

func (msg MsgWithdrawValidatorCommission) Route() string { return ModuleName }
func (msg MsgWithdrawValidatorCommission) Type() string  { return "withdraw_validator_commission" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorCommission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawValidatorCommission) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}
