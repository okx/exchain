package types

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

var (
	_ ibcmsg.Msg = (*ProtobufMsgSubmitProposal)(nil)
)

func NewProtobufMsgSubmitProposal(content ContentAdapter, initialDeposit sdk.SysCoins, proposer sdk.AccAddress) ProtobufMsgSubmitProposal {
	any, err := PackContent(content)
	if nil != err {
		// cant happen
		panic(err)
	}
	cs := make(sdk.CoinAdapters, 0)
	for _, c := range initialDeposit {
		cs = append(cs, sdk.CoinAdapter{
			Denom:  c.Denom,
			Amount: sdk.NewIntFromBigInt(c.Amount.BigInt()),
		})
	}
	ret := ProtobufMsgSubmitProposal{
		Content:        any,
		InitialDeposit: cs,
		Proposer:       proposer.String(),
	}
	return ret
}

func PackContent(c ContentAdapter) (*codectypes.Any, error) {
	msg, ok := c.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", c)
	}

	content, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return content, nil
}

func UnPackContent(any *codectypes.Any) (ContentAdapter, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	clientState, ok := any.GetCachedValue().(ContentAdapter)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into content %T", any)
	}

	return clientState, nil
}

//nolint
func (msg ProtobufMsgSubmitProposal) Route() string { return RouterKey }
func (msg ProtobufMsgSubmitProposal) Type() string  { return TypeMsgSubmitProposal }

// Implements Msg.
func (msg ProtobufMsgSubmitProposal) ValidateBasic() sdk.Error {
	if msg.Content == nil {
		return ErrInvalidProposalContent("content is required")
	}
	cont, err := UnPackContent(msg.Content)
	if nil != err {
		return err
	}
	if cont.ProposalType() == ProposalTypeSoftwareUpgrade {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrInvalidProposalType(cont.ProposalType())
	}

	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if nil != err {
		return err
	}
	if proposer.Empty() {
		return ErrInvalidAddress(proposer.String())
	}

	if !sdk.CoinAdapters(msg.InitialDeposit).IsValid() {
		return ErrInvalidCoins()
	}
	if sdk.CoinAdapters(msg.InitialDeposit).IsAnyNegative() {
		return ErrInvalidCoins()
	}

	if len(msg.InitialDeposit) != 1 || msg.InitialDeposit[0].Denom != sdk.DefaultBondDenom || !sdk.CoinAdapters(msg.InitialDeposit).IsValid() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("must deposit %s but got %s", sdk.DefaultBondDenom,
			sdk.CoinAdapters(msg.InitialDeposit).String()))
	}

	if !IsValidProposalType(cont.ProposalType()) {
		return ErrInvalidProposalType(cont.ProposalType())
	}

	return cont.ValidateBasic()
}

// Implements Msg.
func (msg ProtobufMsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg ProtobufMsgSubmitProposal) GetSigners() []sdk.AccAddress {
	p, _ := sdk.AccAddressFromBech32(msg.Proposer)
	return []sdk.AccAddress{p}
}
