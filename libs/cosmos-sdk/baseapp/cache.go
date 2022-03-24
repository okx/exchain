package baseapp

import "sync"

type blockDataCache struct {
	senderLock sync.RWMutex
	sender     map[string]string
}

func NewBlockDataCache() *blockDataCache {
	return &blockDataCache{
		sender: make(map[string]string),
	}
}

func (cache *blockDataCache) SetSender(tx []byte, sender string) {
	cache.senderLock.Lock()
	cache.sender[string(tx)] = sender
	cache.senderLock.Unlock()
}

func (cache *blockDataCache) GetSender(tx []byte) (sender string, ok bool) {
	cache.senderLock.RLock()
	sender, ok = cache.sender[string(tx)]
	cache.senderLock.RUnlock()
	return
}

func (cache *blockDataCache) Clear() {
	cache.senderLock.Lock()
	for k := range cache.sender {
		delete(cache.sender, k)
	}
	cache.senderLock.Unlock()
}
