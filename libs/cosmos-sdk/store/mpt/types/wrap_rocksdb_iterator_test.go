//go:build rocksdb
// +build rocksdb

package types

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RocksDB_Iterator(t *testing.T) {
	dir := os.TempDir()
	defer os.RemoveAll(dir)

	var cases = []struct {
		num int
	}{
		{0},
		{1},
		{100},
		{1000},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			db, err := NewWrapRocksDB(fmt.Sprintf("test_rocksdb_iterator_%d", i), dir)
			assert.Nil(t, err, "fail to create wrap rocksdb")

			kvs := make(map[string][]byte, c.num)
			batch := db.NewBatch()
			for i := 0; i < c.num; i++ {
				k, v := []byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("value-%d", i))
				batch.Put(k, v)
				kvs[string(k)] = v
			}
			err = batch.Write()
			assert.Nil(t, err, "fail to test wrap rocksdb's batch")

			itr := db.NewIterator(nil, nil)
			defer itr.Release()
			iKvs := make(map[string][]byte, c.num)
			i := 0
			for itr.Next() {
				iKvs[string(itr.Key())] = itr.Value()
				i++
			}
			require.EqualValues(t, kvs, iKvs)
			require.Equal(t, c.num, i)
		})
	}
}
