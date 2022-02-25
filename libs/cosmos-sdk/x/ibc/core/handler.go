package ibc

import (
	"fmt"
	"github.com/okex/exchain/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	clienttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
)

func unmarshalFromRelayMsg(k keeper.Keeper, msg *sdk.RelayMsg) (sdk.MsgAdapter, error) {
	defer func() {
		if e := recover(); nil != e {
			fmt.Println(e)
		}
	}()
	//err := unknownproto.RejectUnknownFieldsStrict(msg.Bytes, adapter, cdc.InterfaceRegistry())
	//if err != nil {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	//}
	return common.UnmarshalGuessss(k.Codec(), msg.Bytes, new(clienttypes.MsgCreateClient),
		new(clienttypes.MsgUpdateClient),
		new(clienttypes.MsgUpgradeClient),
		new(connectiontypes.MsgConnectionOpenInit),
		new(connectiontypes.MsgConnectionOpenTry),
		new(connectiontypes.MsgConnectionOpenAck),
		new(channeltypes.MsgChannelOpenInit),
		new(channeltypes.MsgChannelOpenAck),
		new(channeltypes.MsgChannelCloseConfirm))
	//return common.UnmarshalMsgAdapter(k.Codec(), msg.Bytes)
}

// NewHandler defines the IBC handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, re sdk.Msg) (*sdk.Result, error) {
		m := re.(*sdk.RelayMsg)
		msg, err := unmarshalFromRelayMsg(k, re.(*sdk.RelayMsg))
		if nil != err {
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
