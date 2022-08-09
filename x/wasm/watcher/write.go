package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Save(err error) {
	if !Enable() {
		return
	}
	txCacheMtx.Lock()
	if err == nil {
		for _, msg := range txStateCache {
			blockStateCache[string(msg.Key)] = msg
		}
	}
	txStateCache = txStateCache[:0]
	txCacheMtx.Unlock()
}

func Commit() {
	if !Enable() {
		return
	}
	blockStateCacheCopy := blockStateCache
	blockStateCache = make(map[string]*WatchMessage)
	task := func() {
		batch := db.NewBatch()
		for _, msg := range blockStateCacheCopy {
			if msg.IsDelete {
				batch.Delete(msg.Key)
			} else {
				batch.Set(msg.Key, msg.Value)
			}
		}
		if err := batch.Write(); err != nil {
			logger.Error("wasm watchDB", "batch write error", err)
		}
	}
	tasks <- task
}

var tasks = make(chan func(), 5*3)

func taskRoutine() {
	for task := range tasks {
		task()
	}
}

type writeKVStore struct {
	sdk.KVStore
}

func WrapWriteKVStore(store sdk.KVStore) sdk.KVStore {
	if !Enable() {
		return store
	}

	return &writeKVStore{
		KVStore: store,
	}
}

func (w *writeKVStore) Set(key, value []byte) {
	w.KVStore.Set(key, value)
	txStateCache = append(txStateCache, &WatchMessage{Key: key, Value: value})
}

func (w *writeKVStore) Delete(key []byte) {
	w.KVStore.Delete(key)
	txStateCache = append(txStateCache, &WatchMessage{Key: key, IsDelete: true})
}
