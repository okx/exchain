package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	data := make(map[string]ethcmn.Address)
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	b := bufio.NewReader(f)
	for {
		k, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		v, _, err := b.ReadLine()
		if err != nil {
			panic(err)
		}
		data[string(k)] = ethcmn.HexToAddress(string(v))

	}
	c.data = data
	fmt.Println("verify sig cache size:", len(c.data))
	for k, v := range c.data {
		fmt.Println(k, v.String())
	}
}

func (c *Cache) Save(fileName string) {
	fmt.Println("verify sig cache size:", len(c.data))
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range c.data {
		fmt.Println(k, v.String())
		w.Write([]byte(k))
		w.WriteByte('\n')

		w.WriteString(v.String())
		w.WriteByte('\n')
	}
	w.Flush()
}
