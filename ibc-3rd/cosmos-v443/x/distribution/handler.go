package distribution

import (
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/distribution/keeper"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/distribution/types"
	govtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSetWithdrawAddress:
			res, err := msgServer.SetWithdrawAddress(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdrawDelegatorReward:
			res, err := msgServer.WithdrawDelegatorReward(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdrawValidatorCommission:
			res, err := msgServer.WithdrawValidatorCommission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgFundCommunityPool:
			res, err := msgServer.FundCommunityPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

func NewCommunityPoolSpendProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CommunityPoolSpendProposal:
			return keeper.HandleCommunityPoolSpendProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
