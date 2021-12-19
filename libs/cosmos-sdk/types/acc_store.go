package types

import "github.com/ethereum/go-ethereum/trie"

type AccStore interface {
	Get(key string) (value []byte)
	Set(key string, value []byte)
	Has(key string) bool
	Delete(key string)
	Write()
	Clean()
	NewIterator(startKey []byte) *trie.Iterator
	CreateCacheStore() AccStore
}
