package transient

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"

	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
)

var _ types.Committer = (*Store)(nil)
var _ types.KVStore = (*Store)(nil)

// Store is a wrapper for a MemDB with Commiter implementation
type Store struct {
	dbadapter.Store
}

// Constructs new MemDB adapter
func NewStore() *Store {
	return &Store{Store: dbadapter.Store{DB: dbm.NewMemDB()}}
}

// Implements CommitStore
// Commit cleans up Store.
func (ts *Store) Commit() (id types.CommitID) {
	ts.Store = dbadapter.Store{DB: dbm.NewMemDB()}
	return
}

// Implements CommitStore
func (ts *Store) SetPruning(pruning types.PruningOptions) {
}

// Implements CommitStore
func (ts *Store) LastCommitID() (id types.CommitID) {
	return
}

// Implements Store.
func (ts *Store) GetStoreType() types.StoreType {
	return types.StoreTypeTransient
}

func (ts *Store) GetDBWriteCount() int {
	return 0
}

func (ts *Store) GetDBReadTime() int {
	return 0
}

func (ts *Store) GetDBReadCount() int {
	return 0
}
func (ts *Store) GetNodeReadCount() int {
	return 0
}

func (ts *Store) ResetCount() {
}
