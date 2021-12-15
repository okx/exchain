package types

import (
	"fmt"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/wrap"
	"sync"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type accValue struct {
	value   exported.Account
	deleted bool
	dirty   bool
}

// Store is a wrapper for a MemDB with Commiter implementation
type AccRootKVStore struct {
	mtx    sync.Mutex
	trie   ethstate.Trie
	cache  map[string]*accValue
}

// Constructs new MemDB adapter
func NewAccRootKvStore() *AccRootKVStore {
	return &AccRootKVStore{
		cache: make(map[string]*accValue),
	}
}

func (aks *AccRootKVStore) UpdateTrie(tr ethstate.Trie) {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	aks.trie = tr
}

func (aks *AccRootKVStore) Write() {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	for key, dbValue := range aks.cache {
		if !dbValue.dirty {
			continue
		}

		addr, _ := types.AccAddressFromBech32(key)
		if dbValue.deleted {
			// delete account
			if err :=aks.trie.TryDelete(addr.Bytes()); err != nil {
				panic(err)
			}
		} else {
			data, err := rlp.EncodeToBytes(dbValue.value)
			if err != nil {
				panic(fmt.Errorf("can't encode object at %x: %v", key, err))
			}

			if err = aks.trie.TryUpdate(addr.Bytes(), data); err != nil {
				panic(err)
			}
		}
	}
}

func (aks *AccRootKVStore) Clean() {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	// Clear the cache
	aks.cache = make(map[string]*accValue)
}

func (aks *AccRootKVStore) Get(addr types.AccAddress) exported.Account {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	if cacheValue, ok := aks.cache[addr.String()]; ok {
		return cacheValue.value
	}

	enc, err := aks.trie.TryGet(addr.Bytes())
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
	aks.setCacheValue(addr.String(), wrapAcc.RealAcc, false, false)

	return wrapAcc.RealAcc
}

func (aks *AccRootKVStore) Has(addr types.AccAddress) bool {
	return aks.Get(addr) != nil
}

func (aks *AccRootKVStore) Set(addr types.AccAddress, value exported.Account) {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	aks.setCacheValue(addr.String(), value, false, true)
}

func (aks *AccRootKVStore) Delete(addr types.AccAddress) {
	aks.mtx.Lock()
	defer aks.mtx.Unlock()

	aks.setCacheValue(addr.String(), nil, true, true)
}

func (aks *AccRootKVStore) NewIterator(startKey []byte) *trie.Iterator {
	return trie.NewIterator(aks.trie.NodeIterator(startKey))
}

// Only entrypoint to mutate store.cache.
func (store *AccRootKVStore) setCacheValue(key string, value exported.Account, deleted bool, dirty bool) {
	store.cache[key] = &accValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}