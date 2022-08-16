package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	outtypes "github.com/okex/exchain/x/distribution/types"
)

// ValidatorCommission queries accumulated commission for a validator
func (k Keeper) ValidatorCommission(c context.Context, req *outtypes.QueryValidatorCommissionRequest) (*outtypes.QueryValidatorCommissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	valAdr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	commission := k.GetValidatorAccumulatedCommission(ctx, valAdr)

	return &outtypes.QueryValidatorCommissionResponse{Commission: commission}, nil
}
