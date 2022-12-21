package keeper

import (
	"context"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	icatypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	"github.com/okex/exchain/x/icamauth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterchainAccount implements the Query/InterchainAccount gRPC method
func (k Keeper) InterchainAccount(goCtx context.Context, req *types.QueryInterchainAccountRequest) (*types.QueryInterchainAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not find account: %s", err)
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no account found for portID %s", portID)
	}

	return types.NewQueryInterchainAccountResponse(addr), nil
}
