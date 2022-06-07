package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Commit(err error) {
	if !Enable() {
		return
	}
	if err == nil {
		for _, msg := range txStateCache {
			blockStateCache[string(msg.key)] = msg
		}
	}
	txStateCache = txStateCache[:0]
}

func Flush() {
	if !Enable() {
		return
	}
	blockStateCacheCopy := blockStateCache
	blockStateCache = make(map[string]*watcherMessage)
	task := func() {
		for _, msg := range blockStateCacheCopy {
			if msg.isDelete {
				_ = db.Delete(msg.key)
			} else {
				_ = db.Set(msg.key, msg.value)
			}
		}
	}
	tasks <- task

}

type writeKVStore struct {
	sdk.KVStore
}

func WrapWriteKVStore(store sdk.KVStore) sdk.KVStore {
	once.Do(func() {
		initDB()
	})

	if !Enable() {
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
