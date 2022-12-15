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

func BenchmarkInsertGasTxQueue(b *testing.B) {

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
		gq := NewHeapQueue()
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
		gq := NewHeapQueue()
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
		gq := NewHeapQueue()
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
	hq := NewHeapQueue().(*HeapQueue)
	for i := 0; i < lenght; i++ {
		err := gq.Insert(txs[i])
		require.NoError(t, err)
		err = hq.Insert(txs[i])
		require.NoError(t, err)
	}

	for e := gq.Front(); e != nil; e = e.Next() {
		t.Log("gq from:", e.Address, "nonce", e.Nonce, "gp", e.GasPrice.String())
	}

	hq.Init()
	tx := hq.Peek()
	for tx != nil {
		t.Log("hq from:", tx.from, "nonce", tx.realTx.GetNonce(), "gp", tx.realTx.GetGasPrice().String())
		hq.Shift()
		tx = hq.Peek()
	}

}

func BenchmarkInsertGasTxQueue_1(b *testing.B) {

	b.Run("gas queue reserve ", func(b *testing.B) {
		gqs := make([]*GasTxQueue, 0)
		for i := 0; i < b.N; i++ {
			gqs = append(gqs, NewGasTxQueue(10))
		}
		txs := make([]*mempoolTx, 200000)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < 200000; i++ {
			n := big.NewInt(int64(200000 - i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
			for j := 0; j < b.N; j++ {
				err := gqs[j].Insert(txs[i])
				require.NoError(b, err)
			}
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for e := gqs[i].Front(); e != nil; e = e.Next() {

			}
		}
	})

	b.Run("heap queue reserve", func(b *testing.B) {
		gqs := make([]*HeapQueue, 0)
		for i := 0; i < b.N; i++ {
			gqs = append(gqs, NewHeapQueue().(*HeapQueue))
		}

		txs := make([]*mempoolTx, 200000)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < 200000; i++ {
			n := big.NewInt(int64(i))
			tx := generateMemepool(fmt.Sprintf("%d", i), 0, n)
			txs[i] = tx
			for j := 0; j < b.N; j++ {
				err := gqs[j].Insert(txs[i])
				require.NoError(b, err)
			}
		}

		b.StartTimer()

		for i := 0; i < b.N; i++ {
			gqs[i].Init()
			j := 0
			tx := gqs[i].Peek()
			for tx != nil {
				gqs[i].Shift()
				if j > 20000 {
					break
				}
				j++
				tx = gqs[i].Peek()
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
	hq := NewHeapQueue().(*HeapQueue)
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
