package keeper

import (
	"context"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	internaltypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/internal"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServer = (*AccountKeeper)(nil)
)

func (ak AccountKeeper) Accounts(ctx context.Context, request *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	return nil, nil
}

func (ak AccountKeeper) Account(conte context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "Address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(conte)
	addr, err := sdk.AccAddressFromBech32(req.Address)

	if err != nil {
		return nil, err
	}

	account := ak.GetAccount(ctx, addr)
	if account == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", req.Address)
	}
	//ethA:=account.(*ethermint.EthAccount)
	ba := &internaltypes.BaseAccount{
		Address:       account.GetAddress().String(),
		PubKey:        nil,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}
	any, err := codectypes.NewAnyWithValue(ba)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &types.QueryAccountResponse{Account: any}, nil
}
