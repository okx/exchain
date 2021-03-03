package keeper

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/x/evm/types"
)

// ----------------------------------------------------------------------------
// Setters, only for test use
// ----------------------------------------------------------------------------

// SetBalance calls CommitStateDB.SetBalance using the passed in context
func (k *Keeper) SetBalance(ctx sdk.Context, addr ethcmn.Address, amount *big.Int) {
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	csdb.SetBalance(addr, amount)
	_ = csdb.Finalise(false)
}

// SetNonce calls CommitStateDB.SetNonce using the passed in context
func (k *Keeper) SetNonce(ctx sdk.Context, addr ethcmn.Address, nonce uint64) {
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	csdb.SetNonce(addr, nonce)
	_ = csdb.Finalise(false)
}

func (k *Keeper) SetLogs(ctx sdk.Context, hash ethcmn.Hash, logs []*ethtypes.Log) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLogs)
	bz, err := types.MarshalLogs(logs)
	if err != nil {
		return err
	}

	store.Set(hash.Bytes(), bz)
	return nil
}

// ----------------------------------------------------------------------------
// Getters, for test and query case
// ----------------------------------------------------------------------------

// GetBalance calls CommitStateDB.GetBalance using the passed in context
func (k *Keeper) GetBalance(ctx sdk.Context, addr ethcmn.Address) *big.Int {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).GetBalance(addr)
}

// GetCode calls CommitStateDB.GetCode using the passed in context
func (k *Keeper) GetCode(ctx sdk.Context, addr ethcmn.Address) []byte {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).GetCode(addr)
}

// GetState calls CommitStateDB.GetState using the passed in context
func (k *Keeper) GetState(ctx sdk.Context, addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).GetState(addr, hash)
}

// GetLogs calls CommitStateDB.GetLogs using the passed in context
func (k *Keeper) GetLogs(ctx sdk.Context, hash ethcmn.Hash) ([]*ethtypes.Log, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLogs)
	bz := store.Get(hash.Bytes())
	if len(bz) == 0 {
		// return nil error if logs are not found
		return []*ethtypes.Log{}, nil
	}

	return types.UnmarshalLogs(bz)
}

// AllLogs calls CommitStateDB.AllLogs using the passed in context
func (k *Keeper) AllLogs(ctx sdk.Context) []*ethtypes.Log {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixLogs)
	defer iterator.Close()

	allLogs := []*ethtypes.Log{}
	for ; iterator.Valid(); iterator.Next() {
		var logs []*ethtypes.Log
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &logs)
		allLogs = append(allLogs, logs...)
	}

	return allLogs
}

// ----------------------------------------------------------------------------
// Auxiliary, for test and query case
// ----------------------------------------------------------------------------

// ForEachStorage calls CommitStateDB.ForEachStorage using passed in context
func (k *Keeper) ForEachStorage(ctx sdk.Context, addr ethcmn.Address, cb func(key, value ethcmn.Hash) bool) error {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).ForEachStorage(addr, cb)
}

// IterateStorage calls CommitStateDB.ForEachStorage using passed in context
func (k *Keeper) IterateStorage(ctx sdk.Context, addr ethcmn.Address, cb func(key, value []byte) bool) error {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).IterateStorage(addr, cb)
}

// GetOrNewStateObject calls CommitStateDB.GetOrNetStateObject using the passed in context
func (k *Keeper) GetOrNewStateObject(ctx sdk.Context, addr ethcmn.Address) types.StateObject {
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).GetOrNewStateObject(addr)
}
