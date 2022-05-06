package mpt

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/require"
)

var cases = []struct {
	num int
}{
	{0},
	{1},
	{2},
	{100},
	{1000},
	{10000},
}

func Test_Store_Iterate(t *testing.T) {
	for i, c := range cases {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			trie, kvs := fullFillStore(c.num)
			iter := newMptIterator(trie, nil, nil)
			defer iter.Close()
			count := 0
			iKvs := make(map[string]string, c.num)
			for ; iter.Valid(); iter.Next() {
				require.NotNil(t, iter.Key())
				iKvs[string(iter.Key())] = string(iter.Value())
				count++
			}
			require.EqualValues(t, kvs, iKvs)
			require.Equal(t, c.num, len(iKvs))
			require.Equal(t, c.num, count)
		})
	}
}

func fullFillStore(num int) (ethstate.Trie, map[string]string) {
	db := ethstate.NewDatabase(rawdb.NewMemoryDatabase())
	tr, err := db.OpenTrie(NilHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}

	kvs := make(map[string]string, num)
	for i := 0; i < num; i++ {
		k, v := fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i)
		kvs[k] = v
		if err := tr.TryUpdate([]byte(k), []byte(v)); err != nil {
			panic(err)
		}
	}
	return tr, kvs
}
