package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	err := ValidateBaseAndQuoteAmount(msg.MaxBaseAmount.Denom, msg.QuoteAmount.Denom)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error())
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
func (msg MsgAddLiquidity) GetSwapTokenPairName() string {
	return GetSwapTokenPairName(msg.MaxBaseAmount.Denom, msg.QuoteAmount.Denom)
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
	err := ValidateBaseAndQuoteAmount(msg.MinBaseAmount.Denom, msg.MinQuoteAmount.Denom)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error())
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
func (msg MsgRemoveLiquidity) GetSwapTokenPairName() string {
	return GetSwapTokenPairName(msg.MinBaseAmount.Denom, msg.MinQuoteAmount.Denom)
}

// MsgCreateExchange creates a new exchange with token
type MsgCreateExchange struct {
	Token0Name string          `json:"token0_name"`
	Token1Name string          `json:"token1_name"`
	Sender          sdk.AccAddress `json:"sender"` // Sender
}

// NewMsgCreateExchange create a new exchange with token
func NewMsgCreateExchange(token0Name string, token1Name string, sender sdk.AccAddress) MsgCreateExchange {
	return MsgCreateExchange{
		Token0Name:  token0Name,
		Token1Name:  token1Name,
		Sender:         sender,
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

// GetSwapTokenPair defines token pair
func (msg MsgCreateExchange) GetSwapTokenPairName() string {
	return GetSwapTokenPairName(msg.Token0Name, msg.Token1Name)
}

// MsgTokenToToken define the message for swap between token and DefaultBondDenom
type MsgTokenToToken struct {
	SoldTokenAmount      sdk.DecCoin    `json:"sold_token_amount"`       // Amount of Tokens sold.
	MinBoughtTokenAmount sdk.DecCoin    `json:"min_bought_token_amount"` // Minimum token purchased.
	Deadline             int64          `json:"deadline"`                // Time after which this transaction can no longer be executed.
	Recipient            sdk.AccAddress `json:"recipient"`               // Recipient address,transfer Tokens to recipient.default recipient is sender.
	Sender               sdk.AccAddress `json:"sender"`                  // Sender
}

// NewMsgTokenToToken is a constructor function for MsgTokenOKTSwap
func NewMsgTokenToToken(
	soldTokenAmount, minBoughtTokenAmount sdk.DecCoin, deadline int64, recipient, sender sdk.AccAddress,
) MsgTokenToToken {
	return MsgTokenToToken{
		SoldTokenAmount:      soldTokenAmount,
		MinBoughtTokenAmount: minBoughtTokenAmount,
		Deadline:             deadline,
		Recipient:            recipient,
		Sender:               sender,
	}
}

// Route should return the name of the module
func (msg MsgTokenToToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTokenToToken) Type() string { return TypeMsgTokenSwap }

// ValidateBasic runs stateless checks on the message
func (msg MsgTokenToToken) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
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

	var baseAmountName, quoteAmountName string
	if msg.SoldTokenAmount.Denom < msg.MinBoughtTokenAmount.Denom {
		baseAmountName = msg.SoldTokenAmount.Denom
		quoteAmountName = msg.MinBoughtTokenAmount.Denom
	}else {
		baseAmountName = msg.MinBoughtTokenAmount.Denom
		quoteAmountName = msg.SoldTokenAmount.Denom
	}
	err := ValidateBaseAndQuoteAmount(baseAmountName, quoteAmountName)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTokenToToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTokenToToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSwapTokenPair defines token pair
func (msg MsgTokenToToken) GetSwapTokenPairName() string {
	return GetSwapTokenPairName(msg.MinBoughtTokenAmount.Denom, msg.SoldTokenAmount.Denom)
}