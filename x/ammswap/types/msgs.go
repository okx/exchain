package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
)

// PoolSwap message types and routes
const (
	TypeMsgAddLiquidity = "add_liquidity"
	TypeMsgTokenSwap    = "token_swap"
)

// MsgAddLiquidity Deposit quote_amount and base_amount at current ratio to mint pool tokens.
type MsgAddLiquidity struct {
	MinLiquidity  sdk.Dec        `json:"min_liquidity"`   // Minimum number of sender will mint if total pool token supply is greater than 0.
	MaxBaseAmount sdk.DecCoin    `json:"max_base_amount"` // Maximum number of tokens deposited. Deposits max amount if total pool token supply is 0.
	QuoteAmount   sdk.DecCoin    `json:"quote_amount"`    // Quote token amount
	Deadline      int64          `json:"deadline"`        // Time after which this transaction can no longer be executed.
	Sender        sdk.AccAddress `json:"sender"`          // Sender
}

// NewMsgAddLiquidity is a constructor function for MsgAddLiquidity
func NewMsgAddLiquidity(minLiquidity sdk.Dec, maxBaseAmount, quoteAmount sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgAddLiquidity {
	return MsgAddLiquidity{
		MinLiquidity:  minLiquidity,
		MaxBaseAmount: maxBaseAmount,
		QuoteAmount:   quoteAmount,
		Deadline:      deadline,
		Sender:        sender,
	}
}

// Route should return the name of the module
func (msg MsgAddLiquidity) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddLiquidity) Type() string { return "add_liquidity" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddLiquidity) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if !(msg.MaxBaseAmount.IsPositive() && msg.QuoteAmount.IsPositive()) {
		return sdk.ErrUnknownRequest("token amount must be positive")
	}
	if !msg.MaxBaseAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MaxBaseAmount")
	}
	if !msg.QuoteAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid QuoteAmount")
	}
	if msg.QuoteAmount.Denom != common.NativeToken {
		return sdk.ErrUnknownRequest("quote token only supports " + common.NativeToken)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSwapTokenPair defines token pair
func (msg MsgAddLiquidity) GetSwapTokenPair() string {
	return msg.MaxBaseAmount.Denom + "_" + msg.QuoteAmount.Denom
}

// MsgRemoveLiquidity burns pool tokens to withdraw okt and Tokens at current ratio.
type MsgRemoveLiquidity struct {
	Liquidity      sdk.Dec        `json:"liquidity"`        // Amount of pool token burned.
	MinBaseAmount  sdk.DecCoin    `json:"min_base_amount"`  // Minimum base amount.
	MinQuoteAmount sdk.DecCoin    `json:"min_quote_amount"` // Minimum quote amount.
	Deadline       int64          `json:"deadline"`         // Time after which this transaction can no longer be executed.
	Sender         sdk.AccAddress `json:"sender"`           // Sender
}

// NewMsgRemoveLiquidity is a constructor function for MsgAddLiquidity
func NewMsgRemoveLiquidity(liquidity sdk.Dec, minBaseAmount, minQuoteAmount sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{
		Liquidity:      liquidity,
		MinBaseAmount:  minBaseAmount,
		MinQuoteAmount: minQuoteAmount,
		Deadline:       deadline,
		Sender:         sender,
	}
}

// Route should return the name of the module
func (msg MsgRemoveLiquidity) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRemoveLiquidity) Type() string { return "remove_liquidity" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRemoveLiquidity) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if !(msg.Liquidity.IsPositive()) {
		return sdk.ErrUnknownRequest("token amount must be positive")
	}
	if !msg.MinBaseAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinBaseAmount")
	}
	if !msg.MinQuoteAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinQuoteAmount")
	}
	if msg.MinQuoteAmount.Denom != common.NativeToken {
		return sdk.ErrUnknownRequest("quote token only supports " + common.NativeToken)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSwapTokenPair defines token pair
func (msg MsgRemoveLiquidity) GetSwapTokenPair() string {
	return msg.MinBaseAmount.Denom + "_" + msg.MinQuoteAmount.Denom
}

// MsgCreateExchange creates a new exchange with token
type MsgCreateExchange struct {
	Token  string         `json:"token"`  // Token
	Sender sdk.AccAddress `json:"sender"` // Sender
}

// NewMsgCreateExchange create a new exchange with token
func NewMsgCreateExchange(token string, sender sdk.AccAddress) MsgCreateExchange {
	return MsgCreateExchange{
		Token:  token,
		Sender: sender,
	}
}

// Route should return the name of the module
func (msg MsgCreateExchange) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateExchange) Type() string { return "create_exchange" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateExchange) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if sdk.ValidateDenom(msg.Token) != nil || ValidatePoolTokenName(msg.Token) {
		return sdk.ErrUnknownRequest("invalid Token")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateExchange) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateExchange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgTokenToNativeToken define the message for swap between token and DefaultBondDenom
type MsgTokenToNativeToken struct {
	SoldTokenAmount      sdk.DecCoin    `json:"sold_token_amount"`       // Amount of Tokens sold.
	MinBoughtTokenAmount sdk.DecCoin    `json:"min_bought_token_amount"` // Minimum token purchased.
	Deadline             int64          `json:"deadline"`                // Time after which this transaction can no longer be executed.
	Recipient            sdk.AccAddress `json:"recipient"`               // Recipient address,transfer Tokens to recipient.default recipient is sender.
	Sender               sdk.AccAddress `json:"sender"`                  // Sender
}

// NewMsgTokenToNativeToken is a constructor function for MsgTokenOKTSwap
func NewMsgTokenToNativeToken(
	soldTokenAmount, minBoughtTokenAmount sdk.DecCoin, deadline int64, recipient, sender sdk.AccAddress,
) MsgTokenToNativeToken {
	return MsgTokenToNativeToken{
		SoldTokenAmount:      soldTokenAmount,
		MinBoughtTokenAmount: minBoughtTokenAmount,
		Deadline:             deadline,
		Recipient:            recipient,
		Sender:               sender,
	}
}

// Route should return the name of the module
func (msg MsgTokenToNativeToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTokenToNativeToken) Type() string { return TypeMsgTokenSwap }

// ValidateBasic runs stateless checks on the message
func (msg MsgTokenToNativeToken) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}

	if msg.SoldTokenAmount.Denom != sdk.DefaultBondDenom && msg.MinBoughtTokenAmount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("both token to sell and token to buy do not contain %s,"+
			" quote token only supports %s", sdk.DefaultBondDenom, sdk.DefaultBondDenom))
	}
	if !(msg.SoldTokenAmount.IsPositive()) {
		return sdk.ErrUnknownRequest("token amount must be positive")
	}
	if !msg.SoldTokenAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid SoldTokenAmount")
	}

	if !msg.MinBoughtTokenAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinBoughtTokenAmount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTokenToNativeToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTokenToNativeToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSwapTokenPair defines token pair
func (msg MsgTokenToNativeToken) GetSwapTokenPair() string {
	if msg.SoldTokenAmount.Denom == sdk.DefaultBondDenom {
		return msg.MinBoughtTokenAmount.Denom + "_" + msg.SoldTokenAmount.Denom
	}
	return msg.SoldTokenAmount.Denom + "_" + msg.MinBoughtTokenAmount.Denom
}
