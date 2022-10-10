package typesadapter

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	okc_types "github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
)

var (
	_ txmsg.Msg = &MsgSend{}
	_ txmsg.Msg = &MsgMultiSend{}
	//_ token.TokenTransfer = &MsgSend{}
)

func (msg *MsgSend) ValidateBasic() error {
	_, err := types.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = types.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func (m *MsgSend) GetSigners() []types.AccAddress {
	from, err := types.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{from}
}

func (m *MsgSend) Route() string {
	return "token"
}

func (m *MsgSend) Type() string {
	return "send"
}

func (m MsgSend) GetSignBytes() []byte {
	return types.MustSortJSON(cdc.MustMarshalJSON(m))
}
func (m *MsgSend) GetFrom() sdk.AccAddress {
	from, err := types.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return from
}
func (m *MsgSend) GetTo() sdk.AccAddress {
	to, err := types.AccAddressFromBech32(m.ToAddress)
	if err != nil {
		panic(err)
	}
	return to
}
func (m *MsgSend) GetAmount() []sdk.DecCoin {
	convAmount := make([]sdk.DecCoin, 0)
	for _, am := range m.Amount {
		transferAmountDec := sdk.NewDecFromIntWithPrec(sdk.NewIntFromBigInt(am.Amount.BigInt()), sdk.Precision)
		convAmount = append(convAmount, sdk.NewDecCoinFromDec(am.Denom, transferAmountDec))
	}
	return convAmount
}

func (m *MsgSend) RulesFilter() (sdk.Msg, error) {
	msgSend := *m

	msgSend.Amount = m.Amount.Copy()
	for i, amount := range msgSend.Amount {
		if amount.Denom == sdk.DefaultIbcWei {
			msgSend.Amount[i].Denom = sdk.DefaultBondDenom
		} else if amount.Denom == sdk.DefaultBondDenom {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "MsgSend not support okt denom")
		}
	}
	return &msgSend, nil
}

func (m *MsgSend) Swap(ctx sdk.Context) (sdk.Msg, error) {
	for i, amount := range m.Amount {
		if amount.Denom == sdk.DefaultIbcWei {
			m.Amount[i].Denom = sdk.DefaultBondDenom
		} else if amount.Denom == sdk.DefaultBondDenom {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "MsgSend not support okt denom")
		}
	}

	return m, nil
}

func (msg *MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return okc_types.ErrNoInputs
	}
	if len(msg.Outputs) == 0 {
		return okc_types.ErrNoOutputs
	}
	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

func (m *MsgMultiSend) GetSigners() []types.AccAddress {
	froms := make([]types.AccAddress, 0)
	for i, _ := range m.Inputs {
		from, err := types.AccAddressFromBech32(m.Inputs[i].Address)
		if err != nil {
			panic(err)
		}
		froms = append(froms, from)
	}

	return froms
}

func (m *MsgMultiSend) Route() string {
	return "token"
}

func (m *MsgMultiSend) Type() string {
	return "multi-send"
}

func (m MsgMultiSend) GetSignBytes() []byte {
	return types.MustSortJSON(cdc.MustMarshalJSON(m))
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(in.Address)
	if err != nil {
		return err
	}

	if !in.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}

	if !in.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}

	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
//nolint:interfacer
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	return Input{
		Address: addr.String(),
		Coins:   sdk.CoinsToCoinAdapters(coins),
	}
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(out.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid output address (%s)", err)
	}

	if !out.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}

	if !out.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}

	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
//nolint:interfacer
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	return Output{
		Address: addr.String(),
		Coins:   sdk.CoinsToCoinAdapters(coins),
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) error {
	var totalIn, totalOut sdk.Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}

		inCoins := sdk.CoinAdaptersToCoins(in.Coins)
		totalIn = totalIn.Add(inCoins...)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}

		outCoins := sdk.CoinAdaptersToCoins(out.Coins)
		totalOut = totalOut.Add(outCoins...)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return okc_types.ErrInputOutputMismatch
	}

	return nil
}
