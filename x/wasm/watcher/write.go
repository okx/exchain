package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Reset() {
	if !enableWatcher {
		return
	}
	txStateCache = txStateCache[:0]

}

func Commit() {
	if !enableWatcher {
		return
	}
	for _, msg := range txStateCache {
		blockStateCache[string(msg.key)] = msg
	}
	txStateCache = txStateCache[:0]
}

func Flush() {
	if !enableWatcher {
		return
	}
	for key, msg := range blockStateCache {
		if msg.isDelete {
			_ = db.Delete(msg.key)
		} else {
			_ = db.Set(msg.key, msg.value)
		}
		delete(blockStateCache, key)
	}
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
	txStateCache = append(txStateCache, &watcherMessage{key: key, value: value})
}

func (w *writeKVStore) Delete(key []byte) {
	w.KVStore.Delete(key)
	txStateCache = append(txStateCache, &watcherMessage{key: key, isDelete: true})
}
