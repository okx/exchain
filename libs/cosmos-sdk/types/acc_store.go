package types

import "github.com/ethereum/go-ethereum/trie"

type AccCacheStore interface {
	Get(key AccAddress) (value interface{})
	Set(key AccAddress, value interface{})
	Has(key AccAddress) bool
	Delete(key AccAddress)
	Write()
	NewIterator(startKey []byte) *trie.Iterator
}
