package mempool

import (
	"fmt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func generateMemepool(from string, nonce uint64, gasPrice *big.Int) *mempoolTx {
	builder := strings.Builder{}
	builder.WriteString(from)
	builder.WriteString(fmt.Sprintf("%d", nonce))
	builder.WriteString(gasPrice.String())
	return &mempoolTx{height: 1, gasWanted: 1, tx: []byte(builder.String()), from: from, realTx: abci.MockTx{GasPrice: gasPrice, Nonce: nonce, From: from}}
}

func Benchmark_GasTxQueue_Insert(b *testing.B) {

	b.Run("gas queue", func(b *testing.B) {
		gq := NewGasTxQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			n := big.NewInt(int64(b.N - i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})

	b.Run("gas queue reserve ", func(b *testing.B) {
		gq := NewGasTxQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			n := big.NewInt(int64(i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})

	b.Run("gas queue rand", func(b *testing.B) {
		gq := NewGasTxQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		threhold := big.NewInt(int64(b.N))
		for i := 0; i < b.N; i++ {
			n := threhold.Rand(rand.New(rand.NewSource(time.Now().Unix())), threhold)
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})

	b.Run("heap queue", func(b *testing.B) {
		gq := NewHeapQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			n := big.NewInt(int64(b.N - i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})

	b.Run("heap queue reserve", func(b *testing.B) {
		gq := NewHeapQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			n := big.NewInt(int64(i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})

	b.Run("heap queue rand", func(b *testing.B) {
		gq := NewHeapQueue(10)
		txs := make([]*mempoolTx, b.N)
		b.ResetTimer()
		b.StopTimer()
		threhold := big.NewInt(int64(b.N))
		for i := 0; i < b.N; i++ {
			n := threhold.Rand(rand.New(rand.NewSource(time.Now().Unix())), threhold)
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := gq.Insert(txs[i])
			require.NoError(b, err)
		}
		b.ReportAllocs()
	})
}

func TestQueue_Back(t *testing.T) {
	lenght := 10
	txs := make([]*mempoolTx, 10)
	for i := 0; i < lenght; i++ {
		n := big.NewInt(int64(i))
		tx := generateMemepool(fmt.Sprintf("%d", 1), uint64(i), n)
		txs[i] = tx
	}
	gq := NewGasTxQueue(10)
	hq := NewHeapQueue(10).(*HeapQueue)
	for i := 0; i < lenght; i++ {
		err := gq.Insert(txs[i])
		require.NoError(t, err)
		err = hq.Insert(txs[i])
		require.NoError(t, err)
	}

	for e := gq.Front(); e != nil; e = e.Next() {
		t.Log("gq from:", e.Address, "nonce", e.Nonce, "gp", e.GasPrice.String())
	}

	heads := hq.Init()
	tx := hq.Peek(heads)
	for tx != nil {
		t.Log("hq from:", tx.from, "nonce", tx.realTx.GetNonce(), "gp", tx.realTx.GetGasPrice().String())
		hq.Shift(&heads)
		tx = hq.Peek(heads)
	}

}

func Benchmark_GasTxQueue_Reap(b *testing.B) {
	mempoolTxSize := 200000
	txsize := 20000
	gq := NewGasTxQueue(10)
	hq := NewHeapQueue(10).(*HeapQueue)
	//mod := 2
	for i := 0; i < mempoolTxSize; i++ {
		n := big.NewInt(int64(mempoolTxSize - i))
		tx := generateMemepool(fmt.Sprintf("%d", i), uint64(i), n)
		err := gq.Insert(tx)
		require.NoError(b, err)
		err = hq.Insert(tx)
		require.NoError(b, err)
	}

	b.Run("gas queue reserve ", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			j := 0
			for e := gq.Front(); e != nil; e = e.Next() {
				if j > txsize {
					break
				}
				j++
			}
		}
	})

	b.Run("heap queue reserve", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			heads := hq.Init()
			j := 0
			tx := hq.Peek(heads)
			for tx != nil {
				hq.Shift(&heads)
				if j > txsize {
					break
				}
				j++
				tx = hq.Peek(heads)
			}
		}
	})

}

func TestHeapQueue_CleanItems(t *testing.T) {
	lenght := 10
	txs := make([]*mempoolTx, 2*lenght)
	for i := 0; i < 2*lenght; i++ {
		n := big.NewInt(int64(2*lenght - i))
		tx := generateMemepool(fmt.Sprintf("%d", 1), uint64(i), n)
		txs[i] = tx
	}
	hq := NewHeapQueue(10).(*HeapQueue)
	for i := 0; i < lenght; i++ {
		err := hq.Insert(txs[i])
		require.NoError(t, err)
	}
	done := make(chan int)
	go func() {
		time.Sleep(time.Millisecond * 10)
		hq.CleanItems(fmt.Sprintf("%d", 1), uint64(lenght))
		done <- 1
	}()
	for i := lenght + 2; i < 2*lenght; i++ {
		go func(index int) {
			err := hq.Insert(txs[index])
			require.NoError(t, err)
		}(i)
	}
	<-done

}

func TestHeapQueue_Insert(t *testing.T) {
	length := 10
	hq := NewHeapQueue(10).(*HeapQueue)
	txs := make([]*mempoolTx, length)

	for i := 0; i < length; i++ {
		n := big.NewInt(int64(i))
		tx := generateMemepool(fmt.Sprintf("%d", 1), 0, n)
		txs[i] = tx
	}

	for i := 0; i < length; i++ {
		hq.Insert(txs[i])
		fmtHq(hq, t)
		//require.NoError(t, err)
	}

	heads := hq.Init()
	tx := hq.Peek(heads)
	for tx != nil {
		t.Log("hq from:", tx.from, "nonce", tx.realTx.GetNonce(), "gp", tx.realTx.GetGasPrice().String())
		hq.Shift(&heads)
		tx = hq.Peek(heads)
	}
}

func TestHeapQueue_Insert2(t *testing.T) {
	hq := NewHeapQueue(10).(*HeapQueue)
	address1 := fmt.Sprintf("%d", 1)
	//address2 := fmt.Sprintf("%d", 2)
	tx1 := generateMemepool(address1, 0, big.NewInt(100))
	tx2 := generateMemepool(address1, 1, big.NewInt(100))
	tx3 := generateMemepool(address1, 0, big.NewInt(110))
	tx4 := generateMemepool(address1, 0, big.NewInt(111))

	hq.Insert(tx1)
	hq.Insert(tx2)
	hq.txs[address1].Back().Nonce = 0
	fmtHq(hq, t)

	err := hq.Insert(tx3)
	t.Log(err)
	require.Error(t, err)

	defer func() {
		if r := recover(); r != nil {
			t.Log("unit test", r)
		}
	}()
	err = hq.Insert(tx4)
	require.NoError(t, err)
	//heads := hq.Init()
	//tx := hq.Peek(heads)
	//for tx != nil {
	//	t.Log("hq from:", tx.from, "nonce", tx.realTx.GetNonce(), "gp", tx.realTx.GetGasPrice().String())
	//	hq.Shift(&heads)
	//	tx = hq.Peek(heads)
	//}
}

func fmtHq(hq *HeapQueue, t *testing.T) {
	t.Log("HeapQueue -------------start-------------")
	for k, v := range hq.txs {
		strs := make([]string, 0)
		for e := v.Front(); e != nil; e = e.Next() {
			strs = append(strs, fmt.Sprintf("from: %s nonce:%d gp:%d", e.Address, e.Nonce, e.GasPrice.Int64()))
		}
		t.Log("Address: ", k, "list", strs)

	}
	t.Log("HeapQueue -------------end-------------")
}
