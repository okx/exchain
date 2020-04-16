package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Route() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []sdk.AccAddress
}

type MsgAddLiquidity struct {
	MinLiquidity  sdk.Dec        `json:"min_liquidity"` //Minimum number of sender will mint if total pool token supply is greater than 0.
	MaxBaseTokens sdk.DecCoin    `json:"max_tokens"`    //Maximum number of tokens deposited. Deposits max amount if total pool token supply is 0.
	QuoteTokens   sdk.DecCoin    `json:"base_tokens"`
	Deadline      sdk.Uint       `json:"deadline"` //Time after which this transaction can no longer be executed.
	Sender        sdk.AccAddress `json:"sender"`   //sender
}

// NewMsgAddLiquidity is a constructor function for MsgAddLiquidity
func NewMsgAddLiquidity(minLiquidity sdk.Dec, maxBaseTokens, quoteTokens sdk.DecCoin, deadline sdk.Uint, sender sdk.AccAddress) MsgAddLiquidity {
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