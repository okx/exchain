package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm/types"
)

// GetContractDeploymentWhitelist gets the whole contract deployment whitelist currently
func (k Keeper) GetContractDeploymentWhitelist(ctx sdk.Context) (whitelist types.ContractDeploymentWhitelist) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixContractDeploymentWhitelist)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		whitelist = append(whitelist, types.SplitApprovedDeployerAddress(iterator.Key()))
	}

	return
}

// SetContractDeploymentWhitelistMember sets the deployer address as a member into whitelist
func (k Keeper) SetContractDeploymentWhitelistMember(ctx sdk.Context, distributorAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetContractDeploymentWhitelistMemberKey(distributorAddr), []byte(""))
}

// DeleteContractDeploymentWhitelistMember removes the deployer address from whitelist
func (k Keeper) DeleteContractDeploymentWhitelistMember(ctx sdk.Context, distributorAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetContractDeploymentWhitelistMemberKey(distributorAddr))
}

func (k Keeper) isDeployerInWhitelist(ctx sdk.Context, distributorAddr sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetContractDeploymentWhitelistMemberKey(distributorAddr))
}

// IsContractDeployerQualified verifies the qualification of the contract deployer
func (k Keeper) IsContractDeployerQualified(ctx sdk.Context, from sdk.AccAddress, Recipient *ethcmn.Address) bool {
	if Recipient != nil {
		// not contract creation -> pass
		return true
	}

	return k.isDeployerInWhitelist(ctx, from)
}
