package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
)

type MsgAddLiquidity struct {
	MinLiquidity  sdk.Dec        `json:"min_liquidity"`   //Minimum number of sender will mint if total pool token supply is greater than 0.
	MaxBaseAmount sdk.DecCoin    `json:"max_base_amount"` //Maximum number of tokens deposited. Deposits max amount if total pool token supply is 0.
	QuoteAmount   sdk.DecCoin    `json:"quote_amount"`
	Deadline      int64          `json:"deadline"` //Time after which this transaction can no longer be executed.
	Sender        sdk.AccAddress `json:"sender"`   //sender
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
	if !(msg.MinLiquidity.IsPositive() && msg.MaxBaseAmount.IsPositive() && msg.QuoteAmount.IsPositive()) {
		return sdk.ErrUnknownRequest("tokens must be positive")
	}
	if !msg.MaxBaseAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MaxQuoteAmount")
	}
	if !msg.QuoteAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid BaseTokens")
	}
	if msg.QuoteAmount.Denom != common.NativeToken {
		return sdk.ErrUnknownRequest("quote token only supports okt")
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

func (msg MsgAddLiquidity) GetSwapTokenPair() string {
	return msg.MaxBaseAmount.Denom + "_" + msg.QuoteAmount.Denom
}

// MsgRemoveLiquidity burns pool tokens to withdraw okt and Tokens at current ratio.
type MsgRemoveLiquidity struct {
	Liquidity      sdk.Dec        `json:"liquidity"`        //Amount of pool token burned.
	MinBaseAmount  sdk.DecCoin    `json:"min_base_amount"`  //Minimum base amount.
	MinQuoteAmount sdk.DecCoin    `json:"min_quote_amount"` //Minimum quote amount.
	Deadline       int64          `json:"deadline"`         //Time after which this transaction can no longer be executed.
	Sender         sdk.AccAddress `json:"sender"`           //sender
}

// NewMsgRemoveLiquidity is a constructor function for MsgAddLiquidity
func NewMsgRemoveLiquidity(liquidity sdk.Dec, minBaseAmount, minQuoteAmount sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{
		Liquidity:  liquidity,
		MinBaseAmount: minBaseAmount,
		MinQuoteAmount:   minQuoteAmount,
		Deadline:      deadline,
		Sender:        sender,
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
	if !(msg.Liquidity.IsPositive() && msg.MinBaseAmount.IsPositive() && msg.MinBaseAmount.IsPositive()) {
		return sdk.ErrUnknownRequest("coins must be positive")
	}
	if !msg.MinBaseAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinBaseAmount")
	}
	if !msg.MinQuoteAmount.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinQuoteAmount")
	}
	if msg.MinQuoteAmount.Denom != common.NativeToken {
		return sdk.ErrUnknownRequest("quote token only supports okt")
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

func (msg MsgRemoveLiquidity) GetSwapTokenPair() string {
	return msg.MinBaseAmount.Denom + "_" + msg.MinQuoteAmount.Denom
}
