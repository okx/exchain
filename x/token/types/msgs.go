package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
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

// ValidateBasic Implements Msg.
func (msg MsgTokenIssue) ValidateBasic() sdk.Error {
	// check owner
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	// check original symbol
	if len(msg.OriginalSymbol) == 0 {
		return sdk.ErrUnknownRequest("failed to check issue msg because original symbol cannot be empty")
	}
	if !ValidOriginalSymbol(msg.OriginalSymbol) {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid original symbol: " + msg.OriginalSymbol)
	}

	// check wholeName
	isValid := wholeNameValid(msg.WholeName)
	if !isValid {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid wholename")
	}
	// check desc
	if len(msg.Description) > DescLenLimit {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid desc")
	}
	// check totalSupply
	totalSupply, err := sdk.NewDecFromStr(msg.TotalSupply)
	if err != nil {
		return err
	}
	if totalSupply.GT(sdk.NewDec(TotalSupplyUpperbound)) || totalSupply.LTE(sdk.ZeroDec()) {
		return sdk.ErrUnknownRequest("failed to check issue msg because invalid total supply")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenIssue) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTokenIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTokenBurn burn coins
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
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInsufficientCoins("failed to check burn msg because invalid Coins: " + msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
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
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	amount := msg.Amount.Amount
	if amount.GT(sdk.NewDec(TotalSupplyUpperbound)) {
		return sdk.ErrUnknownRequest("failed to check mint msg because invalid amount")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInsufficientCoins("failed to check mint msg because invalid Coins: " + msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)

	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
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
		return sdk.ErrInvalidAddress(msg.From.String())
	}

	// check transfers
	if len(msg.Transfers) > MultiSendLimit {
		return sdk.ErrUnknownRequest("failed to check multisend msg because restrictions on the number of transfers")
	}
	for _, transfer := range msg.Transfers {
		if !transfer.Coins.IsAllPositive() || !transfer.Coins.IsValid() {
			return sdk.ErrInvalidCoins("failed to check multisend msg because send amount must be positive")
		}

		if transfer.To.Empty() {
			return sdk.ErrInvalidAddress("failed to check multisend msg because address is empty, not valid")
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
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

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("failed to check send msg because miss sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("failed to check send msg because miss recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("failed to check send msg because send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("failed to check send msg because send amount must be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
	//return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
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

// Route Implements Msg.
func (msg MsgTransferOwnership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgTransferOwnership) Type() string { return "transfer" }

// ValidateBasic Implements Msg.
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("failed to check transferownership msg because miss sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("failed to check transferownership msg because miss recipient address")
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("failed to check transferownership msg because symbol cannot be empty")
	}

	if sdk.ValidateDenom(msg.Symbol) != nil {
		return sdk.ErrUnknownRequest("failed to check transferownership msg because invalid token symbol: " + msg.Symbol)
	}

	if !msg.checkMultiSign() {
		return sdk.ErrUnauthorized("failed to check transferownership msg because invalid multi signature")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
	//return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
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
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	// check symbol
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("failed to check modify msg because symbol cannot be empty")
	}
	if sdk.ValidateDenom(msg.Symbol) != nil {
		return sdk.ErrUnknownRequest("failed to check modify msg because invalid token symbol: " + msg.Symbol)
	}
	// check wholeName
	if msg.IsWholeNameModified {
		isValid := wholeNameValid(msg.WholeName)
		if !isValid {
			return sdk.ErrUnknownRequest("failed to check modify msg because invalid wholename")
		}
	}
	// check desc
	if msg.IsDescriptionModified {
		if len(msg.Description) > DescLenLimit {
			return sdk.ErrUnknownRequest("failed to check modify msg because invalid desc")
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTokenModify) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTokenModify) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
