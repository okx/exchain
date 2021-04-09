package state

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/status-im/keycard-go/hexutils"
)

var (
	KeyPrefixPruningRoot       = []byte("pr_0x10")
	KeyPrefixPruningUpdatedKey = []byte("pr_0x11")
)

type cacheKey struct {
	T        ethstate.Trie
	PrevRoot ethcmn.Hash

	//Keys [][]byte
	Keys map[string]int
}

type CacheTrie struct {
	height uint64
	tries  map[string]*cacheKey
}

func NewCacheTrie() CacheTrie {
	return CacheTrie{
		tries:  make(map[string]*cacheKey, 0),
		height: 0,
	}
}

func (c *CacheTrie) UpdateHeight(h uint64) {
	c.height = h
}

func (c *CacheTrie) Add(prefix []byte, trie ethstate.Trie, prevRoot ethcmn.Hash) {
	keys := hexutils.BytesToHex(prefix)
	_, ok := c.tries[keys]
	if !ok {
		c.tries[keys] = &cacheKey{T: trie, PrevRoot: prevRoot, Keys: make(map[string]int, 0)}
	} else {
		c.tries[keys].T = trie
		c.tries[keys].PrevRoot = prevRoot
	}
}

func (c *CacheTrie) InsertDirtyKey(prefix, key []byte) {
	var keyHash [ethcmn.HashLength]byte
	keys := hexutils.BytesToHex(prefix)
	trie, ok := c.tries[keys]
	if !ok {
		c.tries[keys] = &cacheKey{T: nil, Keys: make(map[string]int, 0)}
		trie = c.tries[keys]
	}
	hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
	hasher.Reset()
	hasher.Write(key)
	_, e := hasher.Read(keyHash[:])
	keyStr := hexutils.BytesToHex(keyHash[:])

	_, ok = trie.Keys[keyStr]
	if ok {
		return
	}

	if e == nil {
		trie.Keys[keyStr] = 0
	}
}

func (c *CacheTrie) Reset() {
	c.tries = make(map[string]*cacheKey, 0)
}

func (c *CacheTrie) Commit() {
	var roots []ethcmn.Hash
	for _, v := range c.tries {
		if v.T == nil {
			continue
		}
		hash, e := v.T.Commit(nil)
		if e != nil {
			panic(e)
		}
		roots = append(roots, hash)
	}
	trieDB := InstanceOfStateStore().GetDb().TrieDB()
	for _, root := range roots {
		fmt.Println("commit root hash :" + hexutils.BytesToHex(root.Bytes()))
		e := trieDB.Commit(root, false, nil)
		if e != nil {
			panic(e)
		}
	}
	//save cached key which has been updated by contract
	c.CommitDirtyKey()
	for addr, v := range c.tries {
		fmt.Println("CommitPruning Root new root :" + hexutils.BytesToHex(v.T.Hash().Bytes()))
		fmt.Println("CommitPruning Root old root :" + hexutils.BytesToHex(v.PrevRoot.Bytes()))
		c.CommitPruningRoot(hexutils.HexToBytes(addr), v.PrevRoot.Bytes())
	}

}

func (c *CacheTrie) CommitPruningRoot(addr, prevRoot []byte) {
	InstanceOfStateStore().CommitPruningRoot(c.height, addr, prevRoot)
}

func (c *CacheTrie) CommitDirtyKey() {
	diskDB := InstanceOfStateStore().GetDb().TrieDB().DiskDB()
	prefix := PruningDirtyKeyPrefix(c.height)

	for _, v := range c.tries {
		idx := 0
		for k, _ := range v.Keys {
			diskDB.Put(append(append(prefix, v.T.Hash().Bytes()...), []byte(strconv.Itoa(idx))...), hexutils.HexToBytes(k))
			idx++
		}
	}
}

func PruningRootPrefix(height uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, height)
	return append(KeyPrefixPruningRoot, key...)
}

func PruningDirtyKeyPrefix(height uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, height)
	return append(KeyPrefixPruningUpdatedKey, key...)
}
