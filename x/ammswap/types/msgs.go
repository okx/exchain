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
	MinLiquidity    sdk.Dec        `json:"min_liquidity"`      // Minimum number of sender will mint if total pool token supply is greater than 0.
	MaxAmountTokenA sdk.DecCoin    `json:"max_amount_token_a"` // Maximum number of tokens deposited. Deposits max amount if total pool token supply is 0.
	AmountTokenB    sdk.DecCoin    `json:"amount_token_b"`     // Quote token amount
	Deadline        int64          `json:"deadline"`           // Time after which this transaction can no longer be executed.
	Sender          sdk.AccAddress `json:"sender"`             // Sender
}

// NewMsgAddLiquidity is a constructor function for MsgAddLiquidity
func NewMsgAddLiquidity(minLiquidity sdk.Dec, maxAmountTokenA, amountTokenB sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgAddLiquidity {
	return MsgAddLiquidity{
		MinLiquidity:    minLiquidity,
		MaxAmountTokenA: maxAmountTokenA,
		AmountTokenB:    amountTokenB,
		Deadline:        deadline,
		Sender:          sender,
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
	if !(msg.MaxAmountTokenA.IsPositive() && msg.AmountTokenB.IsPositive()) {
		return sdk.ErrUnknownRequest("token amount must be positive")
	}
	if !msg.MaxAmountTokenA.IsValid() {
		return sdk.ErrUnknownRequest("invalid MaxAmountTokenA")
	}
	if !msg.AmountTokenB.IsValid() {
		return sdk.ErrUnknownRequest("invalid AmountTokenB")
	}
	err := ValidateBaseAndQuoteTokenName(msg.MaxAmountTokenA.Denom, msg.AmountTokenB.Denom)
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
	return GetSwapTokenPairName(msg.MaxAmountTokenA.Denom, msg.AmountTokenB.Denom)
}

// MsgRemoveLiquidity burns pool tokens to withdraw okt and Tokens at current ratio.
type MsgRemoveLiquidity struct {
	Liquidity       sdk.Dec        `json:"liquidity"`        // Amount of pool token burned.
	MinAmountTokenA sdk.DecCoin    `json:"min_amount_token_a"`  // Minimum base amount.
	MinAmountTokenB sdk.DecCoin    `json:"min_amount_token_b"` // Minimum quote amount.
	Deadline        int64          `json:"deadline"`         // Time after which this transaction can no longer be executed.
	Sender          sdk.AccAddress `json:"sender"`           // Sender
}

// NewMsgRemoveLiquidity is a constructor function for MsgAddLiquidity
func NewMsgRemoveLiquidity(liquidity sdk.Dec, minAmountTokenA, minAmountTokenB sdk.DecCoin, deadline int64, sender sdk.AccAddress) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{
		Liquidity:       liquidity,
		MinAmountTokenA: minAmountTokenA,
		MinAmountTokenB: minAmountTokenB,
		Deadline:        deadline,
		Sender:          sender,
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
	if !msg.MinAmountTokenA.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinAmountTokenA")
	}
	if !msg.MinAmountTokenB.IsValid() {
		return sdk.ErrUnknownRequest("invalid MinAmountTokenB")
	}
	err := ValidateBaseAndQuoteTokenName(msg.MinAmountTokenA.Denom, msg.MinAmountTokenB.Denom)
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
	return GetSwapTokenPairName(msg.MinAmountTokenA.Denom, msg.MinAmountTokenB.Denom)
}

// MsgCreateExchange creates a new exchange with token
type MsgCreateExchange struct {
	NameTokenA string         `json:"name_token_a"` // Token
	NameTokenB string         `json:"name_token_b"`
	Sender     sdk.AccAddress `json:"sender"` // Sender
}

// NewMsgCreateExchange create a new exchange with token
func NewMsgCreateExchange(nameTokenA string, nameTokenB string, sender sdk.AccAddress) MsgCreateExchange {
	return MsgCreateExchange{
		NameTokenA: nameTokenA,
		NameTokenB: nameTokenB,
		Sender:     sender,
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
	err := ValidateBaseAndQuoteTokenName(msg.NameTokenA, msg.NameTokenB)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error())
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
	return GetSwapTokenPairName(msg.NameTokenA, msg.NameTokenB)
}

// MsgTokenToToken define the message for swap between token and DefaultBondDenom
type MsgTokenToToken struct {
	SoldTokenAmount      sdk.DecCoin    `json:"sold_token_amount"`       // Amount of Tokens sold.
	MinBoughtTokenAmount sdk.DecCoin    `json:"min_bought_token_amount"` // Minimum token purchased.
	TokenRoute           []string       `json:"token_route"`
	Deadline             int64          `json:"deadline"`  // Time after which this transaction can no longer be executed.
	Recipient            sdk.AccAddress `json:"recipient"` // Recipient address,transfer Tokens to recipient.default recipient is sender.
	Sender               sdk.AccAddress `json:"sender"`    // Sender
}

// NewMsgTokenToToken is a constructor function for MsgTokenOKTSwap
func NewMsgTokenToToken(
	soldTokenAmount, minBoughtTokenAmount sdk.DecCoin, route []string, deadline int64, recipient, sender sdk.AccAddress,
) MsgTokenToToken {
	return MsgTokenToToken{
		SoldTokenAmount:      soldTokenAmount,
		MinBoughtTokenAmount: minBoughtTokenAmount,
		TokenRoute:           route,
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
	if msg.SoldTokenAmount.Denom == msg.MinBoughtTokenAmount.Denom {
		return sdk.ErrUnknownRequest("NameTokenA should not equal to NameTokenB")
	}
	tokenList := msg.TokenRoute
	tokenList = append(tokenList, msg.SoldTokenAmount.Denom, msg.MinBoughtTokenAmount.Denom)
	for _, tokenName := range tokenList {
		err := ValidateSwapTokenName(tokenName)
		if err != nil {
			return sdk.ErrUnknownRequest(err.Error())
		}
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
