package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (csdb *CommitStateDB) ProtectStateDBEnvironment(ctx sdk.Context) {
	subCtx, commit := ctx.CacheContextWithMultiSnapShotRWSet()
	currentGasMeter := subCtx.GasMeter()
	infGasMeter := sdk.GetReusableInfiniteGasMeter()
	subCtx.SetGasMeter(infGasMeter)

	//push dirty object to ctx
	for addr := range csdb.journal.dirties {
		obj, exist := csdb.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		if obj.suicided || obj.empty() {
			csdb.deleteStateObjectForProtect(subCtx, obj)
		} else {
			obj.finaliseForProtect() // Prefetch slots in the background
			obj.commitStateForProtect(subCtx)
			csdb.updateStateObjectForProtect(subCtx, obj)
		}
	}

	//clear state objects and add revert handle
	for addr, preObj := range csdb.stateObjects {
		delete(csdb.stateObjects, addr)
		//when need to revertsnapshot need resetObject to csdb
		csdb.journal.append(resetObjectChange{prev: preObj})
	}
	//commit data to parent ctx
	///when need to revertsnapshot need restore ctx to prev
	csdb.CMChangeCommit(commit)

	subCtx.SetGasMeter(currentGasMeter)
	sdk.ReturnInfiniteGasMeter(infGasMeter)
}

func (csdb *CommitStateDB) CMChangeCommit(writeCacheWithRWSet func() types.MultiSnapShotWSet) {
	cmwSet := writeCacheWithRWSet()
	csdb.journal.append(cmChange{&cmwSet})
}

// updateStateObject writes the given state object to the store.
func (csdb *CommitStateDB) updateStateObjectForProtect(ctx sdk.Context, so *stateObject) error {
	// NOTE: we don't use sdk.NewCoin here to avoid panic on test importer's genesis
	newBalance := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDecFromBigIntWithPrec(so.Balance(), sdk.Precision)} // int2dec
	if !newBalance.IsValid() {
		return fmt.Errorf("invalid balance %s", newBalance)
	}

	//checking and reject tx if address in blacklist
	if csdb.bankKeeper.BlacklistedAddr(so.account.GetAddress()) {
		return fmt.Errorf("address <%s> in blacklist is not allowed", so.account.GetAddress().String())
	}

	coins := so.account.GetCoins()
	balance := coins.AmountOf(newBalance.Denom)
	if balance.IsZero() || !balance.Equal(newBalance.Amount) {
		coins = coins.Add(newBalance)
	}

	if err := so.account.SetCoins(coins); err != nil {
		return err
	}

	csdb.accountKeeper.SetAccount(ctx, so.account)
	if !ctx.IsCheckTx() {
		if ctx.GetWatcher().Enabled() {
			ctx.GetWatcher().SaveAccount(so.account)
		}
	}

	return nil
}

// deleteStateObject removes the given state object from the state store.
func (csdb *CommitStateDB) deleteStateObjectForProtect(ctx sdk.Context, so *stateObject) {
	so.deleted = true
	csdb.accountKeeper.RemoveAccount(ctx, so.account)
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (so *stateObject) finaliseForProtect() {
	for key, value := range so.dirtyStorage {
		so.pendingStorage[key] = value
	}
	if len(so.dirtyStorage) > 0 {
		so.dirtyStorage = make(ethstate.Storage)
	}
}

// commitState commits all dirty storage to a KVStore and resets
// the dirty storage slice to the empty state.
func (so *stateObject) commitStateForProtect(ctx sdk.Context) {
	// Make sure all dirty slots are finalized into the pending storage area
	so.finaliseForProtect() // Don't prefetch any more, pull directly if need be
	if len(so.pendingStorage) == 0 {
		return
	}

	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), AddressStoragePrefix(so.Address()))
	for key, value := range so.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == so.originStorage[key] {
			continue
		}
		so.originStorage[key] = value

		prefixKey := GetStorageByAddressKey(so.Address().Bytes(), key.Bytes())
		if (value == ethcmn.Hash{}) {
			store.Delete(prefixKey.Bytes())
			ctx.Cache().UpdateStorage(so.address, prefixKey, value.Bytes(), true)
			if !ctx.IsCheckTx() {
				if ctx.GetWatcher().Enabled() {
					ctx.GetWatcher().SaveState(so.Address(), prefixKey.Bytes(), ethcmn.Hash{}.Bytes())
				}
			}
		} else {
			store.Set(prefixKey.Bytes(), value.Bytes())
			ctx.Cache().UpdateStorage(so.address, prefixKey, value.Bytes(), true)
			if !ctx.IsCheckTx() {
				if ctx.GetWatcher().Enabled() {
					ctx.GetWatcher().SaveState(so.Address(), prefixKey.Bytes(), value.Bytes())
				}
			}
		}
	}

	if len(so.pendingStorage) > 0 {
		so.pendingStorage = make(ethstate.Storage)
	}

	return
}
