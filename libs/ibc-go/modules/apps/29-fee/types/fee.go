package types

import (
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
)

// NewPacketFee creates and returns a new PacketFee struct including the incentivization fees, refund addres and relayers
func NewPacketFee(fee Fee, refundAddr string, relayers []string) PacketFee {
	return PacketFee{
		Fee:           fee,
		RefundAddress: refundAddr,
		Relayers:      relayers,
	}
}

// Validate performs basic stateless validation of the associated PacketFee
func (p PacketFee) Validate() error {
	_, err := sdk.AccAddressFromBech32(p.RefundAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to convert RefundAddress into sdk.AccAddress")
	}

	// enforce relayers are not set
	if len(p.Relayers) != 0 {
		return ErrRelayersNotEmpty
	}

	if err := p.Fee.Validate(); err != nil {
		return err
	}

	return nil
}

// NewPacketFees creates and returns a new PacketFees struct including a list of type PacketFee
func NewPacketFees(packetFees []PacketFee) PacketFees {
	return PacketFees{
		PacketFees: packetFees,
	}
}

// NewIdentifiedPacketFees creates and returns a new IdentifiedPacketFees struct containing a packet ID and packet fees
func NewIdentifiedPacketFees(packetID channeltypes.PacketId, packetFees []PacketFee) IdentifiedPacketFees {
	return IdentifiedPacketFees{
		PacketId:   packetID,
		PacketFees: packetFees,
	}
}

// NewFee creates and returns a new Fee struct encapsulating the receive, acknowledgement and timeout fees as sdk.Coins
func NewFee(recvFee, ackFee, timeoutFee sdk.CoinAdapters) Fee {
	return Fee{
		RecvFee:    recvFee,
		AckFee:     ackFee,
		TimeoutFee: timeoutFee,
	}
}

// Total returns the total amount for a given Fee
func (f Fee) Total() sdk.CoinAdapters {
	return f.RecvFee.Add(f.AckFee...).Add(f.TimeoutFee...)
}

// Validate asserts that each Fee is valid and all three Fees are not empty or zero
func (fee Fee) Validate() error {
	var errFees []string
	if !fee.AckFee.IsValid() {
		errFees = append(errFees, "ack fee invalid")
	}
	if !fee.RecvFee.IsValid() {
		errFees = append(errFees, "recv fee invalid")
	}
	if !fee.TimeoutFee.IsValid() {
		errFees = append(errFees, "timeout fee invalid")
	}

	if len(errFees) > 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "contains invalid fees: %s", strings.Join(errFees, " , "))
	}

	// if all three fee's are zero or empty return an error
	if fee.AckFee.IsZero() && fee.RecvFee.IsZero() && fee.TimeoutFee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "all fees are zero")
	}

	return nil
}
