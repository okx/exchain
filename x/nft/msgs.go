package nft

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	// TypeMsgSend nft message types
	TypeMsgSend = "send"
)

var _ sdk.MsgProtoAdapter = &MsgSend{}

// ValidateBasic implements the Msg.ValidateBasic method.
func (m MsgSend) ValidateBasic() error {
	if err := ValidateClassID(m.ClassId); err != nil {
		return sdkerrors.Wrapf(ErrInvalidID, "Invalid class id (%s)", m.ClassId)
	}

	if err := ValidateNFTID(m.Id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidID, "Invalid nft id (%s)", m.Id)
	}

	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", m.Sender)
	}

	_, err = sdk.AccAddressFromBech32(m.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid receiver address (%s)", m.Receiver)
	}
	return nil
}

// GetSigners implements Msg
func (m MsgSend) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{signer}
}

func (m MsgSend) Route() string {
	return "nft"
}

func (m MsgSend) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (msg MsgSend) GetSignBytes() []byte {
	panic("nft messages do not support amino")
}
