//go:build rocksdb
// +build rocksdb

package types

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRocksDBBackend(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewWrapRocksDB(name, dir)
	assert.Nil(t, err, "fail to create wrap rocksdb")

	batch := db.NewBatch()
	for i := 0; i < 10000; i++ {
		batch.Put([]byte(fmt.Sprintf("key-%d", i)), []byte(fmt.Sprintf("value-%d", i)))
	}
	err = batch.Write()
	assert.Nil(t, err, "fail to test wrap rocksdb's batch")

	itr := db.NewIterator(nil, nil)
	defer itr.Release()

	for itr.Next() {
		fmt.Println("key is: ", itr.Key(), ", value is: ", itr.Value())
	}
}

const (
	strChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // 62 characters
)

// Str constructs a random alphanumeric string of given length.
func randStr(length int) string {
	chars := []byte{}
MAIN_LOOP:
	for {
		val := rand.Int63() // nolint:gosec // G404: Use of weak random number generator
		for i := 0; i < 10; i++ {
			v := int(val & 0x3f) // rightmost 6 bits
			if v >= 62 {         // only 62 characters in strChars
				val >>= 6
				continue
			} else {
				chars = append(chars, strChars[v])
				if len(chars) == length {
					break MAIN_LOOP
				}
				val >>= 6
			}
		}
	}

	return string(chars)
}
