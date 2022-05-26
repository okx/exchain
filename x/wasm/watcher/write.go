package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Flush() {
	if !enableWatcher {
		return
	}
	for k, v := range wasmStateCache {
		if v != nil {
			_ = db.Set([]byte(k), v)
		} else {
			_ = db.Delete([]byte(k))
		}
	}
	wasmStateCache = make(map[string][]byte)
}

type writeKVStore struct {
	sdk.KVStore
}

func WrapWriteKVStore(store sdk.KVStore) sdk.KVStore {
	once.Do(func() {
		initDB()
	})

	if !enableWatcher {
		return store
	}

	return &writeKVStore{
		KVStore: store,
	}
}

func (w *writeKVStore) Set(key, value []byte) {
	w.KVStore.Set(key, value)
	wasmStateCache[string(key)] = value
}

func (w *writeKVStore) Delete(key []byte) {
	w.KVStore.Delete(key)
	wasmStateCache[string(key)] = nil
}
