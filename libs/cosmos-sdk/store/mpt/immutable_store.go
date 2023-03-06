package mpt

import (
	"io"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
)

type ImmutableMptStore struct {
	trie ethstate.Trie
	db   ethstate.Database
	root ethcmn.Hash
	mtx  sync.Mutex
}

func NewImmutableMptStore(db ethstate.Database, root ethcmn.Hash) (*ImmutableMptStore, error) {
	ms := &ImmutableMptStore{
		db:   db,
		root: root,
	}
	trie, err := ms.db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	ms.trie = trie
	return ms, nil
}

func NewImmutableMptStoreFromTrie(db ethstate.Database, trie ethstate.Trie) *ImmutableMptStore {
	ms := &ImmutableMptStore{
		trie: db.CopyTrie(trie),
		db:   db,
		root: trie.Hash(),
	}
	return ms
}

func (ms *ImmutableMptStore) Get(key []byte) []byte {
	ms.mtx.Lock()
	defer ms.mtx.Unlock()

	value, err := ms.trie.TryGet(key)
	if err != nil {
		return nil
	}
	return value
}

func (ms *ImmutableMptStore) Has(key []byte) bool {
	return ms.Get(key) != nil
}

func (ms *ImmutableMptStore) Set(key []byte, value []byte) {
	panic("immutable store cannot set")
}

func (ms *ImmutableMptStore) Delete(key []byte) {
	panic("immutable store cannot delete")
}

func (ms *ImmutableMptStore) Iterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.db.CopyTrie(ms.trie), start, end)
}

func (ms *ImmutableMptStore) ReverseIterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.db.CopyTrie(ms.trie), start, end)
}

func (ms *ImmutableMptStore) GetStoreType() types.StoreType {
	return StoreTypeMPT
}

func (ms *ImmutableMptStore) CacheWrap() types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(ms)
}

func (ms *ImmutableMptStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(tracekv.NewStore(ms, w, tc))
}

var _ types.KVStore = (*ImmutableMptStore)(nil)
