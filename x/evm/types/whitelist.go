package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// SetContractDeploymentWhitelistMember sets the target address list into whitelist store
func (csdb *CommitStateDB) SetContractDeploymentWhitelist(addrList AddressList) {
	if csdb.Watcher.Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.Watcher.SaveContractDeploymentWhitelistItem(addrList[i])
		}
	}
	store := csdb.ctx.KVStore(csdb.storeKey)
	for i := 0; i < len(addrList); i++ {
		store.Set(GetContractDeploymentWhitelistMemberKey(addrList[i]), []byte(""))
	}
}

// DeleteContractDeploymentWhitelist deletes the target address list from whitelist store
func (csdb *CommitStateDB) DeleteContractDeploymentWhitelist(addrList AddressList) {
	if csdb.Watcher.Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.Watcher.DeleteContractDeploymentWhitelist(addrList[i])
		}
	}
	store := csdb.ctx.KVStore(csdb.storeKey)
	for i := 0; i < len(addrList); i++ {
		store.Delete(GetContractDeploymentWhitelistMemberKey(addrList[i]))
	}
}

// GetContractDeploymentWhitelist gets the whole contract deployment whitelist currently
func (csdb *CommitStateDB) GetContractDeploymentWhitelist() (whitelist AddressList) {
	store := csdb.ctx.KVStore(csdb.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyPrefixContractDeploymentWhitelist)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		whitelist = append(whitelist, splitApprovedDeployerAddress(iterator.Key()))
	}

	return
}

// IsDeployerInWhitelist checks whether the deployer is in the whitelist as a distributor
func (csdb *CommitStateDB) IsDeployerInWhitelist(deployerAddr sdk.AccAddress) bool {
	bs := csdb.dbAdapter.NewStore(csdb.ctx.KVStore(csdb.storeKey), KeyPrefixContractDeploymentWhitelist)
	return bs.Has(deployerAddr)
}

// SetContractBlockedList sets the target address list into blocked list store
func (csdb *CommitStateDB) SetContractBlockedList(addrList AddressList) {
	if csdb.Watcher.Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.Watcher.SaveContractBlockedListItem(addrList[i])
		}
	}
	store := csdb.ctx.KVStore(csdb.storeKey)
	for i := 0; i < len(addrList); i++ {
		store.Set(GetContractBlockedListMemberKey(addrList[i]), []byte(""))
	}
}

// DeleteContractBlockedList deletes the target address list from blocked list store
func (csdb *CommitStateDB) DeleteContractBlockedList(addrList AddressList) {
	if csdb.Watcher.Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.Watcher.DeleteContractBlockedList(addrList[i])
		}
	}
	store := csdb.ctx.KVStore(csdb.storeKey)
	for i := 0; i < len(addrList); i++ {
		store.Delete(GetContractBlockedListMemberKey(addrList[i]))
	}
}

// GetContractBlockedList gets the whole contract blocked list currently
func (csdb *CommitStateDB) GetContractBlockedList() (blockedList AddressList) {
	store := csdb.ctx.KVStore(csdb.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyPrefixContractBlockedList)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		blockedList = append(blockedList, splitBlockedContractAddress(iterator.Key()))
	}

	return
}

// IsContractInBlockedList checks whether the contract address is in the blocked list
func (csdb *CommitStateDB) IsContractInBlockedList(contractAddr sdk.AccAddress) bool {
	bs := csdb.dbAdapter.NewStore(csdb.ctx.KVStore(csdb.storeKey), KeyPrefixContractBlockedList)
	return bs.Has(contractAddr)
}
