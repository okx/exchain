package rootmulti

import (
	"io"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/assert"
)

func TestStore_getCached(t *testing.T) {
	type fields struct {
		db                db.DB
		flatKVDB          db.DB
		lastCommitInfo    commitInfo
		pruningOpts       types.PruningOptions
		storesParams      map[types.StoreKey]storeParams
		stores            map[types.StoreKey]types.CommitKVStore
		keysByName        map[string]types.StoreKey
		lazyLoading       bool
		pruneHeights      []int64
		versions          []int64
		traceWriter       io.Writer
		traceContext      types.TraceContext
		interBlockCache   types.MultiStorePersistentCache
		logger            log.Logger
		upgradeVersion    int64
		commitFilters     []types.StoreFilter
		pruneFilters      []types.StoreFilter
		versionFilters    []types.VersionFilter
		enableAsyncJob    bool
		jobChan           chan func()
		jobDone           *sync.WaitGroup
		metadata          atomic.Value
		commitInfoVersion sync.Map
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &Store{
				db:                tt.fields.db,
				flatKVDB:          tt.fields.flatKVDB,
				lastCommitInfo:    tt.fields.lastCommitInfo,
				pruningOpts:       tt.fields.pruningOpts,
				storesParams:      tt.fields.storesParams,
				stores:            tt.fields.stores,
				keysByName:        tt.fields.keysByName,
				lazyLoading:       tt.fields.lazyLoading,
				pruneHeights:      tt.fields.pruneHeights,
				versions:          tt.fields.versions,
				traceWriter:       tt.fields.traceWriter,
				traceContext:      tt.fields.traceContext,
				interBlockCache:   tt.fields.interBlockCache,
				logger:            tt.fields.logger,
				upgradeVersion:    tt.fields.upgradeVersion,
				commitFilters:     tt.fields.commitFilters,
				pruneFilters:      tt.fields.pruneFilters,
				versionFilters:    tt.fields.versionFilters,
				enableAsyncJob:    tt.fields.enableAsyncJob,
				jobChan:           tt.fields.jobChan,
				jobDone:           tt.fields.jobDone,
				metadata:          tt.fields.metadata,
				commitInfoVersion: tt.fields.commitInfoVersion,
			}
			assert.Equalf(t, tt.want, rs.getCached(), "getCached()")
		})
	}
}
