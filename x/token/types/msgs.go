// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	common "github.com/okex/okchain/x/common/types"
)

const (
	DescLenLimit   = 256
	MultiSendLimit = 1000

	// 90 billion
	TotalSupplyUpperbound = int64(9 * 1e10)
)

//
type MsgTokenIssue struct {
	Description    string         `json:"description"`
	Symbol         string         `json:"symbol"`
	OriginalSymbol string         `json:"original_symbol"`
	WholeName      string         `json:"whole_name"`
	TotalSupply    string         `json:"total_supply"`
	Owner          sdk.AccAddress `json:"owner"`
	Mintable       bool           `json:"mintable"`
}

func NewMsgTokenIssue(tokenDescription, symbol, originalSymbol, wholeName, totalSupply string, owner sdk.AccAddress, mintable bool) MsgTokenIssue {
	return MsgTokenIssue{
		Description:    tokenDescription,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		WholeName:      wholeName,
		TotalSupply:    totalSupply,
		Owner:          owner,
		Mintable:       mintable,
	}
}

func (msg MsgTokenIssue) Route() string { return RouterKey }

func (msg MsgTokenIssue) Type() string { return "issue" }

func (msg MsgTokenIssue) ValidateBasic() sdk.Error {
	// check owner
	if msg.Owner.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "owner")
	}

	// check original symbol
	if len(msg.OriginalSymbol) == 0 {
		return common.ErrEmptyOriginalSymbol(common.AssetCodespace)
	}
	if !ValidOriginalSymbol(msg.OriginalSymbol) {
		return common.ErrInvalidOriginalSymbol(common.AssetCodespace, msg.OriginalSymbol)
	}

	// check wholeName
	isValid := wholeNameValid(msg.WholeName)
	if !isValid {
		return common.ErrInvalidWholeName(common.AssetCodespace, msg.WholeName)
	}
	// check desc
	if len(msg.Description) > DescLenLimit {
		return common.ErrTokenDescriptionExceeds(common.AssetCodespace, DescLenLimit)
	}
	// check totalSupply
	totalSupply, err := sdk.NewDecFromStr(msg.TotalSupply)
	if err != nil {
		return err
	}
	if totalSupply.GT(sdk.NewDec(TotalSupplyUpperbound)) {
		return common.ErrTotalSupplyExceeds(common.AssetCodespace, totalSupply.String(), TotalSupplyUpperbound)
	}
	if totalSupply.LTE(sdk.ZeroDec()) {
		return common.ErrInvalidRequestParam(common.AssetCodespace, "total supply must be positive")
	}

	return nil
}

func (msg MsgTokenIssue) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTokenIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgTokenBurn struct {
	Amount sdk.DecCoin    `json:"amount"`
	Owner  sdk.AccAddress `json:"owner"`
}

func NewMsgTokenBurn(amount sdk.DecCoin, owner sdk.AccAddress) MsgTokenBurn {
	return MsgTokenBurn{
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenBurn) Route() string { return RouterKey }

func (msg MsgTokenBurn) Type() string { return "burn" }

func (msg MsgTokenBurn) ValidateBasic() sdk.Error {
	// check owner
	if msg.Owner.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "owner")
	}
	if !msg.Amount.IsValid() {
		return common.ErrInvalidCoins(common.AssetCodespace)
	}

	return nil
}

func (msg MsgTokenBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTokenBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgTokenMint struct {
	Amount sdk.DecCoin    `json:"amount"`
	Owner  sdk.AccAddress `json:"owner"`
}

func NewMsgTokenMint(amount sdk.DecCoin, owner sdk.AccAddress) MsgTokenMint {
	return MsgTokenMint{
		Amount: amount,
		Owner:  owner,
	}
}

func (msg MsgTokenMint) Route() string { return RouterKey }

func (msg MsgTokenMint) Type() string { return "mint" }

func (msg MsgTokenMint) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "owner")
	}

	amount := msg.Amount.Amount
	if amount.GT(sdk.NewDec(TotalSupplyUpperbound)) {
		return common.ErrMintageAmountExceeds(common.AssetCodespace, TotalSupplyUpperbound)
	}
	if !msg.Amount.IsValid() {
		return common.ErrInvalidCoins(common.AssetCodespace)
	}
	return nil
}

func (msg MsgTokenMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)

	return sdk.MustSortJSON(bz)
}

func (msg MsgTokenMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgMultiSend struct {
	From      sdk.AccAddress `json:"from"`
	Transfers []TransferUnit `json:"transfers"`
}

func NewMsgMultiSend(from sdk.AccAddress, transfers []TransferUnit) MsgMultiSend {
	return MsgMultiSend{
		From:      from,
		Transfers: transfers,
	}
}

func (msg MsgMultiSend) Route() string { return RouterKey }

func (msg MsgMultiSend) Type() string { return "multi-send" }

func (msg MsgMultiSend) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "sender")
	}

	// check transfers
	if len(msg.Transfers) > MultiSendLimit {
		return common.ErrTransfersLengthExceeds(common.AssetCodespace, MultiSendLimit)
	}
	for _, transfer := range msg.Transfers {
		if !(transfer.Coins.IsAllPositive() && transfer.Coins.IsValid()) {
			return common.ErrInvalidCoins(common.AssetCodespace)
		}

		if transfer.To.Empty() {
			return common.ErrInvalidAddress(common.AssetCodespace, "recipient")
		}
	}
	return nil
}

func (msg MsgMultiSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.DecCoins   `json:"amount"`
}

func NewMsgTokenSend(from, to sdk.AccAddress, coins sdk.DecCoins) MsgSend {
	return MsgSend{
		FromAddress: from,
		ToAddress:   to,
		Amount:      coins,
	}
}

func (msg MsgSend) Route() string { return RouterKey }

func (msg MsgSend) Type() string { return "send" }

func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "sender")
	}
	if msg.ToAddress.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "recipient")
	}
	if !(msg.Amount.IsValid() && msg.Amount.IsAllPositive()) {
		return common.ErrInvalidCoins(common.AssetCodespace)
	}

	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgTransferOwnership - high level transaction of the coin module
type MsgTransferOwnership struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
	//FromSignature auth.StdSignature `json:"from_signature"`
	ToSignature auth.StdSignature `json:"to_signature"`
}

func NewMsgTransferOwnership(from, to sdk.AccAddress, symbol string) MsgTransferOwnership {
	return MsgTransferOwnership{
		FromAddress: from,
		ToAddress:   to,
		Symbol:      symbol,
		//FromSignature: auth.StdSignature{},
		ToSignature: auth.StdSignature{},
	}
}

func (msg MsgTransferOwnership) Route() string { return RouterKey }

func (msg MsgTransferOwnership) Type() string { return "transfer" }

func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "sender")
	}
	if msg.ToAddress.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "recipient")
	}
	if len(msg.Symbol) == 0 {
		return common.ErrEmptySymbol(common.AssetCodespace)
	}

	if sdk.ValidateDenom(msg.Symbol) != nil {
		return common.ErrInvalidSymbol(common.AssetCodespace, msg.Symbol)
	}

	if !msg.checkMultiSign() {
		return common.ErrInvalidMultisignCheck(common.AssetCodespace)
	}
	return nil
}

func (msg MsgTransferOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

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

type MsgTokenModify struct {
	Owner                 sdk.AccAddress `json:"owner"`
	Symbol                string         `json:"symbol"`
	Description           string         `json:"description"`
	WholeName             string         `json:"whole_name"`
	IsDescriptionModified bool           `json:"description_modified"`
	IsWholeNameModified   bool           `json:"whole_name_modified"`
}

func NewMsgTokenModify(symbol, desc, wholeName string, isDescEdit, isWholeNameEdit bool, owner sdk.AccAddress) MsgTokenModify {
	return MsgTokenModify{
		Symbol:                symbol,
		IsDescriptionModified: isDescEdit,
		Description:           desc,
		IsWholeNameModified:   isWholeNameEdit,
		WholeName:             wholeName,
		Owner:                 owner,
	}
}

func (msg MsgTokenModify) Route() string { return RouterKey }

func (msg MsgTokenModify) Type() string { return "edit" }

func (msg MsgTokenModify) ValidateBasic() sdk.Error {
	// check owner
	if msg.Owner.Empty() {
		return common.ErrInvalidAddress(common.AssetCodespace, "owner")
	}
	// check symbol
	if len(msg.Symbol) == 0 {
		return common.ErrEmptySymbol(common.AssetCodespace)
	}
	if sdk.ValidateDenom(msg.Symbol) != nil {
		return common.ErrInvalidSymbol(common.AssetCodespace, msg.Symbol)
	}
	// check wholeName
	if msg.IsWholeNameModified {
		isValid := wholeNameValid(msg.WholeName)
		if !isValid {
			return common.ErrInvalidWholeName(common.AssetCodespace, msg.WholeName)
		}
	}
	// check desc
	if msg.IsDescriptionModified {
		if len(msg.Description) > DescLenLimit {
			return common.ErrTokenDescriptionExceeds(common.AssetCodespace, DescLenLimit)
		}
	}
	return nil
}

func (msg MsgTokenModify) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTokenModify) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
