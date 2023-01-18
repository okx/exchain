package baseapp

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type hguTestCase struct {
	key     string
	gasUsed int64
}

func TestHGU(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	testCase := []hguTestCase{
		{"test0", 1},
		{"test0", 2},
		{"test0", 3},
		{"test0", 4},
		{"test0", 5},
		{"test1", 10},
		{"test2", 20},
	}
	expected := []struct {
		key    string
		record HguRecord
	}{
		{
			"test0",
			HguRecord{
				MaxGas:           3,
				MinGas:           3,
				MovingAverageGas: 3,
			},
		},
		{
			"test1",
			HguRecord{
				MaxGas:           10,
				MinGas:           10,
				MovingAverageGas: 10,
			},
		},
		{
			"test2",
			HguRecord{
				MaxGas:           20,
				MinGas:           20,
				MovingAverageGas: 20,
			},
		},
	}
	hguDB := InstanceOfHistoryGasUsedRecordDB()
	for _, c := range testCase {
		hguDB.UpdateGasUsed([]byte(c.key), c.gasUsed)
	}
	for _, c := range expected {
		r, ok := hguDB.getHgu([]byte(c.key))
		require.False(t, ok)
		require.Nil(t, r)
	}

	hguDB.FlushHgu()
	time.Sleep(time.Second)
	for _, c := range expected {
		r, ok := hguDB.getHgu([]byte(c.key))
		require.False(t, ok)
		require.Equal(t, c.record, *r)
	}

	for _, c := range expected {
		r := hguDB.GetHgu([]byte(c.key))
		require.Equal(t, c.record, *r)
	}

	for _, c := range expected {
		r, ok := hguDB.getHgu([]byte(c.key))
		require.True(t, ok)
		require.Equal(t, c.record, *r)
	}

	hguDB.UpdateGasUsed([]byte("test0"), 1)
	hguDB.FlushHgu()
	hguDB.UpdateGasUsed([]byte("test0"), 10)
	hguDB.FlushHgu()
	time.Sleep(time.Second)

	r, ok := hguDB.getHgu([]byte("test0"))
	require.True(t, ok)
	require.Equal(t, int64(10), r.MaxGas)
	require.Equal(t, int64(1), r.MinGas)
	require.Equal(t, int64(5), r.MovingAverageGas)
	require.Equal(t, int64(2), r.BlockNum)
}

func TestMovingAverageGas(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll(HistoryGasUsedDbDir)
	})
	testKey := []byte("test")
	cases := []int64{21000, 23000, 25000, 33000, 37000, 53000}
	hguDB := InstanceOfHistoryGasUsedRecordDB()
	for _, gas := range cases {
		hguDB.UpdateGasUsed(testKey, gas)
		hguDB.FlushHgu()
	}
	time.Sleep(time.Second)

	r := hguDB.GetHgu(testKey)
	require.Equal(t, int64(53000), r.MaxGas)
	require.Equal(t, int64(21373), r.MinGas)
	require.Equal(t, int64(39816), r.MovingAverageGas)
	require.Equal(t, int64(5), r.BlockNum)
}
