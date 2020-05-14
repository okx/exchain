package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
/*
verify interface at compile time
var _ sdk.Msg = &Msg<Action>{}

Msg<Action> - struct for unjailing jailed validator
type Msg<Action> struct {
	ValidatorAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the validator operator
}

NewMsg<Action> creates a new Msg<Action> instance
func NewMsg<Action>(validatorAddr sdk.ValAddress) Msg<Action> {
	return Msg<Action>{
		ValidatorAddr: validatorAddr,
	}
}

const <action>Const = "<action>"

// nolint
func (msg Msg<Action>) Route() string { return RouterKey }
func (msg Msg<Action>) Type() string  { return <action>Const }
func (msg Msg<Action>) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

GetSignBytes gets the bytes for the message signer to sign on
func (msg Msg<Action>) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

ValidateBasic validity check for the AnteHandler
func (msg Msg<Action>) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address"
	}
	return nil
}
*/

type MsgMarginDeposit struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoin    `json:"amount"`
}

func NewMsgMarginDeposit(address sdk.AccAddress, product string, amount sdk.DecCoin) MsgMarginDeposit {
	return MsgMarginDeposit{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

func (msg MsgMarginDeposit) Route() string { return RouterKey }

func (msg MsgMarginDeposit) Type() string { return "margin-deposit" }

func (msg MsgMarginDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgMarginDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMarginDeposit) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

type MsgBorrow struct {
	Address  sdk.AccAddress `json:"address"`
	Product  string         `json:"product"`
	Amount   sdk.DecCoin    `json:"amount"`
	Leverage int            `json:"leverage"`
}

func NewMsgBorrow(address sdk.AccAddress, product string, amount sdk.DecCoin, leverage int) MsgBorrow {
	return MsgBorrow{
		Address:  address,
		Product:  product,
		Amount:   amount,
		Leverage: leverage,
	}
}

func (msg MsgBorrow) Route() string { return RouterKey }

func (msg MsgBorrow) Type() string { return "margin-borrow" }

func (msg MsgBorrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBorrow) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

type MsgRepay struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoin    `json:"amount"`
}

func NewMsgRepay(address sdk.AccAddress, product string, amount sdk.DecCoin) MsgRepay {
	return MsgRepay{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

func (msg MsgRepay) Route() string { return RouterKey }

func (msg MsgRepay) Type() string { return "margin-repay" }

func (msg MsgRepay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgRepay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRepay) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

type MsgWithdraw struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoin    `json:"amount"`
}

func NewMsgWithdraw(address sdk.AccAddress, product string, amount sdk.DecCoin) MsgWithdraw {
	return MsgWithdraw{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

func (msg MsgWithdraw) Route() string { return RouterKey }

func (msg MsgWithdraw) Type() string { return "margin-withdraw" }

func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgWithdraw) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}
