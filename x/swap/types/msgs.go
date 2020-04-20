package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
)

// Swap message types and routes
const (
	TypeMsgAddLiquidity = "add_liquidity"
	TypeMsgTokenOKTSwap = "token_okt_swap"
)

type MsgAddLiquidity struct {
	MinLiquidity  sdk.Dec        `json:"min_liquidity"`   //Minimum number of sender will mint if total pool token supply is greater than 0.
	MaxBaseTokens sdk.DecCoin    `json:"max_base_tokens"` //Maximum number of tokens deposited. Deposits max amount if total pool token supply is 0.
	QuoteTokens   sdk.DecCoin    `json:"quote_tokens"`
	Deadline      int64          `json:"deadline"` //Time after which this transaction can no longer be executed.
	Sender        sdk.AccAddress `json:"sender"`   //sender
}

// NewMsgAddLiquidity is a constructor function for MsgAddLiquidity
func NewMsgAddLiquidity(minLiquidity sdk.Dec, maxBaseTokens, quoteTokens sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgAddLiquidity {
	return MsgAddLiquidity{
		MinLiquidity:  minLiquidity,
		MaxBaseTokens: maxBaseTokens,
		QuoteTokens:   quoteTokens,
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
	if !(msg.MinLiquidity.IsPositive() && msg.MaxBaseTokens.IsPositive() && msg.QuoteTokens.IsPositive()) {
		return sdk.ErrUnknownRequest("tokens must be positive")
	}
	if !msg.MaxBaseTokens.IsValid() {
		return sdk.ErrUnknownRequest("invalid MaxQuoteTokens")
	}
	if !msg.QuoteTokens.IsValid() {
		return sdk.ErrUnknownRequest("invalid BaseTokens")
	}
	if msg.QuoteTokens.Denom != common.NativeToken {
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
	return msg.MaxBaseTokens.Denom + "_" + msg.QuoteTokens.Denom
}

// MsgTokenOKTSwap define the message for swap between token and DefaultBondDenom
type MsgTokenOKTSwap struct {
	SoldTokenAmount      sdk.DecCoin    `json:"sold_token_amount"`       //Amount of Tokens sold.
	MinBoughtTokenAmount sdk.DecCoin    `json:"min_bought_token_amount"` //Minimum token purchased.
	Deadline             int64          `json:"deadline"`                //Time after which this transaction can no longer be executed.
	Recipient            sdk.AccAddress `json:"recipient"`               //if give Recipient address,transfers Tokens to recipient.default recipient is sender
	Sender               sdk.AccAddress `json:"sender"`                  //sender
}

// NewMsgTokenOKTSwap is a constructor function for MsgTokenOKTSwap
func NewMsgTokenOKTSwap(
	soldTokenAmount, minBoughtTokenAmount sdk.DecCoin, deadline int64, recipient, sender sdk.AccAddress,
) MsgTokenOKTSwap {
	return MsgTokenOKTSwap{
		SoldTokenAmount:      soldTokenAmount,
		MinBoughtTokenAmount: minBoughtTokenAmount,
		Deadline:             deadline,
		Recipient:            recipient,
		Sender:               sender,
	}
}

// Route should return the name of the module
func (msg MsgTokenOKTSwap) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTokenOKTSwap) Type() string { return TypeMsgTokenOKTSwap }

// ValidateBasic runs stateless checks on the message
func (msg MsgTokenOKTSwap) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}

	if msg.SoldTokenAmount.Denom != sdk.DefaultBondDenom && msg.MinBoughtTokenAmount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("both token to sell and token to buy do not contain %s,"+
			" quote token only supports okt", sdk.DefaultBondDenom))
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
func (msg MsgTokenOKTSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTokenOKTSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgTokenOKTSwap) GetSwapTokenPair() string {
	if msg.SoldTokenAmount.Denom == sdk.DefaultBondDenom {
		return msg.MinBoughtTokenAmount.Denom + "_" + msg.SoldTokenAmount.Denom
	}
	return msg.SoldTokenAmount.Denom + "_" + msg.MinBoughtTokenAmount.Denom
}
