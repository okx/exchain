package watcher

import (
	"container/list"
	"encoding/hex"
	"sync"
)

type MessageCache struct {
	mtx sync.RWMutex
	mp  map[string]WatchMessage // if the key of value WatchMessage is nil, this key should del on db batch write
}

func newMessageCache() *MessageCache {
	return &MessageCache{
		mp: make(map[string]WatchMessage),
	}
}

func (c *MessageCache) Set(wsg WatchMessage) {
	if wsg == nil {
		return
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.mp[hex.EncodeToString(wsg.GetKey())] = wsg
}

func (c *MessageCache) BatchDel(keys [][]byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, k := range keys {
		c.mp[hex.EncodeToString(k)] = &Batch{Key: k, TypeValue: TypeDelete}
	}
}

func (c *MessageCache) BatchSet(wsgs []WatchMessage) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, wsg := range wsgs {
		if wsg == nil {
			continue
		}
		c.mp[hex.EncodeToString(wsg.GetKey())] = wsg
	}
}

func (c *MessageCache) BatchSetEx(batchs []*Batch) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, b := range batchs {
		if b == nil {
			continue
		}
		c.mp[hex.EncodeToString(b.GetKey())] = b
	}
}

func (c *MessageCache) Get(key []byte) (WatchMessage, bool) {
	if len(key) == 0 {
		return nil, false
	}
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if v, ok := c.mp[hex.EncodeToString(key)]; ok {
		return v, true
	}
	return nil, false
}

func (c *MessageCache) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for k := range c.mp {
		delete(c.mp, k)
	}
}

type MessageCacheEvent struct {
	*MessageCache
	version int64
}

type commitCache struct {
	mtx sync.RWMutex
	m   map[int64]*list.Element
	l   *list.List // in the value is *MessageCacheEvent
}

func newCommitCache() *commitCache {
	return &commitCache{
		m: make(map[int64]*list.Element),
		l: list.New(),
	}
}

func (cc *commitCache) pushBack(version int64, ca *MessageCacheEvent) {
	cc.mtx.Lock()
	defer cc.mtx.Unlock()
	if elm, ok := cc.m[version]; ok {
		elm.Value = ca
		return
	}
	elm := cc.l.PushBack(ca)
	cc.m[version] = elm
}

func (cc *commitCache) remove(version int64) *MessageCacheEvent {
	cc.mtx.Lock()
	defer cc.mtx.Unlock()
	if elm, ok := cc.m[version]; ok {
		value := cc.l.Remove(elm)
		delete(cc.m, version)
		return value.(*MessageCacheEvent)
	}
	return nil
}

func (cc *commitCache) getTop() (*MessageCacheEvent, bool) {
	cc.mtx.Lock()
	defer cc.mtx.Unlock()
	elm := cc.l.Front()
	if elm == nil {
		return nil, false
	}
	return elm.Value.(*MessageCacheEvent), true
}

func (cc *commitCache) getElementFromCache(key []byte) (WatchMessage, bool) {
	cc.mtx.RLock()
	defer cc.mtx.RUnlock()
	for e := cc.l.Back(); e != nil; e = e.Prev() {
		if v, ok := e.Value.(*MessageCacheEvent).Get(key); ok {
			return v, true
		}
	}
	return nil, false
}

func (cc *commitCache) size() int {
	return len(cc.m)
}
