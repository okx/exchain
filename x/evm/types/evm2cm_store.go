package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// SetSysContractAddress set system contract address to store
func (csdb *CommitStateDB) SetSysContractAddress(addr sdk.AccAddress) sdk.Error {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetSysContractAddressKey()
	store.Set(key, addr)
	return nil
}

// DelSysContractAddress del system contract address to store
func (csdb *CommitStateDB) DelSysContractAddress() sdk.Error {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetSysContractAddressKey()
	store.Delete(key)
	return nil
}

func (csdb *CommitStateDB) GetSysContractAddress() (sdk.AccAddress, sdk.Error) {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetSysContractAddressKey()
	value := store.Get(key)
	if value == nil {
		return nil, ErrSysContractAddressIsNotExist(ErrKeyNotFound)
	}
	return value, nil
}
