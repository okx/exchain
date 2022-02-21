package env

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
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
	mtx  sync.RWMutex
	data map[string]ethcmn.Address
}

func newCache() *Cache {
	return &Cache{
		data: make(map[string]ethcmn.Address, cacheSize),
	}
}

func (c *Cache) Get(key string) (ethcmn.Address, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if !validateKey(key) {
		return ethcmn.Address{}, false
	}
	if value, ok := c.data[key]; ok {
		return value, true
	}
	return ethcmn.Address{}, false
}
func (c *Cache) Add(key string, value ethcmn.Address) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if !validateKey(key) {
		return
	}
	c.data[key] = value
}

func validateKey(key string) bool {
	if key == "" {
		return false
	}
	return true
}

func (c *Cache) Load(fileName string) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var data map[string]ethcmn.Address
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	c.data = data
}

func (c *Cache) Save(fileName string) {
	content, err := json.Marshal(c.data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fileName, content, fs.ModePerm)
	if err != nil {
		panic(err)
	}
}
