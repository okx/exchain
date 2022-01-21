package rootmulti

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/tendermint/tm-db"
	"io"
)

var commithash = []byte("FAKE_HASH")

var _ types.Committer = (*CommitDBStoreAdapter)(nil)
var _ types.KVStore = (*CommitDBStoreAdapter)(nil)

//----------------------------------------
// CommitDBStoreWrapper should only be used for simulation/debugging,
// as it doesn't compute any commit hash, and it cannot load older state.

// Wrapper type for dbm.Db with implementation of KVStore
type CommitDBStoreAdapter struct {
	dbm.DB
	KvCache *fastcache.Cache
}

// Constructs new MemDB adapter
func NewCommitDBStore(db dbm.DB) *CommitDBStoreAdapter {
	return &CommitDBStoreAdapter{
		DB: db,
		KvCache: fastcache.New(2*1024*1024*1024),
	}
}

func (cdsa *CommitDBStoreAdapter) Commit(*iavl.TreeDelta, []byte) (types.CommitID, iavl.TreeDelta, []byte) {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}, iavl.TreeDelta{}, nil
}

func (cdsa *CommitDBStoreAdapter) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa *CommitDBStoreAdapter) SetPruning(_ types.PruningOptions) {}

func (cdsa *CommitDBStoreAdapter) GetDBReadTime() int   { return 0 }
func (cdsa *CommitDBStoreAdapter) GetDBWriteCount() int { return 0 }

func (cdsa *CommitDBStoreAdapter) GetDBReadCount() int   { return 0 }
func (cdsa *CommitDBStoreAdapter) GetNodeReadCount() int { return 0 }

func (cdsa *CommitDBStoreAdapter) ResetCount() {}

// Get wraps the underlying DB's Get method panicing on error.
func (cdsa *CommitDBStoreAdapter) Get(key []byte) []byte {
	if enc := cdsa.KvCache.Get(nil, key); len(enc) > 0 {
		return enc
	}
	v, err := cdsa.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return v
}

// Has wraps the underlying DB's Has method panicing on error.
func (cdsa *CommitDBStoreAdapter) Has(key []byte) bool {
	if cdsa.KvCache.Has(key) {
		return true
	}
	ok, err := cdsa.DB.Has(key)
	if err != nil {
		panic(err)
	}

	return ok
}

// Set wraps the underlying DB's Set method panicing on error.
func (cdsa *CommitDBStoreAdapter) Set(key, value []byte) {
	if err := cdsa.DB.Set(key, value); err != nil {
		panic(err)
	}
	cdsa.KvCache.Set(key, value)
}

// Delete wraps the underlying DB's Delete method panicing on error.
func (cdsa *CommitDBStoreAdapter) Delete(key []byte) {
	if err := cdsa.DB.Delete(key); err != nil {
		panic(err)
	}
	cdsa.KvCache.Del(key)
}

// Iterator wraps the underlying DB's Iterator method panicing on error.
func (cdsa *CommitDBStoreAdapter) Iterator(start, end []byte) types.Iterator {
	iter, err := cdsa.DB.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	return iter
}

// ReverseIterator wraps the underlying DB's ReverseIterator method panicing on error.
func (cdsa *CommitDBStoreAdapter) ReverseIterator(start, end []byte) types.Iterator {
	iter, err := cdsa.DB.ReverseIterator(start, end)
	if err != nil {
		panic(err)
	}

	return iter
}

// GetStoreType returns the type of the store.
func (cdsa *CommitDBStoreAdapter) GetStoreType() types.StoreType {
	return types.StoreTypeDB
}

// CacheWrap cache wraps the underlying store.
func (cdsa *CommitDBStoreAdapter) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(cdsa)
}

// CacheWrapWithTrace implements KVStore.
func (cdsa *CommitDBStoreAdapter) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(cdsa, w, tc))
}
