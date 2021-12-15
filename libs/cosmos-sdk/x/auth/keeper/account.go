package keeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/wrap"
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
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) (acc exported.Account) {
	//store := ctx.KVStore(ak.key)
	//bz := store.Get(types.AddressStoreKey(addr))
	//if bz == nil {
	//	return nil
	//}
	//acc := ak.decodeAccount(bz)
	//return acc

	store := types.NewGasKvStore(ctx.GetAccCacheStore(), types2.KVGasConfig(), ctx.GasMeter())
	return store.Get(addr)
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
	//addr := acc.GetAddress()
	//store := ctx.KVStore(ak.key)
	//bz, err := ak.cdc.MarshalBinaryBare(acc)
	//if err != nil {
	//	panic(err)
	//}
	//store.Set(types.AddressStoreKey(addr), bz)

	store := types.NewGasKvStore(ctx.GetAccCacheStore(), types2.KVGasConfig(), ctx.GasMeter())
	store.Set(acc)

	if ak.observers != nil && !ctx.IsCheckTx() {
		for _, observer := range ak.observers {
			if observer != nil {
				observer.OnAccountUpdated(acc)
			}
		}
	}
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc exported.Account) {
	//addr := acc.GetAddress()
	//store := ctx.KVStore(ak.key)
	//store.Delete(types.AddressStoreKey(addr))

	store := types.NewGasKvStore(ctx.GetAccCacheStore(), types2.KVGasConfig(), ctx.GasMeter())
	store.Delete(acc)
}

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, cb func(account exported.Account) (stop bool)) {
	//store := ctx.KVStore(ak.key)
	//iterator := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)
	//
	//defer iterator.Close()
	//for ; iterator.Valid(); iterator.Next() {
	//	account := ak.decodeAccount(iterator.Value())
	//
	//	if cb(account) {
	//		break
	//	}
	//}

	store := types.NewGasKvStore(ctx.GetAccCacheStore(), types2.KVGasConfig(), ctx.GasMeter())
	itr := store.NewIterator(nil)
	for itr.Next() {
		val := itr.Value()
		var wrapAcc wrap.WrapAccount
		if err := rlp.DecodeBytes(val, &wrapAcc); err != nil {
			continue
		}

		if cb(wrapAcc.RealAcc) {
			break
		}
	}
}

func (ak *AccountKeeper) NewCacheStore(ctx sdk.Context) sdk.AccCacheStore {
	if ctx.IsCheckTx() {
		return types.NewCacheStore(ak.checkRootStore)
	} else {
		return types.NewCacheStore(ak.deliverRootStore)
	}
}

func (ak *AccountKeeper) PushData2Database(ctx sdk.Context, root ethcmn.Hash) {
	triedb := ak.db.TrieDB()
	// Full but not archive node, do proper garbage collection
	triedb.Reference(root, ethcmn.Hash{}) // metadata reference to keep trie alive
	ak.triegc.Push(root, -int64(ctx.BlockHeight()))

	if types.TrieDirtyDisabled {
		if err := triedb.Commit(root, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
		ak.SetLatestStoredBlockHeight(uint64(ctx.BlockHeight()))
	} else {
		if ctx.BlockHeight() > core.TriesInMemory {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = ethcmn.StorageSize(256) * 1024 * 1024
			)

			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			chosen := ctx.BlockHeight() - core.TriesInMemory

			// If the header is missing (canonical chain behind), we're reorging a low
			// diff sidechain. Suspend committing until this operation is completed.
			chRoot := ak.GetRootMptHash(uint64(chosen))
			if chRoot == (ethcmn.Hash{}) {
				ak.Logger(ctx).Debug("Reorg in progress, trie commit postponed", "number", chosen)
			} else {
				ak.SetLatestStoredBlockHeight(uint64(chosen))

				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				if err := triedb.Commit(chRoot, true, nil); err != nil {
					panic("fail to commit mpt data: " + err.Error())
				}
			}

			// Garbage collect anything below our required write retention
			for !ak.triegc.Empty() {
				root, number := ak.triegc.Pop()
				if int64(-number) > chosen {
					ak.triegc.Push(root, number)
					break
				}
				triedb.Dereference(root.(ethcmn.Hash))
			}
		}
	}

	ak.checkRootStore.Clean()
	ak.deliverRootStore.Clean()
}
