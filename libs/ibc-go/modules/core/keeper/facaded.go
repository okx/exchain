package keeper

import (
	"context"
	"errors"

	"github.com/okex/exchain/libs/ibc-go/modules/core/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channeltyeps "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
)

var _ IBCServerKeeper = (*FacadedKeeper)(nil)

var errMisSpecificKeeper = errors.New("mis ")

type IBCServerKeeper interface {
	channeltyeps.QueryServer
	channeltyeps.MsgServer
	clienttypes.MsgServer
	connectiontypes.MsgServer

	GetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) (string, bool)
	GetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) []byte
}

// TODO, CONSTRUCTOR
type FacadedKeeper struct {
	keepers map[int64]IBCServerKeeper
}

func (f *FacadedKeeper) GetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) []byte {
	k, err := f.doGetByCtx(ctx)
	if nil != err {
		panic(types.ErrInternalConfigError)
	}
	return k.GetPacketCommitment(ctx, portID, channelID, sequence)
}

func (f *FacadedKeeper) GetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) (string, bool) {
	k, err := f.doGetByCtx(ctx)
	if nil != err {
		panic(types.ErrInternalConfigError)
	}
	return k.GetPacketReceipt(ctx, portID, channelID, sequence)
}

func (f *FacadedKeeper) Channel(goCtx context.Context, request *channeltyeps.QueryChannelRequest) (*channeltyeps.QueryChannelResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.Channel(goCtx, request)
}

func (f *FacadedKeeper) Channels(goCtx context.Context, request *channeltyeps.QueryChannelsRequest) (*channeltyeps.QueryChannelsResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.Channels(goCtx, request)
}

func (f *FacadedKeeper) ConnectionChannels(goCtx context.Context, request *channeltyeps.QueryConnectionChannelsRequest) (*channeltyeps.QueryConnectionChannelsResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ConnectionChannels(goCtx, request)
}

func (f *FacadedKeeper) ChannelClientState(goCtx context.Context, request *channeltyeps.QueryChannelClientStateRequest) (*channeltyeps.QueryChannelClientStateResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelClientState(goCtx, request)
}

func (f *FacadedKeeper) ChannelConsensusState(goCtx context.Context, request *channeltyeps.QueryChannelConsensusStateRequest) (*channeltyeps.QueryChannelConsensusStateResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelConsensusState(goCtx, request)
}

func (f *FacadedKeeper) PacketCommitment(goCtx context.Context, request *channeltyeps.QueryPacketCommitmentRequest) (*channeltyeps.QueryPacketCommitmentResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.PacketCommitment(goCtx, request)
}

func (f *FacadedKeeper) PacketCommitments(goCtx context.Context, request *channeltyeps.QueryPacketCommitmentsRequest) (*channeltyeps.QueryPacketCommitmentsResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.PacketCommitments(goCtx, request)
}

func (f *FacadedKeeper) PacketReceipt(goCtx context.Context, request *channeltyeps.QueryPacketReceiptRequest) (*channeltyeps.QueryPacketReceiptResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.PacketReceipt(goCtx, request)
}

func (f *FacadedKeeper) PacketAcknowledgement(goCtx context.Context, request *channeltyeps.QueryPacketAcknowledgementRequest) (*channeltyeps.QueryPacketAcknowledgementResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.PacketAcknowledgement(goCtx, request)
}

func (f *FacadedKeeper) PacketAcknowledgements(goCtx context.Context, request *channeltyeps.QueryPacketAcknowledgementsRequest) (*channeltyeps.QueryPacketAcknowledgementsResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.PacketAcknowledgements(goCtx, request)
}

func (f *FacadedKeeper) UnreceivedPackets(goCtx context.Context, request *channeltyeps.QueryUnreceivedPacketsRequest) (*channeltyeps.QueryUnreceivedPacketsResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.UnreceivedPackets(goCtx, request)
}

func (f *FacadedKeeper) UnreceivedAcks(goCtx context.Context, request *channeltyeps.QueryUnreceivedAcksRequest) (*channeltyeps.QueryUnreceivedAcksResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.UnreceivedAcks(goCtx, request)
}

func (f *FacadedKeeper) NextSequenceReceive(goCtx context.Context, request *channeltyeps.QueryNextSequenceReceiveRequest) (*channeltyeps.QueryNextSequenceReceiveResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.NextSequenceReceive(goCtx, request)
}

func (f *FacadedKeeper) ChannelOpenInit(goCtx context.Context, init *channeltyeps.MsgChannelOpenInit) (*channeltyeps.MsgChannelOpenInitResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelOpenInit(goCtx, init)
}

func (f *FacadedKeeper) ChannelOpenTry(goCtx context.Context, try *channeltyeps.MsgChannelOpenTry) (*channeltyeps.MsgChannelOpenTryResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelOpenTry(goCtx, try)
}

func (f *FacadedKeeper) ChannelOpenAck(goCtx context.Context, ack *channeltyeps.MsgChannelOpenAck) (*channeltyeps.MsgChannelOpenAckResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelOpenAck(goCtx, ack)
}

func (f *FacadedKeeper) ChannelOpenConfirm(goCtx context.Context, confirm *channeltyeps.MsgChannelOpenConfirm) (*channeltyeps.MsgChannelOpenConfirmResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelOpenConfirm(goCtx, confirm)
}

func (f *FacadedKeeper) ChannelCloseInit(goCtx context.Context, init *channeltyeps.MsgChannelCloseInit) (*channeltyeps.MsgChannelCloseInitResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelCloseInit(goCtx, init)
}

func (f *FacadedKeeper) ChannelCloseConfirm(goCtx context.Context, confirm *channeltyeps.MsgChannelCloseConfirm) (*channeltyeps.MsgChannelCloseConfirmResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ChannelCloseConfirm(goCtx, confirm)
}

func (f *FacadedKeeper) RecvPacket(goCtx context.Context, packet *channeltyeps.MsgRecvPacket) (*channeltyeps.MsgRecvPacketResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.RecvPacket(goCtx, packet)
}

func (f *FacadedKeeper) Timeout(goCtx context.Context, timeout *channeltyeps.MsgTimeout) (*channeltyeps.MsgTimeoutResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.Timeout(goCtx, timeout)
}

func (f *FacadedKeeper) TimeoutOnClose(goCtx context.Context, onClose *channeltyeps.MsgTimeoutOnClose) (*channeltyeps.MsgTimeoutOnCloseResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.TimeoutOnClose(goCtx, onClose)
}

func (f *FacadedKeeper) Acknowledgement(goCtx context.Context, acknowledgement *channeltyeps.MsgAcknowledgement) (*channeltyeps.MsgAcknowledgementResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.Acknowledgement(goCtx, acknowledgement)
}

func (f *FacadedKeeper) CreateClient(goCtx context.Context, client *clienttypes.MsgCreateClient) (*clienttypes.MsgCreateClientResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.CreateClient(goCtx, client)
}

func (f *FacadedKeeper) UpdateClient(goCtx context.Context, client *clienttypes.MsgUpdateClient) (*clienttypes.MsgUpdateClientResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.UpdateClient(goCtx, client)
}

func (f *FacadedKeeper) UpgradeClient(goCtx context.Context, client *clienttypes.MsgUpgradeClient) (*clienttypes.MsgUpgradeClientResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.UpgradeClient(goCtx, client)
}

func (f *FacadedKeeper) SubmitMisbehaviour(goCtx context.Context, misbehaviour *clienttypes.MsgSubmitMisbehaviour) (*clienttypes.MsgSubmitMisbehaviourResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.SubmitMisbehaviour(goCtx, misbehaviour)
}

func (f *FacadedKeeper) ConnectionOpenInit(goCtx context.Context, init *connectiontypes.MsgConnectionOpenInit) (*connectiontypes.MsgConnectionOpenInitResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ConnectionOpenInit(goCtx, init)
}

func (f *FacadedKeeper) ConnectionOpenTry(goCtx context.Context, try *connectiontypes.MsgConnectionOpenTry) (*connectiontypes.MsgConnectionOpenTryResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ConnectionOpenTry(goCtx, try)
}

func (f *FacadedKeeper) ConnectionOpenAck(goCtx context.Context, ack *connectiontypes.MsgConnectionOpenAck) (*connectiontypes.MsgConnectionOpenAckResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ConnectionOpenAck(goCtx, ack)
}

func (f *FacadedKeeper) ConnectionOpenConfirm(goCtx context.Context, confirm *connectiontypes.MsgConnectionOpenConfirm) (*connectiontypes.MsgConnectionOpenConfirmResponse, error) {
	specificK, err := f.getHeightKeeper(goCtx)
	if nil != err {
		return nil, err
	}
	return specificK.ConnectionOpenConfirm(goCtx, confirm)
}

func (f *FacadedKeeper) getHeightKeeper(goCtx context.Context) (IBCServerKeeper, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return f.doGetByCtx(ctx)
}

func (f *FacadedKeeper) doGetByCtx(ctx sdk.Context) (IBCServerKeeper, error) {
	h := ctx.BlockHeight()
	if tmtypes.HigherThanVenus3(h) {
		// veneus3 keeper
		return f.doGet(tmtypes.GetVenus3Height())
	}
	return f.doGet(tmtypes.GetVenus1Height())
}

func (f *FacadedKeeper) doGet(h int64) (IBCServerKeeper, error) {
	ret, exist := f.keepers[h]
	if !exist {
		return nil, errMisSpecificKeeper
	}
	return ret, nil
}
