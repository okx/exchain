package baseapp

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testKey = []byte("testKey")

func TestHguBasic(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	jobQueueLen = 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var hgu *HguRecord
	for i := 0; i < 100000; i++ {
		gu := r.Int63n(31)*10 + 100 // [100, 400], step=10
		InstanceOfHistoryGasUsedRecordDB().UpdateGasUsed(testKey, gu)
		InstanceOfHistoryGasUsedRecordDB().FlushHgu()
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		if hgu == nil {
			continue
		}
		require.True(t, hgu.MaxGas >= hgu.HighGas, fmt.Sprintf("MaxGas: %d, HighGas: %d\n", hgu.MaxGas, hgu.HighGas))
		require.True(t, hgu.HighGas >= hgu.StandardGas, fmt.Sprintf("HighGas: %d, StandardGas: %d\n", hgu.HighGas, hgu.StandardGas))
	}
	t.Log(hgu)
}

func TestHguSmooth(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var hgu *HguRecord
	cases := []struct {
		gasLimit int64
		floor    int64
		ceil     int64
	}{
		{
			20000,
			20000,
			20000,
		},
		{
			21000,
			21000,
			21000,
		},
		{
			21050,
			21000,
			21050,
		},
		{
			23000,
			21000,
			22000,
		},
	}

	const baseGasUsed = 21000
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: baseGasUsed, key: testKey})

	for i := 0; i < 100; i++ {
		gu := baseGasUsed + r.Int63n(101) // [21000, 21100]
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.True(t, hgu.MaxGas >= hgu.HighGas, fmt.Sprintf("MaxGas: %d, HighGas: %d\n", hgu.MaxGas, hgu.HighGas))
		require.True(t, hgu.HighGas >= hgu.StandardGas, fmt.Sprintf("HighGas: %d, StandardGas: %d\n", hgu.HighGas, hgu.StandardGas))

		for _, c := range cases {
			gc := estimateGas(c.gasLimit, hgu)
			require.LessOrEqual(t, c.floor, gc, fmt.Sprintf("index: %d, hgu: %v\n", i, hgu))
			require.LessOrEqual(t, gc, c.ceil, fmt.Sprintf("hgu: %v\n", hgu))
		}
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: gu, key: testKey})
	}
	t.Log(hgu)
}

func TestHgu10PercentChange(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})

	const (
		gasLimit    = int64(50000000)
		baseGasUsed = int64(6746245)
		gasFloor    = int64(21000)
	)

	gasUsed := baseGasUsed

	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: gasUsed, key: testKey})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var hgu *HguRecord
	_ = hgu
	var totalChangedPercent int64
	var highChangedCount int
	totalCount := 100000
	for i := 0; i < totalCount; i++ {
		if gasUsed < gasFloor+gasFloor/5 {
			// add 10%
			gasUsed = gasUsed + gasUsed/10
		} else if gasUsed > gasLimit-gasLimit/5 {
			// sub 10%
			gasUsed = gasUsed - gasUsed/10
		} else {
			// random add 10% or sub 10%
			gasUsed = gasUsed - gasUsed/10 + gasUsed/5*r.Int63n(2)
		}
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		eg := estimateGas(gasUsed*3/2, hgu)
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: gasUsed, key: testKey})
		changedPercent := absInt64(gasUsed*10000/eg - 10000)
		totalChangedPercent += changedPercent
		if changedPercent > 2000 {
			highChangedCount += 1
		}
	}
	t.Logf("Mean deviation rate: %.2f%%, deviation more than 20%% rate: %.2f%%\n", float64(totalChangedPercent)/100/float64(totalCount), float64(highChangedCount)*100/float64(totalCount))

}

func TestHguRandomChange(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})

	const (
		gasLimit    = int64(50000000)
		baseGasUsed = int64(6746245)
		gasFloor    = int64(21000)
	)

	gasUsed := baseGasUsed

	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: gasUsed, key: testKey})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var hgu *HguRecord
	var totalChangedPercent int64
	var highChangedCount int
	totalCount := 100000
	for i := 0; i < totalCount; i++ {
		if gasUsed < gasFloor+gasFloor/5 {
			// add 10%
			gasUsed = gasUsed + gasUsed/10
		} else if gasUsed > gasLimit-gasLimit/5 {
			// sub 10%
			gasUsed = gasUsed - gasUsed/10
		} else if r.Int63n(100) == 0 {
			// random change
			gasUsed = gasFloor + r.Int63n(baseGasUsed)
		} else {
			// random add 10% or sub 10%
			gasUsed = gasUsed - gasUsed/10 + gasUsed/5*r.Int63n(2)
		}
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		eg := estimateGas(gasUsed*3/2, hgu)
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: gasUsed, key: testKey})
		changedPercent := absInt64(gasUsed*10000/eg - 10000)
		totalChangedPercent += changedPercent
		if changedPercent > 5000 {
			highChangedCount += 1
		}
	}
	t.Logf("Mean deviation rate: %.2f%%, deviation more than 50%% rate: %.2f%%\n", float64(totalChangedPercent)/100/float64(totalCount), float64(highChangedCount)*100/float64(totalCount))
}

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// TestHguXEN tests `t` method of 0x6f0a55cd633Cc70BeB0ba7874f3B010C002ef59f on okc mainnet
// mainnet height: [15432438, 15432445]
func TestHguXEN(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	cases := []struct {
		height   int64
		gasUsed  int64
		gasLimit int64
	}{
		{
			15432438,
			33571582,
			50000000,
		},
		{
			15432439,
			33571582,
			50000000,
		},
		{
			15432440,
			33571582,
			50000000,
		},
		{
			15432441,
			6746245,
			8400000,
		},
		{
			15432442,
			33571582,
			50000000,
		},
		{
			15432443,
			33571582,
			50000000,
		},
		{
			15432444,
			33571582,
			50000000,
		},
		{
			15432445,
			33571582,
			50000000,
		},
	}
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: cases[0].gasUsed, key: testKey})
	var hgu *HguRecord
	for _, c := range cases {
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.LessOrEqual(t, hgu.HighGas, hgu.MaxGas)
		require.LessOrEqual(t, hgu.StandardGas, hgu.HighGas)
		t.Log(hgu, c.gasUsed, estimateGas(c.gasLimit, hgu))
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: c.gasUsed, key: testKey})
	}
}

func TestHguMock1(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	cases := []struct {
		gasUsed  int64
		gasLimit int64
	}{
		{
			33571582,
			50000000,
		},
		{
			33571582,
			50000000,
		},
		{
			6746245,
			8400000,
		},
		{
			6746245,
			8400000,
		},
		{
			6746245,
			8400000,
		},
		{
			6746245,
			8400000,
		},
		{
			33571582,
			50000000,
		},
		{
			33571582,
			50000000,
		},
		{
			33571582,
			50000000,
		},
	}
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: cases[0].gasUsed, key: testKey})
	var hgu *HguRecord
	for _, c := range cases {
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.LessOrEqual(t, hgu.HighGas, hgu.MaxGas)
		require.LessOrEqual(t, hgu.StandardGas, hgu.HighGas)
		t.Log(hgu, c.gasUsed, estimateGas(c.gasLimit, hgu))
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: c.gasUsed, key: testKey})
	}
}

func TestHguMock2(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	cases := []struct {
		height   int64
		gasUsed  int64
		gasLimit int64
	}{
		{
			15432438,
			6746245,
			8400000,
		},
		{
			15432439,
			6746245,
			8400000,
		},
		{
			15432440,
			6746245,
			8400000,
		},
		{
			15432441,
			33571582,
			50000000,
		},
		{
			15432442,
			33571582,
			50000000,
		},
		{
			15432443,
			33571582,
			50000000,
		},
		{
			15432444,
			33571582,
			50000000,
		},
		{
			15432445,
			6746245,
			8400000,
		},
		{
			15432446,
			6746245,
			8400000,
		},
	}
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: cases[0].gasUsed, key: testKey})
	var hgu *HguRecord
	for _, c := range cases {
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.LessOrEqual(t, hgu.HighGas, hgu.MaxGas)
		require.LessOrEqual(t, hgu.StandardGas, hgu.HighGas)
		t.Log(hgu, c.gasUsed, estimateGas(c.gasLimit, hgu))
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: c.gasUsed, key: testKey})
	}
}

func TestHguMock3(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	cases := []struct {
		height   int64
		gasUsed  int64
		gasLimit int64
	}{
		{
			15432438,
			33571582,
			50000000,
		},
		{
			15432439,
			33571582,
			50000000,
		},
		{
			15432440,
			6746245,
			50000000,
		},
		{
			15432441,
			6746245,
			50000000,
		},
		{
			15432442,
			6746245,
			50000000,
		},
		{
			15432443,
			6746245,
			50000000,
		},
		{
			15432444,
			33571582,
			50000000,
		},
		{
			15432445,
			33571582,
			50000000,
		},
		{
			15432446,
			33571582,
			50000000,
		},
	}
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: cases[0].gasUsed, key: testKey})
	var hgu *HguRecord
	for _, c := range cases {
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.LessOrEqual(t, hgu.HighGas, hgu.MaxGas)
		require.LessOrEqual(t, hgu.StandardGas, hgu.HighGas)
		t.Log(hgu, c.gasUsed, estimateGas(c.gasLimit, hgu))
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: c.gasUsed, key: testKey})
	}
}

func TestHguMock4(t *testing.T) {
	t.Cleanup(func() {
		InstanceOfHistoryGasUsedRecordDB().close()
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	cases := []struct {
		height   int64
		gasUsed  int64
		gasLimit int64
	}{
		{
			15432438,
			6746245,
			50000000,
		},
		{
			15432439,
			6746245,
			50000000,
		},
		{
			15432440,
			33571582,
			50000000,
		},
		{
			15432441,
			33571582,
			50000000,
		},
		{
			15432442,
			33571582,
			50000000,
		},
		{
			15432443,
			33571582,
			50000000,
		},
		{
			15432444,
			6746245,
			50000000,
		},
		{
			15432445,
			6746245,
			50000000,
		},
		{
			15432446,
			6746245,
			50000000,
		},
	}
	InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: cases[0].gasUsed, key: testKey})
	var hgu *HguRecord
	for _, c := range cases {
		hgu = InstanceOfHistoryGasUsedRecordDB().GetHgu(testKey)
		require.NotNil(t, hgu)
		require.LessOrEqual(t, hgu.HighGas, hgu.MaxGas)
		require.LessOrEqual(t, hgu.StandardGas, hgu.HighGas)
		t.Log(hgu, c.gasUsed, estimateGas(c.gasLimit, hgu))
		InstanceOfHistoryGasUsedRecordDB().flushHgu(gasKey{gas: c.gasUsed, key: testKey})
	}
}
