package types

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	commonType "github.com/okex/okchain/x/common/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	typeMsgDeposit           = "deposit"
	typeMsgWithdraw          = "withdraw"
	typeMsgTransferOwnership = "transferOwnership"
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

// MsgDelist - high level transaction of the dex module
type MsgDelist struct {
	Owner   sdk.AccAddress `json:"owner"`
	Product string         `json:"product"`
}

// NewMsgDelist creates a new MsgDelist
func NewMsgDelist(owner sdk.AccAddress, product string) MsgDelist {
	return MsgDelist{
		Owner:   owner,
		Product: product,
	}
}

// Route Implements Msg
func (msg MsgDelist) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDelist) Type() string { return "delist" }

// ValidateBasic Implements Msg
func (msg MsgDelist) ValidateBasic() sdk.Error {
	if msg.Product == "" || len(msg.Product) == 0 {
		return commonType.ErrEmptyProduct(commonType.SpotCodespace)
	}
	return nil
}

// GetSignBytes Implements Msg
func (msg MsgDelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgDeposit - high level transaction of the dex module
type MsgDeposit struct {
	Product   string         `json:"product"`   // product for trading pair in full name of the tokens
	Amount    sdk.DecCoin    `json:"amount"`    // Coins to add to the deposit
	Depositor sdk.AccAddress `json:"depositor"` // Address of the depositor
}

// NewMsgDeposit creates a new MsgDeposit
func NewMsgDeposit(product string, amount sdk.DecCoin, depositor sdk.AccAddress) MsgDeposit {
	return MsgDeposit{product, amount, depositor}
}

// Route Implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgDeposit) Type() string { return typeMsgDeposit }

// ValidateBasic Implements Msg
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return commonType.ErrInvalidAddress(commonType.SpotCodespace, "depositor")
	}
	if !msg.Amount.IsValid() {
		return commonType.ErrInvalidCoins(commonType.SpotCodespace)
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
	Amount    sdk.DecCoin    `json:"amount"`    // Coins to add to the deposit
	Depositor sdk.AccAddress `json:"depositor"` // Address of the depositor
}

// NewMsgWithdraw creates a new MsgWithdraw
func NewMsgWithdraw(product string, amount sdk.DecCoin, depositor sdk.AccAddress) MsgWithdraw {
	return MsgWithdraw{product, amount, depositor}
}

// Route Implements Msg
func (msg MsgWithdraw) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgWithdraw) Type() string { return typeMsgWithdraw }

// ValidateBasic Implements Msg
func (msg MsgWithdraw) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return commonType.ErrInvalidAddress(commonType.SpotCodespace, "depositor")
	}
	if !msg.Amount.IsValid() {
		return commonType.ErrInvalidCoins(commonType.SpotCodespace)
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
	FromAddress sdk.AccAddress    `json:"from_address"`
	ToAddress   sdk.AccAddress    `json:"to_address"`
	Product     string            `json:"product"`
	ToSignature auth.StdSignature `json:"to_signature"`
}

// NewMsgTransferOwnership create a new MsgTransferOwnership
func NewMsgTransferOwnership(from, to sdk.AccAddress, product string) MsgTransferOwnership {
	return MsgTransferOwnership{
		FromAddress: from,
		ToAddress:   to,
		Product:     product,
		ToSignature: auth.StdSignature{},
	}
}

// Route Implements Msg
func (msg MsgTransferOwnership) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferOwnership) Type() string { return typeMsgTransferOwnership }

// ValidateBasic Implements Msg
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return commonType.ErrInvalidAddress(commonType.SpotCodespace, "sender")
	}

	if msg.ToAddress.Empty() {
		return commonType.ErrInvalidAddress(commonType.SpotCodespace, "recipient")
	}

	if msg.Product == "" {
		return commonType.ErrEmptyProduct(commonType.SpotCodespace)
	}

	if !msg.checkMultiSign() {
		return commonType.ErrInvalidMultisignCheck(commonType.SpotCodespace)
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

func (msg MsgTransferOwnership) checkMultiSign() bool {
	// check pubkey
	if msg.ToSignature.PubKey == nil {
		return false
	}

	if !sdk.AccAddress(msg.ToSignature.PubKey.Address()).Equals(msg.ToAddress) {
		return false
	}

	// check multisign
	toSignature := msg.ToSignature
	msg.ToSignature = auth.StdSignature{}
	toValid := toSignature.VerifyBytes(msg.GetSignBytes(), toSignature.Signature)
	return toValid
}
