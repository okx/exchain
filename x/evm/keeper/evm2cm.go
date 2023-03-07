package keeper

import (
	ethermint "github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/evm/types"
)

// SetSysContractAddress set system contract address to store
func (k *Keeper) SetSysContractAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	store := k.paramSpace.CustomKVStore(ctx)
	key := types.GetSysContractAddressKey()
	store.Set(key, addr)
	return nil
}

// DelSysContractAddress del system contract address to store
func (k *Keeper) DelSysContractAddress(ctx sdk.Context) sdk.Error {
	store := k.paramSpace.CustomKVStore(ctx)
	key := types.GetSysContractAddressKey()
	store.Delete(key)
	return nil
}

func (k *Keeper) GetSysContractAddress(ctx sdk.Context) (sdk.AccAddress, sdk.Error) {
	store := k.paramSpace.CustomKVStore(ctx)
	key := types.GetSysContractAddressKey()
	value := store.Get(key)
	if value == nil {
		return nil, types.ErrSysContractAddressIsNotExist(types.ErrKeyNotFound)
	}
	return value, nil
}

func (k *Keeper) IsMatchSysContractAddress(ctx sdk.Context, addr sdk.AccAddress) bool {
	iaddr, err := k.GetSysContractAddress(ctx)
	if err != nil {
		return false
	}
	return iaddr.Equals(addr)
}

func (k Keeper) IsContractAccount(ctx sdk.Context, addr sdk.AccAddress) bool {
	acct := k.accountKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return false
	}
	ethAcct, ok := acct.(*ethermint.EthAccount)
	if !ok {
		return false
	}
	return ethAcct.IsContract()
}

func querySysContractAddress(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res, err := keeper.GetSysContractAddress(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
