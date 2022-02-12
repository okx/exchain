package transfer

import (
	"fmt"
	"github.com/okex/exchain/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/types"
)

// NewHandler returns sdk.Handler for IBC token transfer module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, re sdk.Msg) (*sdk.Result, error) {
		msg,err:=common.UnmarshalMsgAdapter(k.Codec(),re.(*sdk.RelayMsg).Bytes)
		if nil != err {
			return nil, err
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgTransferAdapter:
			res, err := k.Transfer(sdk.WrapSDKContext(ctx), msg.ToMsgTransfer())
			fmt.Println(res)
			return sdk.WrapServiceResult(ctx, nil, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 transfer message type: %T", msg)
		}
	}
}
