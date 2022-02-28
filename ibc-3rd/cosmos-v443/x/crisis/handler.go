package crisis

import (
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/crisis/types"
)

// RouterKey
const RouterKey = types.ModuleName

func NewHandler(k types.MsgServer) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgVerifyInvariant:
			res, err := k.VerifyInvariant(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized crisis message type: %T", msg)
		}
	}
}
