package state

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/status-im/keycard-go/hexutils"
)

type CacheTrie struct {
	tries map[string]ethstate.Trie
}

func NewCacheTrie() CacheTrie {
	return CacheTrie{
		tries: make(map[string]ethstate.Trie, 0),
	}
}

func (c *CacheTrie) Add(prefix []byte, trie ethstate.Trie) {
	_, ok := c.tries[hexutils.BytesToHex(prefix)]
	if ok {
		return
	}
	c.tries[hexutils.BytesToHex(prefix)] = trie
}

func (c *CacheTrie) Reset() {
	c.tries = make(map[string]ethstate.Trie, 0)
}

func (c *CacheTrie) Commit() {
	var roots []ethcmn.Hash
	for _, v := range c.tries {
		hash, e := v.Commit(nil)
		if e != nil {
			panic(e)
		}
		roots = append(roots, hash)
	}
	trieDB := InstanceOfStateStore().GetDb().TrieDB()
	for _, root := range roots {
		e := trieDB.Commit(root, false, nil)
		if e != nil {
			panic(e)
		}
	}
}
