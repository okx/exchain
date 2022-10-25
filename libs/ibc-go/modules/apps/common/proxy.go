package common

import (
	"fmt"
	"os"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

var (
	_ porttypes.Middleware = (*HeightProxyMiddleware)(nil)
	_ porttypes.Middleware = (*DisaleProxyMiddleware)(nil)
)

type HeightProxyMiddleware struct {
	h        int64
	logger   log.Logger
	internal porttypes.Middleware
}

func NewHeightProxyMiddleware(h int64, module string, internal porttypes.Middleware) porttypes.Middleware {
	ret := &HeightProxyMiddleware{h: h, internal: internal}
	ret.logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", module)
	return ret
}

func (h *HeightProxyMiddleware) GetInternal() porttypes.Middleware {
	return h.internal
}
func (h *HeightProxyMiddleware) GetLogger() log.Logger {
	return h.logger
}
func (h *HeightProxyMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanOpenInit is disabled", "available", h.h, "now", ctx.BlockHeight())
		return "", ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (h *HeightProxyMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version, counterpartyVersion string) (string, error) {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanOpenTry is disabled", "available", h.h, "now", ctx.BlockHeight())
		return "", ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
}

func (h *HeightProxyMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanOpenAck is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (h *HeightProxyMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanOpenConfirm is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanOpenConfirm(ctx, portID, channelID)
}

func (h *HeightProxyMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanCloseInit is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanCloseInit(ctx, portID, channelID)
}

func (h *HeightProxyMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnChanCloseConfirm is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnChanCloseConfirm(ctx, portID, channelID)
}

func (h *HeightProxyMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnRecvPacket is disabled", "available", h.h, "now", ctx.BlockHeight())
		return channeltypes.NewErrorAcknowledgement(fmt.Sprintf("OnRecvPacket is disabled until:%d", h.h))
	}
	return h.internal.OnRecvPacket(ctx, packet, relayer)
}

func (h *HeightProxyMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnAcknowledgementPacket is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (h *HeightProxyMiddleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("OnTimeoutPacket is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.OnTimeoutPacket(ctx, packet, relayer)
}

func (h *HeightProxyMiddleware) NegotiateAppVersion(ctx sdk.Context, order channeltypes.Order, connectionID string, portID string, counterparty channeltypes.Counterparty, proposedVersion string) (version string, err error) {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("NegotiateAppVersion is disabled", "available", h.h, "now", ctx.BlockHeight())
		return "", ErrDisableProxyBeforeHeight
	}
	return h.internal.NegotiateAppVersion(ctx, order, connectionID, portID, counterparty, proposedVersion)
}

func (h *HeightProxyMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("SendPacket is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.SendPacket(ctx, chanCap, packet)
}

func (h *HeightProxyMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("WriteAcknowledgement is disabled", "available", h.h, "now", ctx.BlockHeight())
		return ErrDisableProxyBeforeHeight
	}
	return h.internal.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (h *HeightProxyMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	if !h.higherThan(ctx.BlockHeight()) {
		h.logger.Error("GetAppVersion is disabled", "available", h.h, "now", ctx.BlockHeight())
		return "", false
	}
	return h.internal.GetAppVersion(ctx, portID, channelID)
}

func (h *HeightProxyMiddleware) higherThan(hh int64) bool {
	if h.h == 0 {
		return false
	}
	return hh >= h.h
}

//////
type DisaleProxyMiddleware struct {
}

func NewDisaleProxyMiddleware() porttypes.Middleware {
	return &DisaleProxyMiddleware{}
}

func (d *DisaleProxyMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return "", ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version, counterpartyVersion string) (string, error) {
	return "", ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(fmt.Sprintf("OnRecvPacket is disabled"))
}

func (d *DisaleProxyMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) NegotiateAppVersion(ctx sdk.Context, order channeltypes.Order, connectionID string, portID string, counterparty channeltypes.Counterparty, proposedVersion string) (version string, err error) {
	return "", ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return ErrDisableProxyBeforeHeight
}

func (d *DisaleProxyMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return "", false
}
