package mpt

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
)

/*
 * these tests are copied from the go-ethereum/core/state/trie_prefetcher_test.go
 */
func filledStateDB() (*ethstate.StateDB, ethstate.Database, common.Hash) {
	db := ethstate.NewDatabase(rawdb.NewMemoryDatabase())
	originalRoot := common.Hash{}
	state, _ := ethstate.New(originalRoot, db, nil)

	// Create an account and check if the retrieved balance is correct
	addr := common.HexToAddress("0xaffeaffeaffeaffeaffeaffeaffeaffeaffeaffe")
	skey := common.HexToHash("aaa")
	sval := common.HexToHash("bbb")

	state.SetBalance(addr, big.NewInt(42)) // Change the account trie
	state.SetCode(addr, []byte("hello"))   // Change an external metadata
	state.SetState(addr, skey, sval)       // Change the storage trie
	for i := 0; i < 100; i++ {
		sk := common.BigToHash(big.NewInt(int64(i)))
		state.SetState(addr, sk, sk) // Change the storage trie
	}
	return state, db, originalRoot
}

func TestCopyAndClose(t *testing.T) {
	_, db, originalRoot := filledStateDB()
	prefetcher := NewTriePrefetcher(db, originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	prefetcher.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	time.Sleep(1 * time.Second)
	a := prefetcher.Trie(originalRoot)
	prefetcher.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	b := prefetcher.Trie(originalRoot)
	cpy := prefetcher.Copy()
	cpy.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	cpy.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	c := cpy.Trie(originalRoot)
	prefetcher.Close()
	cpy2 := cpy.Copy()
	cpy2.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	d := cpy2.Trie(originalRoot)
	cpy.Close()
	cpy2.Close()
	if a.Hash() != b.Hash() || a.Hash() != c.Hash() || a.Hash() != d.Hash() {
		t.Fatalf("Invalid trie, hashes should be equal: %v %v %v %v", a.Hash(), b.Hash(), c.Hash(), d.Hash())
	}
}

func TestUseAfterClose(t *testing.T) {
	_, db, originalRoot := filledStateDB()
	prefetcher := NewTriePrefetcher(db, originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	a := prefetcher.Trie(originalRoot)
	prefetcher.Close()
	b := prefetcher.Trie(originalRoot)
	if a == nil {
		t.Fatal("Prefetching before close should not return nil")
	}
	if b != nil {
		t.Fatal("Trie after close should return nil")
	}
}

func TestCopyClose(t *testing.T) {
	_, db, originalRoot := filledStateDB()
	prefetcher := NewTriePrefetcher(db, originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.Prefetch(originalRoot, [][]byte{skey.Bytes()})
	cpy := prefetcher.Copy()
	a := prefetcher.Trie(originalRoot)
	b := cpy.Trie(originalRoot)
	prefetcher.Close()
	c := prefetcher.Trie(originalRoot)
	d := cpy.Trie(originalRoot)
	if a == nil {
		t.Fatal("Prefetching before close should not return nil")
	}
	if b == nil {
		t.Fatal("Copy trie should return nil")
	}
	if c != nil {
		t.Fatal("Trie after close should return nil")
	}
	if d == nil {
		t.Fatal("Copy trie should not return nil")
	}
}
