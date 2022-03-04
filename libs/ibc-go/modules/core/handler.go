package ibc

import (
	"github.com/okex/exchain/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
)

var (
	MsgDetailUpdateClient = "MsgDetailUpdateClient"

	MsgConnectionOpenInit    = "MsgConnectionOpenInit"
	MsgConnectionOpenTry     = "MsgConnectionOpenTry"
	MsgConnectionOpenConfirm = "MsgConnectionOpenConfirm"
	MsgConnectionOpenAck     = "MsgConnectionOpenAck"

	MsgChannelOpenInit    = "MsgChannelOpenInit"
	MsgChannelOpenTry     = "MsgChannelOpenTry"
	MsgChannelOpenAck     = "MsgChannelOpenAck"
	MsgChannelOpenConfirm = "MsgChannelOpenConfirm"
)

func unmarshalFromRelayMsg(k keeper.Keeper, msg *sdk.RelayMsg) (sdk.MsgAdapter, error) {
	defer func() {
		if e := recover(); nil != e {
			panic(e)
		}
	}()
	//err := unknownproto.RejectUnknownFieldsStrict(msg.Bytes, adapter, cdc.InterfaceRegistry())
	//if err != nil {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	//}
	ms := make([]sdk.MsgProtoAdapter, 0)

	switch msg.MsgType {
	case MsgDetailUpdateClient:
		ms = append(ms, new(clienttypes.MsgUpdateClient))
	case MsgConnectionOpenTry:
		ms = append(ms, new(connectiontypes.MsgConnectionOpenTry))
	case MsgConnectionOpenConfirm:
		ms = append(ms, new(connectiontypes.MsgConnectionOpenConfirm))
	case MsgConnectionOpenInit:
		ms = append(ms, new(connectiontypes.MsgConnectionOpenInit))

	case MsgChannelOpenInit:
		ms = append(ms, new(channeltypes.MsgChannelOpenInit))
	case MsgChannelOpenTry:
		ms = append(ms, new(channeltypes.MsgChannelOpenTry))
	case MsgConnectionOpenAck:
		ms = append(ms, new(connectiontypes.MsgConnectionOpenAck))
	case MsgChannelOpenAck:
		ms = append(ms, new(channeltypes.MsgChannelOpenAck))
	case MsgChannelOpenConfirm:
		ms = append(ms, new(channeltypes.MsgChannelOpenConfirm))

	default:
		ms = append(ms, new(clienttypes.MsgCreateClient),
			new(channeltypes.MsgChannelCloseConfirm),
			new(clienttypes.MsgUpgradeClient),
			new(channeltypes.MsgChannelOpenAck),
		)
	}
	return common.UnmarshalGuessss(k.Codec(), msg.Bytes, ms...,
	)
	//return common.UnmarshalMsgAdapter(k.Codec(), msg.Bytes)
}

// NewHandler defines the IBC handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, re sdk.Msg) (*sdk.Result, error) {
		m := re.(*sdk.RelayMsg)
		msg, err := unmarshalFromRelayMsg(k, re.(*sdk.RelayMsg))
		if nil != err {
			panic(err)
			aaa := new(types.ClientState)
			err := k.Codec().GetProtocMarshal().UnmarshalBinaryBare(m.Bytes, aaa)
			err = k.Codec().GetProtocMarshal().UnmarshalInterface(m.Bytes, &aaa)
			k.Codec().GetProtocMarshal().UnmarshalBinaryLengthPrefixed(m.Bytes, aaa)
			err = aaa.Unmarshal(m.Bytes)
			if nil != err {
				return nil, err
			}
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// IBC client msg interface types
		case *clienttypes.MsgCreateClient:
			res, err := k.CreateClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgUpdateClient:
			res, err := k.UpdateClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgUpgradeClient:
			res, err := k.UpgradeClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgSubmitMisbehaviour:
			res, err := k.SubmitMisbehaviour(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// IBC connection msgs
		case *connectiontypes.MsgConnectionOpenInit:
			res, err := k.ConnectionOpenInit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenTry:
			res, err := k.ConnectionOpenTry(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenAck:
			res, err := k.ConnectionOpenAck(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenConfirm:
			res, err := k.ConnectionOpenConfirm(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// IBC channel msgs
		case *channeltypes.MsgChannelOpenInit:
			res, err := k.ChannelOpenInit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenTry:
			res, err := k.ChannelOpenTry(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenAck:
			res, err := k.ChannelOpenAck(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenConfirm:
			res, err := k.ChannelOpenConfirm(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelCloseInit:
			res, err := k.ChannelCloseInit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelCloseConfirm:
			res, err := k.ChannelCloseConfirm(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// IBC packet msgs get routed to the appropriate module callback
		case *channeltypes.MsgRecvPacket:
			res, err := k.RecvPacket(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgAcknowledgement:
			res, err := k.Acknowledgement(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgTimeout:
			res, err := k.Timeout(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgTimeoutOnClose:
			res, err := k.TimeoutOnClose(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
		}
	}
}
