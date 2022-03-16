package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/erc20/types"
	"github.com/okex/exchain/x/params"
)

// Keeper wraps the CommitStateDB, allowing us to pass in SDK context while adhering
// to the StateDB interface.
type Keeper struct {
	cdc            *codec.Codec
	storeKey       sdk.StoreKey
	paramSpace     Subspace
	accountKeeper  AccountKeeper
	supplyKeeper   SupplyKeeper
	bankKeeper     BankKeeper
	govKeeper      GovKeeper
	evmKeeper      EvmKeeper
	transferKeeper TransferKeeper
}

// NewKeeper generates new erc20 module keeper
func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace,
	ak AccountKeeper, sk SupplyKeeper, bk BankKeeper,
	ek EvmKeeper, tk TransferKeeper) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     paramSpace,
		accountKeeper:  ak,
		supplyKeeper:   sk,
		bankKeeper:     bk,
		evmKeeper:      ek,
		transferKeeper: tk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetGovKeeper sets keeper of gov
func (k *Keeper) SetGovKeeper(gk GovKeeper) {
	k.govKeeper = gk
}

// SetExternalContractForDenom set the external contract for native denom,
// 1. if any existing for denom, replace the old one.
// 2. if any existing for contract, return error.
func (k Keeper) SetExternalContractForDenom(ctx sdk.Context, denom string, contract common.Address) error {
	// check the contract is not registered already
	_, found := k.GetDenomByContract(ctx, contract)
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

// GetExternalContracts returns all external contract mappings
func (k Keeper) GetExternalContracts(ctx sdk.Context) (out []types.TokenMapping) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixDenomToExternalContract)

	for ; iter.Valid(); iter.Next() {
		out = append(out, types.TokenMapping{
			Denom:    string(iter.Key()),
			Contract: common.BytesToAddress(iter.Value()).Hex(),
		})
	}
	return
}

// getExternalContractByDenom find the corresponding external contract for the denom
func (k Keeper) getExternalContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DenomToExternalContractKey(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}
	return common.BytesToAddress(bz), true
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

// SetAutoContractForDenom set the auto deployed contract for native denom
func (k Keeper) SetAutoContractForDenom(ctx sdk.Context, denom string, contract common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DenomToAutoContractKey(denom), contract.Bytes())
	store.Set(types.ContractToDenomKey(contract.Bytes()), []byte(denom))
}

// GetAutoContracts returns all auto-deployed contract mappings
func (k Keeper) GetAutoContracts(ctx sdk.Context) (out []types.TokenMapping) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixDenoToAutoContract)
	for ; iter.Valid(); iter.Next() {
		out = append(out, types.TokenMapping{
			Denom:    string(iter.Key()),
			Contract: common.BytesToAddress(iter.Value()).Hex(),
		})
	}
	return
}

// getAutoContractByDenom find the corresponding auto-deployed contract for the denom
func (k Keeper) getAutoContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DenomToAutoContractKey(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}
	return common.BytesToAddress(bz), true
}

// GetDenomByContract find native denom by contract address
func (k Keeper) GetDenomByContract(ctx sdk.Context, contract common.Address) (denom string, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ContractToDenomKey(contract.Bytes()))
	if len(bz) == 0 {
		return "", false
	}
	return string(bz), true
}

// GetContractByDenom find the corresponding contract for the denom,
// external contract is taken in preference to auto-deployed one
func (k Keeper) GetContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	contract, found = k.getExternalContractByDenom(ctx, denom)
	if !found {
		contract, found = k.getAutoContractByDenom(ctx, denom)
	}
	return
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
