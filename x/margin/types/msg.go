package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgDexDeposit{}
var _ sdk.Msg = &MsgDexWithdraw{}
var _ sdk.Msg = &MsgDexSet{}
var _ sdk.Msg = &MsgDexSave{}
var _ sdk.Msg = &MsgDexReturn{}
var _ sdk.Msg = &MsgDeposit{}
var _ sdk.Msg = &MsgBorrow{}
var _ sdk.Msg = &MsgRepay{}
var _ sdk.Msg = &MsgWithdraw{}

// MsgDexDeposit - struct for dex depositing for a product
type MsgDexDeposit struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoin    `json:"amount"`
}

// NewMsgDexDeposit creates a new MsgDeposit instance
func NewMsgDexDeposit(address sdk.AccAddress, product string, amount sdk.DecCoin) MsgDexDeposit {
	return MsgDexDeposit{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDexDeposit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDexDeposit) Type() string { return "dexDeposit" }

// ValidateBasic Implements Msg
func (msg MsgDexDeposit) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because deposits only support %s token", sdk.DefaultBondDenom))
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDexDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDexDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDexWithdraw - struct for dex withdrawing from a product
type MsgDexWithdraw struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoin    `json:"amount"`
}

// NewMsgDexWithdraw creates a new MsgWithdraw instance
func NewMsgDexWithdraw(address sdk.AccAddress, product string, amount sdk.DecCoin) MsgDexWithdraw {
	return MsgDexWithdraw{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDexWithdraw) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDexWithdraw) Type() string { return "dexWithdraw" }

// ValidateBasic Implements Msg
func (msg MsgDexWithdraw) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to withdraws because deposits only support %s token", sdk.DefaultBondDenom))
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDexWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDexWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDexSet - struct for dex setting params for a product
type MsgDexSet struct {
	Address                sdk.AccAddress `json:"address"`
	Product                string         `json:"product"`
	MaxLeverage            int64          `json:"max-leverage"`
	BorrowRate             sdk.Dec        `json:"borrow-rate"`
	MaintenanceMarginRatio sdk.Dec        `json:"maintenance-margin-ratio"`
}

// NewMsgDexSet creates a new MsgDexSet instance
func NewMsgDexSet(address sdk.AccAddress, product string, maxLeverage int64, borrowRate sdk.Dec, maintenanceMarginRatio sdk.Dec) MsgDexSet {
	return MsgDexSet{
		Address:                address,
		Product:                product,
		MaxLeverage:            maxLeverage,
		BorrowRate:             borrowRate,
		MaintenanceMarginRatio: maintenanceMarginRatio,
	}
}

// Route Implements Msg
func (msg MsgDexSet) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDexSet) Type() string { return "dexSet" }

// ValidateBasic Implements Msg
func (msg MsgDexSet) ValidateBasic() sdk.Error {
	if msg.MaxLeverage < 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("invald max leverage:%d", msg.MaxLeverage))
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDexSet) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDexSet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDexSave - struct for dex saving  for a product
type MsgDexSave struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoins   `json:"amount"`
}

// NewMsgDexSave creates a new MsgDexSave instance
func NewMsgDexSave(address sdk.AccAddress, product string, amount sdk.DecCoins) MsgDexSave {
	return MsgDexSave{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDexSave) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDexSave) Type() string { return "dexSave" }

// ValidateBasic Implements Msg
func (msg MsgDexSave) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDexSave) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDexSave) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDexReturn - struct for dex returning from a product
type MsgDexReturn struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoins   `json:"amount"`
}

// NewMsgDexReturn creates a new NewMsgDexReturn instance
func NewMsgDexReturn(address sdk.AccAddress, product string, amount sdk.DecCoins) MsgDexReturn {
	return MsgDexReturn{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDexReturn) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDexReturn) Type() string { return "dexReturn" }

// ValidateBasic Implements Msg
func (msg MsgDexReturn) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDexReturn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgDexReturn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDeposit - struct for depositing for a product
type MsgDeposit struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoins   `json:"amount"`
}

// NewMsgDeposit creates a new NewMsgDeposit instance
func NewMsgDeposit(address sdk.AccAddress, product string, amount sdk.DecCoins) MsgDeposit {
	return MsgDeposit{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDeposit) Type() string { return "deposit" }

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

// MsgBorrow - struct for borrowing for a product
type MsgBorrow struct {
	Address  sdk.AccAddress `json:"address"`
	Product  string         `json:"product"`
	Amount   sdk.DecCoin    `json:"amount"`
	Leverage sdk.Dec        `json:"leverage"`
}

// NewMsgBorrow creates a new MsgBorrow instance
func NewMsgBorrow(address sdk.AccAddress, product string, amount sdk.DecCoin, leverage sdk.Dec) MsgBorrow {
	return MsgBorrow{
		Address:  address,
		Product:  product,
		Amount:   amount,
		Leverage: leverage,
	}
}

// Route Implements Msg
func (msg MsgBorrow) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBorrow) Type() string { return "borrow" }

// ValidateBasic Implements Msg
func (msg MsgBorrow) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgBorrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgRepay - struct for repaying for a product
type MsgRepay struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoins   `json:"amount"`
}

// NewMsgRepay creates a new MsgRepay instance
func NewMsgRepay(address sdk.AccAddress, product string, amount sdk.DecCoins) MsgRepay {
	return MsgRepay{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgRepay) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgRepay) Type() string { return "repay" }

// ValidateBasic Implements Msg
func (msg MsgRepay) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgRepay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgRepay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgWithdraw - struct for withdraying for a product
type MsgWithdraw struct {
	Address sdk.AccAddress `json:"address"`
	Product string         `json:"product"`
	Amount  sdk.DecCoins   `json:"amount"`
}

// NewMsgWithdraw creates a new MsgWithdraw instance
func NewMsgWithdraw(address sdk.AccAddress, product string, amount sdk.DecCoins) MsgWithdraw {
	return MsgWithdraw{
		Address: address,
		Product: product,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgWithdraw) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgWithdraw) Type() string { return "withdraw" }

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
