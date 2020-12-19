package types

import (
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	typeMsgDeposit           = "deposit"
	typeMsgWithdraw          = "withdraw"
	typeMsgTransferOwnership = "transferOwnership"
	typeMsgUpdateOperator    = "updateOperator"
	typeMsgCreateOperator    = "createOperator"
)

// MsgList - high level transaction of the dex module
type MsgList struct {
	Owner      sdk.AccAddress `json:"owner"`
	ListAsset  string         `json:"list_asset"`  //  Symbol of asset listed on Dex.
	QuoteAsset string         `json:"quote_asset"` //  Symbol of asset quoted by asset listed on Dex.
	InitPrice  sdk.Dec        `json:"init_price"`
}

// NewMsgList creates a new MsgList
func NewMsgList(owner sdk.AccAddress, listAsset, quoteAsset string, initPrice sdk.Dec) MsgList {
	return MsgList{
		Owner:      owner,
		ListAsset:  listAsset,
		QuoteAsset: quoteAsset,
		InitPrice:  initPrice,
	}
}

// Route Implements Msg
func (msg MsgList) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgList) Type() string { return "list" }

// ValidateBasic Implements Msg
func (msg MsgList) ValidateBasic() sdk.Error {
	if msg.ListAsset == msg.QuoteAsset {
		return ErrInvalidCoins()
	}

	if !msg.InitPrice.IsPositive() {
		return ErrInitPriceIsNotPositive()
	}

	if msg.Owner.Empty() {
		return ErrInvalidAddress(msg.Owner.String())
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgDeposit - high level transaction of the dex module
type MsgDeposit struct {
	Product   string         `json:"product"`   // product for trading pair in full name of the tokens
	Amount    sdk.SysCoin    `json:"amount"`    // Coins to add to the deposit
	Depositor sdk.AccAddress `json:"depositor"` // Address of the depositor
}

// NewMsgDeposit creates a new MsgDeposit
func NewMsgDeposit(product string, amount sdk.SysCoin, depositor sdk.AccAddress) MsgDeposit {
	return MsgDeposit{product, amount, depositor}
}

// Route Implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDeposit) Type() string { return typeMsgDeposit }

// ValidateBasic Implements Msg
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return ErrInvalidCoins()
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
	return []sdk.AccAddress{msg.Depositor}
}

// MsgWithdraw - high level transaction of the dex module
type MsgWithdraw struct {
	Product   string         `json:"product"`   // product for trading pair in full name of the tokens
	Amount    sdk.SysCoin    `json:"amount"`    // Coins to add to the deposit
	Depositor sdk.AccAddress `json:"depositor"` // Address of the depositor
}

// NewMsgWithdraw creates a new MsgWithdraw
func NewMsgWithdraw(product string, amount sdk.SysCoin, depositor sdk.AccAddress) MsgWithdraw {
	return MsgWithdraw{product, amount, depositor}
}

// Route Implements Msg
func (msg MsgWithdraw) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgWithdraw) Type() string { return typeMsgWithdraw }

// ValidateBasic Implements Msg
func (msg MsgWithdraw) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return ErrInvalidCoins()
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
	return []sdk.AccAddress{msg.Depositor}
}

// MsgTransferOwnership - high level transaction of the dex module
type MsgTransferOwnership struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Product     string         `json:"product"`
	//ToSignature auth.StdSignature `json:"to_signature"`
}

// NewMsgTransferOwnership create a new MsgTransferOwnership
func NewMsgTransferOwnership(from, to sdk.AccAddress, product string) MsgTransferOwnership {
	return MsgTransferOwnership{
		FromAddress: from,
		ToAddress:   to,
		Product:     product,
	}
}

// Route Implements Msg
func (msg MsgTransferOwnership) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferOwnership) Type() string { return typeMsgTransferOwnership }

// ValidateBasic Implements Msg
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return ErrInvalidAddress("missing sender address")
	}

	if msg.ToAddress.Empty() {
		return ErrInvalidAddress("missing recipient address")
	}

	if msg.Product == "" {
		return ErrTokenPairIsRequired()
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

//MsgConfirmOwnership - high level transaction of the coin module
type MsgConfirmOwnership struct {
	Product string         `json:"product"`
	Address sdk.AccAddress `json:"new_owner"`
}

func NewMsgConfirmOwnership(newOwner sdk.AccAddress, product string) MsgConfirmOwnership {
	return MsgConfirmOwnership{
		Product: product,
		Address: newOwner,
	}
}

func (msg MsgConfirmOwnership) Route() string { return RouterKey }

func (msg MsgConfirmOwnership) Type() string { return "confirm" }

func (msg MsgConfirmOwnership) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrInvalidAddress("failed to check MsgConfirmOwnership msg because miss sender address")
	}
	if len(msg.Product) == 0 {
		return ErrTokenPairIsRequired()
	}
	return nil
}

func (msg MsgConfirmOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgConfirmOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgCreateOperator register a new DEXOperator or update it
// Addr represent an DEXOperator
// if DEXOperator not exist, register a new DEXOperator
// else update Website or HandlingFeeAddress
type MsgCreateOperator struct {
	Owner              sdk.AccAddress `json:"owner"`
	Website            string         `json:"website"`
	HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
}

// NewMsgCreateOperator creates a new MsgCreateOperator
func NewMsgCreateOperator(website string, owner, handlingFeeAddress sdk.AccAddress) MsgCreateOperator {
	if handlingFeeAddress.Empty() {
		handlingFeeAddress = owner
	}
	return MsgCreateOperator{owner, strings.TrimSpace(website), handlingFeeAddress}
}

// Route Implements Msg
func (msg MsgCreateOperator) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgCreateOperator) Type() string { return typeMsgCreateOperator }

// ValidateBasic Implements Msg
func (msg MsgCreateOperator) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return ErrInvalidAddress("missing owner address")
	}
	if msg.HandlingFeeAddress.Empty() {
		return ErrInvalidAddress("missing handling fee address")
	}
	return checkWebsite(msg.Website)
}

// GetSignBytes Implements Msg
func (msg MsgCreateOperator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgCreateOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgUpdateOperator register a new DEXOperator or update it
// Addr represent an DEXOperator
// if DEXOperator not exist, register a new DEXOperator
// else update Website or HandlingFeeAddress
type MsgUpdateOperator struct {
	Owner              sdk.AccAddress `json:"owner"`
	Website            string         `json:"website"`
	HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
}

// NewMsgUpdateOperator creates a new MsgUpdateOperator
func NewMsgUpdateOperator(website string, owner, handlingFeeAddress sdk.AccAddress) MsgUpdateOperator {
	if handlingFeeAddress.Empty() {
		handlingFeeAddress = owner
	}
	return MsgUpdateOperator{owner, strings.TrimSpace(website), handlingFeeAddress}
}

// Route Implements Msg
func (msg MsgUpdateOperator) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgUpdateOperator) Type() string { return typeMsgUpdateOperator }

// ValidateBasic Implements Msg
func (msg MsgUpdateOperator) ValidateBasic() sdk.Error {
	if msg.HandlingFeeAddress.Empty() {
		return ErrInvalidAddress("missing handling fee address")
	}
	return checkWebsite(msg.Website)
}

// GetSignBytes Implements Msg
func (msg MsgUpdateOperator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg
func (msg MsgUpdateOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func checkWebsite(website string) sdk.Error {
	if len(website) == 0 {
		return nil
	}
	if len(website) > 1024 {
		return ErrInvalidWebsiteLength(len(website), 1024)
	}
	u, err := url.Parse(website)
	if err != nil {
		return ErrInvalidWebsiteURL(err.Error())
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidWebsiteURL(fmt.Sprintf("got: %s, expected: http or https", u.Scheme))
	}
	return nil
}
