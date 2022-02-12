package mem

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/listenkv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"io"

	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
)

var (
	_ types.KVStore   = (*Store)(nil)
	_ types.Committer = (*Store)(nil)
)

// Store implements an in-memory only KVStore. Entries are persisted between
// commits and thus between blocks. State in Memory store is not committed as part of app state but maintained privately by each node
type Store struct {
	dbadapter.Store
}

func (s *Store) CommitterCommit(_ *iavl.TreeDelta) (_ types.CommitID, _ *iavl.TreeDelta) {
	return
}

func (s *Store) GetDBReadTime() int {
	return 0
}

func (s *Store) GetDBWriteCount() int {
	return 0
}

func (s *Store) GetDBReadCount() int {
	return 0
}

func (s *Store) GetNodeReadCount() int {
	return 0
}

func (s *Store) GetFlatKVReadTime() int {
	return 0
}

func (s *Store) GetFlatKVWriteTime() int {
	return 0
}

func (s *Store) GetFlatKVReadCount() int {
	return 0
}

func (s *Store) GetFlatKVWriteCount() int {
	return 0
}

func (s *Store) ResetCount() {}

func NewStore() *Store {
	return NewStoreWithDB(dbm.NewMemDB())
}

func NewStoreWithDB(db *dbm.MemDB) *Store { // nolint: interfacer
	return &Store{Store: dbadapter.Store{DB: db}}
}

// GetStoreType returns the Store's type.
func (s Store) GetStoreType() types.StoreType {
	return types.StoreTypeMemory
}

// CacheWrap branches the underlying store.
func (s Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(s)
}

// CacheWrapWithTrace implements KVStore.
func (s Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(s, w, tc))
}

// CacheWrapWithListeners implements the CacheWrapper interface.
func (s Store) CacheWrapWithListeners(storeKey types.StoreKey, listeners []types.WriteListener) types.CacheWrap {
	return cachekv.NewStore(listenkv.NewStore(s, storeKey, listeners))
}

// Commit performs a no-op as entries are persistent between commitments.
func (s *Store) Commit() (id types.CommitID) { return }

func (s *Store) SetPruning(pruning types.PruningOptions) {}

// GetPruning is a no-op as pruning options cannot be directly set on this store.
// They must be set on the root commit multi-store.
func (s *Store) GetPruning() types.PruningOptions { return types.PruningOptions{} }

func (s Store) LastCommitID() (id types.CommitID) { return }
