package types

import (
	"bytes"
	"fmt"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/okbchain/app/types"
)

const (
	FlagTrieUseCompositeKey = "trie.use-composite-key"
	ContractStateCache      = 2048 // MB
)

func (so *stateObject) deepCopyMpt(db *CommitStateDB) *stateObject {
	acc := so.account.DeepCopy().(*types.EthAccount)
	newStateObj := newStateObject(db, acc)
	if so.trie != nil {
		newStateObj.trie = db.db.CopyTrie(so.trie)
	}

	newStateObj.code = make(types.Code, len(so.code))
	copy(newStateObj.code, so.code)
	newStateObj.dirtyStorage = so.dirtyStorage.Copy()
	newStateObj.originStorage = so.originStorage.Copy()
	newStateObj.pendingStorage = so.pendingStorage.Copy()
	newStateObj.suicided = so.suicided
	newStateObj.dirtyCode = so.dirtyCode
	newStateObj.deleted = so.deleted

	return newStateObj
}

func (so *stateObject) GetCommittedStateMpt(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if so.fakeStorage != nil {
		return so.fakeStorage[key]
	}
	// If we have a pending write or clean cached, return that
	if value, pending := so.pendingStorage[key]; pending {
		return value
	}
	if value, cached := so.originStorage[key]; cached {
		return value
	}

	var (
		enc   []byte
		value ethcmn.Hash
	)

	ctx := &so.stateDB.ctx
	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), mpt.AddressStoragePrefixMpt(so.address, so.account.StateRoot))
	enc = store.Get(key.Bytes())

	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			so.setError(err)
		}
		value.SetBytes(content)
	}

	so.originStorage[key] = value
	return value
}

func (so *stateObject) GetCommittedStateMptForQuery(db ethstate.Database, key ethcmn.Hash) []byte {
	ctx := &so.stateDB.ctx
	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), mpt.AddressStoragePrefixMpt(so.address, so.account.StateRoot))
	enc := store.Get(key.Bytes())
	return enc
}

func (so *stateObject) CodeInRawDB(db ethstate.Database) []byte {
	if so.code != nil {
		return so.code
	}
	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(so.addrHash, ethcmn.BytesToHash(so.CodeHash()))
	if err != nil {
		so.setError(fmt.Errorf("can't load code hash %x: %v", so.CodeHash(), err))
	} else {
		so.code = code
	}

	return code
}

func (so *stateObject) getTrie(db ethstate.Database) ethstate.Trie {
	if so.trie == nil {
		// Try fetching from prefetcher first
		// We don't prefetch empty tries
		if so.account.StateRoot != types2.EmptyRootHash && so.stateDB.prefetcher != nil {
			// When the miner is creating the pending state, there is no
			// prefetcher
			so.trie = so.stateDB.prefetcher.Trie(so.account.StateRoot)
		}

		if so.trie == nil {
			var err error
			so.trie, err = db.OpenStorageTrie(so.addrHash, so.account.StateRoot)
			if err != nil {
				so.setError(fmt.Errorf("failed to open storage trie: %v for addr: %s", err, so.account.EthAddress().String()))

				so.trie, _ = db.OpenStorageTrie(so.addrHash, ethcmn.Hash{})
				so.setError(fmt.Errorf("can't create storage trie: %v", err))
			}
		}
	}
	return so.trie
}

// UpdateRoot sets the trie root to the current root hash of
func (so *stateObject) updateRoot(db ethstate.Database) {
	// If nothing changed, don't bother with hashing anything
	so.updateTrie(db)
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been made
func (so *stateObject) updateTrie(db ethstate.Database) (updated bool) {
	// Make sure all dirty slots are finalized into the pending storage area
	so.finalise(false) // Don't prefetch any more, pull directly if need be
	updated = false
	if len(so.pendingStorage) == 0 {
		return
	}

	// Insert all the pending updates into the trie
	ctx := &so.stateDB.ctx
	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), mpt.AddressStoragePrefixMpt(so.address, so.account.StateRoot))
	usedStorage := make([][]byte, 0, len(so.pendingStorage))
	for key, value := range so.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == so.originStorage[key] {
			continue
		}
		updated = true
		so.originStorage[key] = value

		usedStorage = append(usedStorage, ethcmn.CopyBytes(key[:])) // Copy needed for closure
		if (value == ethcmn.Hash{}) {
			store.Delete(key[:])
			if !so.stateDB.ctx.IsCheckTx() {
				if so.stateDB.ctx.GetWatcher().Enabled() {
					so.stateDB.ctx.GetWatcher().SaveState(so.Address(), key[:], ethcmn.Hash{}.Bytes())
				}
			}
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			store.Set(key[:], v)
			if !so.stateDB.ctx.IsCheckTx() {
				if so.stateDB.ctx.GetWatcher().Enabled() {
					so.stateDB.ctx.GetWatcher().SaveState(so.Address(), key[:], v)
				}
			}
		}
	}
	if so.stateDB.prefetcher != nil {
		so.stateDB.prefetcher.Used(so.account.StateRoot, usedStorage)
	}

	if len(so.pendingStorage) > 0 {
		so.pendingStorage = make(ethstate.Storage)
	}
	return
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (so *stateObject) CommitTrie(db ethstate.Database) error {
	// If nothing changed, don't bother with hashing anything
	if updated := so.updateTrie(db); !updated {
		return nil
	}

	if so.dbErr != nil {
		return so.dbErr
	}

	return nil
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (so *stateObject) finalise(prefetch bool) {
	if so.stateDB.prefetcher != nil && prefetch && so.account.StateRoot != types2.EmptyRootHash {
		slotsToPrefetch := make([][]byte, 0, len(so.dirtyStorage))
		for key, value := range so.dirtyStorage {
			so.pendingStorage[key] = value
			if value != so.originStorage[key] {
				slotsToPrefetch = append(slotsToPrefetch, ethcmn.CopyBytes(key[:])) // Copy needed for closure
			}
		}
		if len(slotsToPrefetch) > 0 {
			so.stateDB.prefetcher.Prefetch(so.account.StateRoot, slotsToPrefetch)
		}
	} else {
		for key, value := range so.dirtyStorage {
			so.pendingStorage[key] = value
		}
	}

	if len(so.dirtyStorage) > 0 {
		so.dirtyStorage = make(ethstate.Storage)
	}
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (so *stateObject) CodeSize(db ethstate.Database) int {
	if so.code != nil {
		return len(so.code)
	}
	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(so.addrHash, ethcmn.BytesToHash(so.CodeHash()))
	if err != nil {
		so.setError(fmt.Errorf("can't load code size %x: %v", so.CodeHash(), err))
	}
	return size
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (so *stateObject) SetStorage(storage map[ethcmn.Hash]ethcmn.Hash) {
	// Allocate fake storage if it's nil.
	if so.fakeStorage == nil {
		so.fakeStorage = make(ethstate.Storage)
	}
	for key, value := range storage {
		so.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

func (so *stateObject) UpdateAccInfo() error {
	accProto := so.stateDB.accountKeeper.GetAccount(so.stateDB.ctx, so.account.Address)
	if accProto != nil {
		if ethAccount, ok := accProto.(*types.EthAccount); ok {
			so.account = ethAccount
			return nil
		}
	}
	return fmt.Errorf("fail to update account for address: %s", so.account.Address.String())
}
