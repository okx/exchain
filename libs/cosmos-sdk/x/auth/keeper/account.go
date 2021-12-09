package keeper

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
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

	defer func(){
		ctx.GasMeter().ConsumeGas(ak.gsConfig.ReadCostFlat, types2.GasReadCostFlatDesc)
		ctx.GasMeter().ConsumeGas(ak.gsConfig.ReadCostPerByte*types2.Gas(estimateAccByteLenForGasConsume(acc)), types2.GasReadPerByteDesc)
	}()

	if ctx.IsCheckTx() {
		if val := ak.checkTxStore.Get(addr.String()); val != nil {
			return val
		}
	} else {
		if val := ak.deliverTxStore.Get(addr.String()); val != nil {
			return val
		}

		if val, ok := ak.accLRU.Get(addr.String()); ok {
			return val.(exported.Account)
		}
	}

	enc, err := ak.trie.TryGet(addr.Bytes())
	if err != nil {
		return nil
	}
	if len(enc) == 0 {
		return nil
	}

	var wrapAcc wrap.WrapAccount
	err = rlp.DecodeBytes(enc, &wrapAcc)
	if err != nil {
		return nil
	}

	return wrapAcc.RealAcc
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

	defer func(){
		ctx.GasMeter().ConsumeGas(ak.gsConfig.WriteCostFlat, types2.GasWriteCostFlatDesc)
		ctx.GasMeter().ConsumeGas(ak.gsConfig.WriteCostPerByte*types2.Gas(estimateAccByteLenForGasConsume(acc)), types2.GasWritePerByteDesc)
	}()

	if ak.observers != nil && !ctx.IsCheckTx() {
		for _, observer := range ak.observers {
			if observer != nil {
				observer.OnAccountUpdated(acc)
			}
		}
	}

	if ctx.IsCheckTx() {
		ak.checkTxStore.Set(acc.GetAddress().String(), acc)
	} else {
		ak.deliverTxStore.Set(acc.GetAddress().String(), acc)
	}
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc exported.Account) {
	//addr := acc.GetAddress()
	//store := ctx.KVStore(ak.key)
	//store.Delete(types.AddressStoreKey(addr))

	defer func(){
		ctx.GasMeter().ConsumeGas(ak.gsConfig.DeleteCost, types2.GasDeleteDesc)
	}()

	if ctx.IsCheckTx() {
		ak.checkTxStore.Delete(acc.GetAddress().String())
	} else {
		ak.deliverTxStore.Delete(acc.GetAddress().String())
	}
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

	it := trie.NewIterator(ak.trie.NodeIterator(nil))
	for it.Next() {
		ctx.GasMeter().ConsumeGas(ak.gsConfig.IterNextCostFlat, types2.GasIterNextCostFlatDesc)
		if len(it.Value) > 0 {
			var wrapAcc wrap.WrapAccount
			if err := rlp.DecodeBytes(it.Value, &wrapAcc); err != nil {
				continue
			}

			ctx.GasMeter().ConsumeGas(ak.gsConfig.ReadCostPerByte*types2.Gas(estimateAccByteLenForGasConsume(wrapAcc.RealAcc)), types2.GasValuePerByteDesc)
			if cb(wrapAcc.RealAcc) {
				break
			}
		}
	}
}

func (ak *AccountKeeper) Update(ctx sdk.Context, err error) {
	if !ctx.IsCheckTx() && err == nil {
		ak.deliverTxStore.IteratorCache(func(key string, acc exported.Account, isDirty bool, isDelete bool) {
			if !isDirty {
				return
			}

			accKey,_ := sdk.AccAddressFromBech32(key)
			if isDelete {
				ak.accLRU.Remove(key)

				// delete account
				ak.trie.TryDelete(accKey)

			} else {
				ak.accLRU.Add(key, acc)

				data, err := rlp.EncodeToBytes(acc)
				if err != nil {
					panic(fmt.Errorf("can't encode object at %x: %v", key, err))
				}

				if err = ak.trie.TryUpdate(accKey, data); err != nil {
					panic(err)
				}
			}
		})
	}

	ak.CleanCacheStore()
}

func (ak *AccountKeeper) CleanCacheStore() {
	ak.checkTxStore.Clean()
	ak.deliverTxStore.Clean()
}

func (ak *AccountKeeper) PushData2Database(ctx sdk.Context, root ethcmn.Hash) {
	triedb := ak.db.TrieDB()
	// Full but not archive node, do proper garbage collection
	triedb.Reference(root, ethcmn.Hash{}) // metadata reference to keep trie alive
	ak.triegc.Push(root, -int64(ctx.BlockHeight()))

	if types.TrieDirtyDisabled {
		triedb.Commit(root, false, nil)
		ak.SetLatestBlockHeight(uint64(ctx.BlockHeight()))
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
				ak.SetLatestBlockHeight(uint64(chosen))

				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				triedb.Commit(chRoot, true, nil)
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
}

func estimateAccByteLenForGasConsume(acc exported.Account) int64{
	if acc == nil {
		return 0
	}

	if acc.IsEthAccount() {
		return 150
	}

	return 70
}
