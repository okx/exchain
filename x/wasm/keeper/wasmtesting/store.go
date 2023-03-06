package wasmtesting

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	storetypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// MockCommitMultiStore mock with a CacheMultiStore to capture commits
type MockCommitMultiStore struct {
	sdk.CommitMultiStore
	Committed []bool
}

func (m *MockCommitMultiStore) CacheMultiStore() storetypes.CacheMultiStore {
	m.Committed = append(m.Committed, false)
	return &mockCMS{m, &m.Committed[len(m.Committed)-1]}
}

type mockCMS struct {
	sdk.CommitMultiStore
	committed *bool
}

func (m *mockCMS) GetRWSet(mp types.MsRWSet) {
	panic("implement me")
}

func (m *mockCMS) DisableCacheReadList() {
	panic("implement me")
}

func (m *mockCMS) Clear() {
	panic("implement me")
}

func (m *mockCMS) IteratorCache(isdirty bool, cb func(key string, value []byte, isDirty bool, isDelete bool, storeKey storetypes.StoreKey) bool, sKey storetypes.StoreKey) bool {
	panic("implement me")
}

func (m *mockCMS) Write() {
	*m.committed = true
}
