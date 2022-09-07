package watcher

import (
	"container/list"
	"encoding/hex"
	"fmt"
	"sync"
)

type MessageCache struct {
	mtx   sync.RWMutex
	count int
	mp    map[string]WatchMessage // if the key of value WatchMessage is nil, this key should del on db batch write
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
	c.count++
	c.mp[hex.EncodeToString(wsg.GetKey())] = wsg
}

func (c *MessageCache) BatchDel(keys [][]byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count += len(keys)
	for _, k := range keys {
		c.mp[hex.EncodeToString(k)] = &Batch{Key: k, TypeValue: TypeDelete}
	}
}

func (c *MessageCache) BatchSet(wsgs []WatchMessage) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.count += len(wsgs)
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
	c.count += len(batchs)
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

type Stat struct {
	count      int // kv count
	dbSize     int // db storage
	structSize int //
}

func (c *MessageCache) Clear() map[string]*Stat {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	static := make(map[string]*Stat)
	for k, v := range c.mp {
		delete(c.mp, k)
		//just for test
		Statistic(v, static)
	}
	for k, v := range static {
		dbsize := float64(v.dbSize) / float64(1024*1024)
		structsize := float64(v.structSize) / float64(1024*1024)
		fmt.Printf("**** lyh ****** static %s, count %d, dbSize %.3f, structSize %.3f \n", k, v.count, dbsize, structsize)
	}
	return static
}

func Statistic(wmg WatchMessage, stat map[string]*Stat) {
	key := wmg.GetKey()
	value := []byte(wmg.GetValue())
	if len(key) < 1 {
		return
	}
	switch string(key[0:1]) {
	case string(prefixTx):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixTx"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixTx"] = v
	case string(prefixBlock):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixBlock"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixBlock"] = v
	case string(prefixReceipt):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixReceipt"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixReceipt"] = v
	case string(prefixCode):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixCode"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixCode"] = v
	case string(prefixBlockInfo):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixBlockInfo"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixBlockInfo"] = v
	case string(prefixLatestHeight):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixLatestHeight"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixLatestHeight"] = v
	case string(prefixAccount):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixAccount"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixAccount"] = v
	case string(PrefixState):
		var v *Stat
		var ok bool
		if v, ok = stat["PrefixState"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["PrefixState"] = v
	case string(prefixCodeHash):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixCodeHash"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixCodeHash"] = v
	case string(prefixParams):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixParams"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixParams"] = v
	case string(prefixWhiteList):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixWhiteList"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixWhiteList"] = v
	case string(prefixBlackList):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixBlackList"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixBlackList"] = v
	case string(prefixRpcDb):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixRpcDb"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixRpcDb"] = v
	case string(prefixTxResponse):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixTxResponse"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixTxResponse"] = v
	case string(prefixStdTxHash):
		var v *Stat
		var ok bool
		if v, ok = stat["prefixStdTxHash"]; !ok {
			v = &Stat{}
		}
		v.count++
		v.dbSize += len(key) + len(value)
		v.structSize += getSize(wmg)
		stat["prefixStdTxHash"] = v
	default:
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
	cc.mtx.RLock()
	defer cc.mtx.RUnlock()
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
	cc.mtx.RLock()
	defer cc.mtx.RUnlock()
	return len(cc.m)
}
