package keeper

import (
	"context"

	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	internaltypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/typesadapter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServer = (*AccountKeeper)(nil)
)

func (ak AccountKeeper) Accounts(c context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(ak.key)
	accountsStore := prefix.NewStore(store, types.AddressStoreKeyPrefix)

	var accounts []*codectypes.Any
	pageRes, err := query.Paginate(accountsStore, req.Pagination, func(key, value []byte) error {
		account := ak.decodeAccount(value)
		ba := convEthAccountToBaseAccount(account)
		any, err := codectypes.NewAnyWithValue(ba)
		if err != nil {
			return err
		}

		accounts = append(accounts, any)
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	return &types.QueryAccountsResponse{Accounts: accounts, Pagination: pageRes}, err

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

// Params returns parameters of auth module
func (ak AccountKeeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := ak.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func convEthAccountToBaseAccount(account exported.Account) *internaltypes.BaseAccount {
	ba := &internaltypes.BaseAccount{
		Address:       account.GetAddress().String(),
		PubKey:        nil,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}
	return ba
}
