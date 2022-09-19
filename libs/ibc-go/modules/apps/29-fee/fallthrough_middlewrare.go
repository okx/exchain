package fee

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

var (
	_ porttypes.Middleware = (*FallThroughMiddleware)(nil)
)

type FallThroughMiddleware struct {
	h1 int64
	h2 int64

	left   porttypes.Middleware
	middle porttypes.Middleware
	right  porttypes.Middleware
}

// TODO, 添加ut
func NewFallThroughMiddleware(h1 int64, h2 int64, left porttypes.Middleware, middle porttypes.Middleware, right porttypes.Middleware) porttypes.Middleware {
	if h1 > 0 && h2 > 0 {
		if h1 > h2 {
			panic("illegal constructor parameter")
		}
	}
	ret := &FallThroughMiddleware{h1: h1, h2: h2, left: left, middle: middle, right: right}
	return ret
}

func (f *FallThroughMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
	} else {
		return f.left.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
	}
}

func (f *FallThroughMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version, counterpartyVersion string) (string, error) {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
	} else {
		return f.left.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
	}
}

func (f *FallThroughMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
	} else {
		return f.left.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
	}
}

func (f *FallThroughMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanOpenConfirm(ctx, portID, channelID)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanOpenConfirm(ctx, portID, channelID)
	} else {
		return f.left.OnChanOpenConfirm(ctx, portID, channelID)
	}
}

func (f *FallThroughMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanCloseInit(ctx, portID, channelID)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanCloseInit(ctx, portID, channelID)
	} else {
		return f.left.OnChanCloseInit(ctx, portID, channelID)
	}
}

func (f *FallThroughMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnChanCloseConfirm(ctx, portID, channelID)
	} else if f.higherThanH1(h) {
		return f.middle.OnChanCloseConfirm(ctx, portID, channelID)
	} else {
		return f.left.OnChanCloseConfirm(ctx, portID, channelID)
	}
}

func (f *FallThroughMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnRecvPacket(ctx, packet, relayer)
	} else if f.higherThanH1(h) {
		return f.middle.OnRecvPacket(ctx, packet, relayer)
	} else {
		return f.left.OnRecvPacket(ctx, packet, relayer)
	}
}

func (f *FallThroughMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	} else if f.higherThanH1(h) {
		return f.middle.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	} else {
		return f.left.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}
}

func (f *FallThroughMiddleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.OnTimeoutPacket(ctx, packet, relayer)
	} else if f.higherThanH1(h) {
		return f.middle.OnTimeoutPacket(ctx, packet, relayer)
	} else {
		return f.left.OnTimeoutPacket(ctx, packet, relayer)
	}
}

func (f *FallThroughMiddleware) NegotiateAppVersion(ctx sdk.Context, order channeltypes.Order, connectionID string, portID string, counterparty channeltypes.Counterparty, proposedVersion string) (version string, err error) {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.NegotiateAppVersion(ctx, order, connectionID, portID, counterparty, proposedVersion)
	} else if f.higherThanH1(h) {
		return f.middle.NegotiateAppVersion(ctx, order, connectionID, portID, counterparty, proposedVersion)
	} else {
		return f.left.NegotiateAppVersion(ctx, order, connectionID, portID, counterparty, proposedVersion)
	}
}

func (f *FallThroughMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.SendPacket(ctx, chanCap, packet)
	} else if f.higherThanH1(h) {
		return f.middle.SendPacket(ctx, chanCap, packet)
	} else {
		return f.left.SendPacket(ctx, chanCap, packet)
	}
}

func (f *FallThroughMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.WriteAcknowledgement(ctx, chanCap, packet, ack)
	} else if f.higherThanH1(h) {
		return f.middle.WriteAcknowledgement(ctx, chanCap, packet, ack)
	} else {
		return f.left.WriteAcknowledgement(ctx, chanCap, packet, ack)
	}
}

func (f *FallThroughMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	h := ctx.BlockHeight()
	if f.higherThanH2(h) {
		return f.right.GetAppVersion(ctx, portID, channelID)
	} else if f.higherThanH1(h) {
		return f.middle.GetAppVersion(ctx, portID, channelID)
	} else {
		return f.left.GetAppVersion(ctx, portID, channelID)
	}
}

//nolint
func (f *FallThroughMiddleware) higherThanH2(h int64) bool {
	if f.h2 == 0 {
		return false
	}
	return h >= f.h2
}

//nolint
func (f *FallThroughMiddleware) higherThanH1(h int64) bool {
	if f.h1 == 0 {
		return false
	}
	return h >= f.h1
}
