package transient

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	"github.com/okx/okbchain/libs/iavl"
	dbm "github.com/okx/okbchain/libs/tm-db"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/dbadapter"
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
func (ts *Store) Commit(*iavl.TreeDelta, []byte) (id types.CommitID, _ iavl.TreeDelta, _ []byte) {
	ts.Store = dbadapter.Store{DB: dbm.NewMemDB()}
	return
}

func (ts *Store) CommitterCommit(interface{}) (id types.CommitID, _ interface{}) {
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

func (ts *Store) LastCommitVersion() (v int64) {
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

func (ts *Store) GetFlatKVReadTime() int {
	return 0
}

func (ts *Store) GetFlatKVWriteTime() int {
	return 0
}

func (ts *Store) GetFlatKVReadCount() int {
	return 0
}

func (ts *Store) GetFlatKVWriteCount() int {
	return 0
}

func (ts *Store) SetUpgradeVersion(int64) {

}
