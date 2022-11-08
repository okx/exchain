package types

import (
	fmt "fmt"
	"strings"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	proto "github.com/gogo/protobuf/proto"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg             = &MsgRegisterAccount{}
	_ sdk.HeightSensitive = MsgRegisterAccount{}
	_ sdk.Msg             = &MsgSubmitTx{}
	_ sdk.HeightSensitive = MsgSubmitTx{}

	_ codectypes.UnpackInterfacesMessage = MsgSubmitTx{}
)

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(owner, connectionID, version string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionID,
		Version:      version,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

func (m MsgRegisterAccount) Route() string {
	return RouterKey
}

func (m MsgRegisterAccount) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgRegisterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// NewMsgSubmitTx creates and returns a new MsgSubmitTx instance
func NewMsgSubmitTx(sdkMsg sdk.Msg, connectionID, owner string) (*MsgSubmitTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitTx{
		ConnectionId: connectionID,
		Owner:        owner,
		Msg:          any,
	}, nil
}

// PackTxMsgAny marshals the sdk.Msg payload to a protobuf Any type
func PackTxMsgAny(sdkMsg sdk.Msg) (*codectypes.Any, error) {
	msg, ok := sdkMsg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", sdkMsg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return any, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsg sdk.Msg

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitTx) GetTxMsg() sdk.MsgAdapter {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.MsgAdapter)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

func (m MsgSubmitTx) Route() string {
	return RouterKey
}

func (m MsgSubmitTx) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m *MsgSubmitTx) GetSignBytes() []byte {
	panic("MsgSubmitTx messages do not support amino")
}

//////////
func (msg MsgRegisterAccount) ValidWithHeight(h int64) error {
	return common.MsgNotSupportBeforeHeight(&msg, h)
}

func (msg MsgSubmitTx) ValidWithHeight(h int64) error {
	return common.MsgNotSupportBeforeHeight(&msg, h)
}
