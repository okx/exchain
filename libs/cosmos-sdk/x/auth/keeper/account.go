package keeper

import (
	"context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	types2 "github.com/okex/exchain/temp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewAccountWithAddress implements sdk.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		panic(err)
	}
	return ak.NewAccount(ctx, acc)
}

// NewAccount sets the next account number to a given account interface
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc exported.Account) exported.Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// GetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	if data, gas, ok := ctx.Cache().GetAccount(ethcmn.BytesToAddress(addr)); ok {
		ctx.GasMeter().ConsumeGas(gas, "x/auth/keeper/account.go/GetAccount")
		if data == nil {
			return nil
		}

		return data.Copy().(exported.Account)
	}

	store := ctx.KVStore(ak.key)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		ctx.Cache().UpdateAccount(addr, nil, len(bz), false)
		return nil
	}
	acc := ak.decodeAccount(bz)
	ctx.Cache().UpdateAccount(addr, acc, len(bz), false)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) (accounts []exported.Account) {
	ak.IterateAccounts(ctx,
		func(acc exported.Account) (stop bool) {
			accounts = append(accounts, acc)
			return false
		})
	return accounts
}

// SetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBareWithRegisteredMarshaller(acc)
	if err != nil {
		bz, err = ak.cdc.MarshalBinaryBare(acc)
	}
	if err != nil {
		panic(err)
	}
	store.Set(types.AddressStoreKey(addr), bz)
	ctx.Cache().UpdateAccount(acc.GetAddress(), acc, len(bz), true)

	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		if ak.observers != nil {
			for _, observer := range ak.observers {
				if observer != nil {
					observer.OnAccountUpdated(acc)
				}
			}
		}
	}
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(types.AddressStoreKey(addr))
	ctx.Cache().UpdateAccount(addr, nil, 0, true)
}

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, cb func(account exported.Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iterator := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		account := ak.decodeAccount(iterator.Value())

		if cb(account) {
			break
		}
	}
}

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
	req.Address = "ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9"
	addr, err := sdk.AccAddressFromBech32(req.Address)

	if err != nil {
		return nil, err
	}

	account := ak.GetAccount(ctx, addr)
	if account == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", req.Address)
	}
	//ethA:=account.(*ethermint.EthAccount)
	ba := &types2.BaseAccount{
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

//func(ak AccountKeeper)getProtobufAccount(ctx sdk.Context, addr sdk.AccAddress) exported.AccountAdapter{
//if data, gas, ok := ctx.Cache().GetAccount(ethcmn.BytesToAddress(addr)); ok {
//		ctx.GasMeter().ConsumeGas(gas, "x/auth/keeper/account.go/GetAccount")
//		if data == nil {
//			return nil
//		}
//
//		return data.Copy().(exported.AccountAdapter)
//	}
//
//	store := ctx.KVStore(ak.key)
//	bz := store.Get(types.AddressStoreKey(addr))
//	if bz == nil {
//		ctx.Cache().UpdateAccount(addr, nil, len(bz), false)
//		return nil
//	}
//	acc := ak.decodeAccount(bz)
//	ctx.Cache().UpdateAccount(addr, acc, len(bz), false)
//
//	return acc
//}
