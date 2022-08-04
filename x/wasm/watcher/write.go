package watcher

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Commit(err error) {
	if !Enable() {
		return
	}
	txCacheMtx.Lock()
	if err == nil {
		for _, msg := range txStateCache {
			blockStateCache[string(msg.key)] = msg
		}
	}
	txStateCache = txStateCache[:0]
	txCacheMtx.Unlock()
}

func Flush() {
	if !Enable() {
		return
	}
	blockStateCacheCopy := blockStateCache
	blockStateCache = make(map[string]*watcherMessage)
	task := func() {
		batch := db.NewBatch()
		for _, msg := range blockStateCacheCopy {
			if msg.isDelete {
				batch.Delete(msg.key)
			} else {
				batch.Set(msg.key, msg.value)
			}
		}
		if err := batch.Write(); err != nil {
			fmt.Println("batch write error", err)
		}
	}
	tasks <- task

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
	txStateCache = append(txStateCache, &watcherMessage{key: key, value: value})
}

func (w *writeKVStore) Delete(key []byte) {
	w.KVStore.Delete(key)
	txStateCache = append(txStateCache, &watcherMessage{key: key, isDelete: true})
}
