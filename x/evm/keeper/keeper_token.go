package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// OnMintVouchers After minting vouchers on this chain, convert these vouchers into evm tokens.
func (k Keeper) OnMintVouchers(ctx sdk.Context, vouchers sdk.SysCoins, receiver string) {
	cacheCtx, commit := ctx.CacheContext()
	err := k.ConvertVouchersToEvmTokens(cacheCtx, receiver, vouchers)
	if err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf("Failed to convert vouchers to evm tokens for receiver %s, coins %s. Receive error %s",
				receiver, vouchers.String(), err))
	}
	commit()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
}

func (k Keeper) ConvertVouchersToEvmTokens(ctx sdk.Context, from string, vouchers sdk.SysCoins) error {
	return nil
}

// DeleteExternalContractForDenom delete the external contract mapping for native denom,
// returns false if mapping not exists.
func (k Keeper) DeleteExternalContractForDenom(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	existingContract, found := k.getExternalContractByDenom(ctx, denom)
	if !found {
		return false
	}
	store.Delete(types.ContractToDenomKey(existingContract.Bytes()))
	store.Delete(types.DenomToExternalContractKey(denom))
	return true
}

// SetExternalContractForDenom set the external contract for native denom,
// 1. if any existing for denom, replace the old one.
// 2. if any existing for contract, return error.
func (k Keeper) SetExternalContractForDenom(ctx sdk.Context, denom string, contract common.Address) error {
	// check the contract is not registered already
	_, found := k.getDenomByContract(ctx, contract)
	if found {
		return types.ErrRegisteredContract(contract.String())
	}

	store := ctx.KVStore(k.storeKey)
	existingContract, found := k.getExternalContractByDenom(ctx, denom)
	if found {
		// delete existing mapping
		store.Delete(types.ContractToDenomKey(existingContract.Bytes()))
	}
	store.Set(types.DenomToExternalContractKey(denom), contract.Bytes())
	store.Set(types.ContractToDenomKey(contract.Bytes()), []byte(denom))
	return nil
}

func (k Keeper) getDenomByContract(ctx sdk.Context, contract common.Address) (denom string, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ContractToDenomKey(contract.Bytes()))
	if len(bz) == 0 {
		return "", false
	}
	return string(bz), true
}

// IterateMapping iterates over all the stored mapping and performs a callback function
func (k Keeper) IterateMapping(ctx sdk.Context, cb func(denom, contract string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixContractToDenom)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Value())
		conotract := common.BytesToAddress(iterator.Key()).String()

		if cb(denom, conotract) {
			break
		}
	}
}

func (k Keeper) getExternalContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DenomToExternalContractKey(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}
	return common.BytesToAddress(bz), true
}
