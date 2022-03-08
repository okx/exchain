package types

import (
	"bytes"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

const (
	FlagContractStateCache = "contract-state-cache"
	FlagUseCompositeKey string = "use-composite-key"
)

var (
	ContractStateCache uint = 2048 // MB
	UseCompositeKey = true
)

func (so *stateObject) deepCopyMpt(db *CommitStateDB) *stateObject {
	acc := db.accountKeeper.NewAccountWithAddress(db.ctx, so.account.Address)
	newStateObj := newStateObject(db, acc)
	if so.trie != nil {
		newStateObj.trie = db.db.CopyTrie(so.trie)
	}

	newStateObj.code = so.code
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
		err   error
		value ethcmn.Hash
	)

	prefixKey := AssembleCompositeKey(so.address.Bytes(), key.Bytes())
	if enc = so.stateDB.StateCache.Get(nil, prefixKey.Bytes()); len(enc) > 0 {
		value.SetBytes(enc)
	} else {
		tmpKey := key
		if UseCompositeKey {
			tmpKey = so.GetStorageByAddressKey(key.Bytes())
		}

		if enc, err = so.getTrie(db).TryGet(tmpKey.Bytes()); err != nil {
			so.setError(err)
			return ethcmn.Hash{}
		}

		if len(enc) > 0 {
			_, content, _, err := rlp.Split(enc)
			if err != nil {
				so.setError(err)
			}
			value.SetBytes(content)
		}
	}

	so.originStorage[key] = value
	return value
}

func (so *stateObject) CodeMpt(db ethstate.Database) []byte {
	if so.code != nil {
		return so.code
	}
	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(so.addrHash, ethcmn.BytesToHash(so.CodeHash()))
	if err != nil {
		so.setError(fmt.Errorf("can't load code hash %x: %v", so.CodeHash(), err))
	}
	so.code = code

	return code
}

func (so *stateObject) getTrie(db ethstate.Database) ethstate.Trie {
	if so.trie == nil {
		var err error
		so.trie, err = db.OpenStorageTrie(so.addrHash, so.account.StateRoot)
		if err != nil {
			so.setError(fmt.Errorf("failed to open storage trie: %v for addr: %s", err, so.account.EthAddress().String()))

			so.trie, _ = db.OpenStorageTrie(so.addrHash, ethcmn.Hash{})
			so.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return so.trie
}

// UpdateRoot sets the trie root to the current root hash of
func (so *stateObject) updateRoot(db ethstate.Database) {
	// If nothing changed, don't bother with hashing anything
	if so.updateTrie(db) == nil {
		return
	}
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been made
func (so *stateObject) updateTrie(db ethstate.Database) ethstate.Trie {
	// Make sure all dirty slots are finalized into the pending storage area
	so.finalise(false) // Don't prefetch any more, pull directly if need be
	if len(so.pendingStorage) == 0 {
		return so.trie
	}

	// Insert all the pending updates into the trie
	tr := so.getTrie(db)
	for key, value := range so.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == so.originStorage[key] {
			continue
		}
		so.originStorage[key] = value

		prefixKey := AssembleCompositeKey(so.address.Bytes(), key.Bytes())
		if UseCompositeKey {
			key = so.GetStorageByAddressKey(key.Bytes())
		}
		if (value == ethcmn.Hash{}) {
			so.setError(tr.TryDelete(key[:]))
			so.stateDB.StateCache.Del(prefixKey.Bytes())
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			so.setError(tr.TryUpdate(key[:], v))
			so.stateDB.StateCache.Set(prefixKey.Bytes(), value.Bytes())
		}
	}

	if len(so.pendingStorage) > 0 {
		so.pendingStorage = make(ethstate.Storage)
	}
	return tr
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (so *stateObject) CommitTrie(db ethstate.Database) error {
	// If nothing changed, don't bother with hashing anything
	if so.updateTrie(db) == nil {
		return nil
	}
	if so.dbErr != nil {
		return so.dbErr
	}

	root, err := so.trie.Commit(nil)
	if err == nil {
		so.account.StateRoot = root
	}
	return err
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (so *stateObject) finalise(prefetch bool) {
	for key, value := range so.dirtyStorage {
		so.pendingStorage[key] = value
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
	if !tmtypes.HigherThanMars(so.stateDB.ctx.BlockHeight()) {
		return len(so.Code(db))
	} else {
		if bytes.Equal(so.CodeHash(), emptyCodeHash) {
			return 0
		}
		size, err := db.ContractCodeSize(so.addrHash, ethcmn.BytesToHash(so.CodeHash()))
		if err != nil {
			so.setError(fmt.Errorf("can't load code size %x: %v", so.CodeHash(), err))
		}
		return size
	}
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

func AssembleCompositeKey(prefix, key []byte) ethcmn.Hash {
	compositeKey := make([]byte, (len(prefix)+len(key))/2)

	copy(compositeKey, prefix[:len(prefix)/2])
	copy(compositeKey[len(prefix)/2:], key[len(key)/2:])
	return ethcmn.BytesToHash(compositeKey)
}
