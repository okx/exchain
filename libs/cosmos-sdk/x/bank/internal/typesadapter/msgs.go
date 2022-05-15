package typesadapter

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

var (
	_ txmsg.Msg = &MsgSend{}
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
