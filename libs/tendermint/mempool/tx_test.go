package mempool

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestMempoolTxListList_Reset(t *testing.T) {
	list := NewMempoolTxList(func(wtx1, wtx2 *mempoolTx) bool {
		return wtx1.height >= wtx2.height
	})

	require.Zero(t, list.Size())

	for i := 0; i < 100; i++ {
		list.Insert(&mempoolTx{height: int64(i)})
	}

	require.Equal(t, 100, list.Size())

	list.Reset()
	require.Zero(t, list.Size())
}

func TestMempoolTxList_Insert(t *testing.T) {
	list := NewMempoolTxList(func(wtx1, wtx2 *mempoolTx) bool {
		return wtx1.height >= wtx2.height
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var expected []int
	for i := 0; i < 100; i++ {
		height := rng.Int63n(10000)
		expected = append(expected, int(height))
		list.Insert(&mempoolTx{height: height})

		if i%10 == 0 {
			list.Insert(&mempoolTx{height: height})
			expected = append(expected, int(height))
		}
	}

	got := make([]int, list.Size())
	for i, wtx := range list.txs {
		got[i] = int(wtx.height)
	}

	sort.Ints(expected)
	require.Equal(t, expected, got)
}

func TestMempoolTxList_Remove(t *testing.T) {
	list := NewMempoolTxList(func(wtx1, wtx2 *mempoolTx) bool {
		return wtx1.height >= wtx2.height
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var txs []*mempoolTx
	for i := 0; i < 100; i++ {
		height := rng.Int63n(10000)
		tx := &mempoolTx{height: height}

		txs = append(txs, tx)
		list.Insert(tx)

		if i%10 == 0 {
			tx = &mempoolTx{height: height}
			list.Insert(tx)
			txs = append(txs, tx)
		}
	}

	// remove a tx that does not exist
	list.Remove(&mempoolTx{height: 20000})

	// remove a tx that exists (by height) but not referenced
	list.Remove(&mempoolTx{height: txs[0].height})

	// remove a few existing txs
	for i := 0; i < 25; i++ {
		j := rng.Intn(len(txs))
		list.Remove(txs[j])
		txs = append(txs[:j], txs[j+1:]...)
	}

	expected := make([]int, len(txs))
	for i, tx := range txs {
		expected[i] = int(tx.height)
	}

	got := make([]int, list.Size())
	for i, wtx := range list.txs {
		got[i] = int(wtx.height)
	}

	sort.Ints(expected)
	require.Equal(t, expected, got)
}
