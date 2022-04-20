package typesadapter

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	okctypes "github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
)

var (
	_ txmsg.Msg = &MsgSend{}
	//_ sdk.Msg   = (*MsgSend)(nil)
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
	return "bank"
}

func (m *MsgSend) Type() string {
	return "send"
}

func (m *MsgSend) GetSignBytes() []byte {
	return types.MustSortJSON(okctypes.ModuleCdc.MustMarshalJSON(&m))
}
