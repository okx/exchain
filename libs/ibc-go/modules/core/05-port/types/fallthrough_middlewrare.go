package types

import (
	"github.com/okex/exchain/libs/ibc-go/modules/core/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

var (
	_ Middleware = (*FallThroughMiddleware)(nil)
)

type FallThroughMiddleware struct {
	*common.SelectorStrategy
}

func NewFallThroughMiddleware(defaultMiddleware Middleware, factories ...common.SelectorFactory) Middleware {
	ret := FallThroughMiddleware{}
	ret.SelectorStrategy = common.NewSelectorStrategy(defaultMiddleware)
	ret.SelectorStrategy.RegisterSelectors(factories...)
	ret.SelectorStrategy.Seal()

	return &ret
}

func (f *FallThroughMiddleware) getProxy(ctx sdk.Context) Middleware {
	return f.SelectorStrategy.GetProxy(ctx).(Middleware)
}

func (f *FallThroughMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return f.getProxy(ctx).OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (f *FallThroughMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version, counterpartyVersion string) (string, error) {
	return f.getProxy(ctx).OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
}

func (f *FallThroughMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return f.getProxy(ctx).OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (f *FallThroughMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return f.getProxy(ctx).OnChanOpenConfirm(ctx, portID, channelID)
}

func (f *FallThroughMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return f.getProxy(ctx).OnChanCloseInit(ctx, portID, channelID)
}

func (f *FallThroughMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return f.getProxy(ctx).OnChanCloseConfirm(ctx, portID, channelID)
}

func (f *FallThroughMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	return f.getProxy(ctx).OnRecvPacket(ctx, packet, relayer)
}

func (f *FallThroughMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return f.getProxy(ctx).OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (f *FallThroughMiddleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return f.getProxy(ctx).OnTimeoutPacket(ctx, packet, relayer)
}

func (f *FallThroughMiddleware) NegotiateAppVersion(ctx sdk.Context, order channeltypes.Order, connectionID string, portID string, counterparty channeltypes.Counterparty, proposedVersion string) (version string, err error) {
	return f.getProxy(ctx).NegotiateAppVersion(ctx, order, connectionID, portID, counterparty, proposedVersion)
}

func (f *FallThroughMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	return f.getProxy(ctx).SendPacket(ctx, chanCap, packet)
}

func (f *FallThroughMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return f.getProxy(ctx).WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (f *FallThroughMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return f.getProxy(ctx).GetAppVersion(ctx, portID, channelID)
}
