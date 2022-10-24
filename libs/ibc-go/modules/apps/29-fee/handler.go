package fee

import (
	"fmt"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/keeper"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		if !tmtypes.HigherThanVenus4(ctx.BlockHeight()) {
			errMsg := fmt.Sprintf("ibc ica is not supported at height %d", ctx.BlockHeight())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		ctx.SetEventManager(sdk.NewEventManager())

		switch detailMsg := msg.(type) {
		case *types.MsgPayPacketFee:
			res, err := k.PayPacketFee(sdk.WrapSDKContext(ctx), detailMsg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRegisterCounterpartyPayee:
			res, err := k.RegisterCounterpartyPayee(sdk.WrapSDKContext(ctx), detailMsg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgPayPacketFeeAsync:
			res, err := k.PayPacketFeeAsync(sdk.WrapSDKContext(ctx), detailMsg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRegisterPayee:
			res, err := k.RegisterPayee(sdk.WrapSDKContext(ctx), detailMsg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized  fee message type: %T", msg)
		}

	}
}
