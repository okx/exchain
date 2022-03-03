package types

import (
	"errors"
	"fmt"
	ethermint "github.com/okex/exchain/app/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types2 "github.com/okex/exchain/libs/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/mpt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types2 "github.com/okex/exchain/libs/types"
)

func (csdb *CommitStateDB) CommitMpt(deleteEmptyObjects bool) (ethcmn.Hash, error) {
	// Finalize any pending changes and merge everything into the tries
	csdb.IntermediateRoot(deleteEmptyObjects)

	// Commit objects to the trie, measuring the elapsed time
	codeWriter := csdb.db.TrieDB().DiskDB().NewBatch()
	for addr := range csdb.stateObjectsDirty {
		if obj := csdb.stateObjects[addr]; !obj.deleted {
			// Write any contract code associated with the state object
			if obj.code != nil && obj.dirtyCode {
				rawdb.WriteCode(codeWriter, ethcmn.BytesToHash(obj.CodeHash()), obj.code)
				obj.dirtyCode = false
			}

			// Write any storage changes in the state object to its storage trie
			if err := obj.CommitTrie(csdb.db); err != nil {
				return ethcmn.Hash{}, err
			}

			if tmtypes.HigherThanMars(csdb.ctx.BlockHeight()) || types2.EnableDoubleWrite {
				accProto := csdb.accountKeeper.GetAccount(csdb.ctx, obj.account.Address)
				if ethermintAccount, ok := accProto.(*ethermint.EthAccount); ok {
					ethermintAccount.StateRoot = obj.account.StateRoot
					csdb.accountKeeper.SetAccount(csdb.ctx, ethermintAccount)
				}
			}
		}
	}

	if len(csdb.stateObjectsDirty) > 0 {
		csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	}

	if codeWriter.ValueSize() > 0 {
		if err := codeWriter.Write(); err != nil {
			csdb.SetError(fmt.Errorf("failed to commit dirty codes: %s", err.Error()))
		}
	}

	return ethcmn.Hash{}, nil
}

func (csdb *CommitStateDB) ForEachStorageMpt(so *stateObject, cb func(key, value ethcmn.Hash) (stop bool)) error {
	it := trie.NewIterator(so.getTrie(csdb.db).NodeIterator(nil))
	for it.Next() {
		key := ethcmn.BytesToHash(so.trie.GetKey(it.Key))
		if value, dirty := so.dirtyStorage[key]; dirty {
			if cb(key, value) {
				return nil
			}
			continue
		}

		if len(it.Value) > 0 {
			_, content, _, err := rlp.Split(it.Value)
			if err != nil {
				return err
			}
			if cb(key, ethcmn.BytesToHash(content)) {
				return nil
			}
		}
	}

	return nil
}

func (csdb *CommitStateDB) GetStateByKeyMpt(addr ethcmn.Address, key ethcmn.Hash) ethcmn.Hash {
	var (
		enc []byte
		err error
	)

	if enc, err = csdb.StorageTrie(addr).TryGet(key.Bytes()); err != nil {
		return ethcmn.Hash{}
	}

	var value ethcmn.Hash
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			return ethcmn.Hash{}
		}
		value.SetBytes(content)
	}

	return value
}

func (csdb *CommitStateDB) GetCodeByHashMpt(hash ethcmn.Hash) []byte {
	code, err := csdb.db.ContractCode(ethcmn.Hash{}, hash)
	if err != nil {
		return nil
	}

	return code
}

// getDeletedStateObject is similar to getStateObject, but instead of returning
// nil for a deleted state object, it returns the actual object with the deleted
// flag set. This is needed by the state journal to revert to the correct s-
// destructed object instead of wiping all knowledge about the state object.
func (csdb *CommitStateDB) getDeletedStateObject(addr ethcmn.Address) *stateObject {
	// Prefer live objects if any is available
	if obj := csdb.stateObjects[addr]; obj != nil {
		if _, ok := csdb.updatedAccount[addr]; ok {
			delete(csdb.updatedAccount, addr)
			if err := obj.UpdateAccInfo(); err != nil {
				csdb.SetError(err)
				return nil
			}
		}
		return obj
	}

	// otherwise, attempt to fetch the account from the account mapper
	acc := csdb.accountKeeper.GetAccount(csdb.ctx, sdk.AccAddress(addr.Bytes()))
	if acc == nil {
		csdb.SetError(fmt.Errorf("no account found for address: %s", addr.String()))
		return nil
	}

	// insert the state object into the live set
	so := newStateObject(csdb, acc)
	csdb.setStateObject(so)

	return so
}


func (csdb *CommitStateDB) MarkUpdatedAcc(addList []ethcmn.Address) {
	for _, addr := range addList {
		csdb.updatedAccount[addr] = struct{}{}
	}
}

// ----------------------------------------------------------------------------
// Proof related
// ----------------------------------------------------------------------------
//// GetProof returns the Merkle proof for a given account.
//func (csdb *CommitStateDB) GetProof(addr ethcmn.Address) ([][]byte, error) {
//	return csdb.GetProofByHash(crypto.Keccak256Hash(addr.Bytes()))
//}
//
//// GetProofByHash returns the Merkle proof for a given account.
//func (csdb *CommitStateDB) GetProofByHash(addrHash ethcmn.Hash) ([][]byte, error) {
//	var proof mpt.ProofList
//	err := csdb.trie.Prove(addrHash[:], 0, &proof)
//	return proof, err
//}

// GetStorageProof returns the Merkle proof for given storage slot.
func (csdb *CommitStateDB) GetStorageProof(a ethcmn.Address, key ethcmn.Hash) ([][]byte, error) {
	var proof mpt.ProofList
	addrTrie := csdb.StorageTrie(a)
	if addrTrie == nil {
		return proof, errors.New("storage trie for requested address does not exist")
	}
	err := addrTrie.Prove(crypto.Keccak256(key.Bytes()), 0, &proof)
	return proof, err
}

func (csdb *CommitStateDB) Logger() log.Logger {
	return csdb.ctx.Logger().With("module", ModuleName)
}
