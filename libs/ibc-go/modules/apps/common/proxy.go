package common

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

var (
	_ porttypes.Middleware = (*DisaleProxyMiddleware)(nil)
)

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
