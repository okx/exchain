package vesting

import (
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/keeper"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/vesting/types"
)

// NewHandler returns a handler for x/auth message types.
func NewHandler(ak keeper.AccountKeeper, bk types.BankKeeper) sdk.Handler {
	msgServer := NewMsgServerImpl(ak, bk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateVestingAccount:
			res, err := msgServer.CreateVestingAccount(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
