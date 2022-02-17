package env

import (
	"container/list"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

var (
	VerifySigCache *Cache
	once           sync.Once
)

const cacheSize = 1000000

func init() {
	once.Do(func() {
		VerifySigCache = newCache()
	})
}

type Cache struct {
	mtx   sync.RWMutex
	data  map[string]*list.Element
	queue *list.List
}

type cacheNode struct {
	key   string
	value ethcmn.Address
}

func newCache() *Cache {
	return &Cache{
		data:  make(map[string]*list.Element, cacheSize),
		queue: list.New(),
	}
}

func (c *Cache) Get(key string) (ethcmn.Address, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if elem, ok := c.data[key]; ok {
		c.queue.MoveToBack(elem)
		return elem.Value.(*cacheNode).value, true
	}
	return ethcmn.Address{}, false
}
func (c *Cache) Add(key string, value ethcmn.Address) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	node := &cacheNode{key, value}
	elem := c.queue.PushBack(node)
	c.data[key] = elem

	for c.queue.Len() > cacheSize {
		oldest := c.queue.Front()
		oldKey := c.queue.Remove(oldest).(*cacheNode).key
		delete(c.data, oldKey)
	}
}
