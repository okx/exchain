package watcher_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/x/evm/watcher"
)

func newTestWatchData() *watcher.WatchData {
	w := setupTest()
	w.app.EvmKeeper.Watcher.Commit()
	time.Sleep(time.Second * 1)

	// get WatchData
	wdFunc := w.app.EvmKeeper.Watcher.GetWatchDataFunc()
	wd, _ := wdFunc()

	// unmarshal to raw watch data
	data := &watcher.WatchData{}
	data.UnmarshalFromAmino(nil, wd)
	return data
}
func TestWatchDataEncoder(t *testing.T) {
	w := setupTest()
	w.app.EvmKeeper.Watcher.Commit()
	time.Sleep(time.Second * 1)

	// get WatchData
	wdFunc := w.app.EvmKeeper.Watcher.GetWatchDataFunc()
	wd, err := wdFunc()
	require.NoError(t, err)

	// unmarshal to raw watch data
	data := &watcher.WatchData{}
	err = data.UnmarshalFromAmino(nil, wd)
	require.NoError(t, err)
}

// benchmark encode performance
func BenchmarkAminoEncodeDelta(b *testing.B) { benchmarkEncodeDelta(b, newEncoder("amino")) }
func BenchmarkJsonEncodeDelta(b *testing.B)  { benchmarkEncodeDelta(b, newEncoder("json")) }
func benchmarkEncodeDelta(b *testing.B, enc encoder) {
	// produce WatchData
	wd := newTestWatchData()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.encodeFunc(wd)
	}

}

// benchmark decode performance
func BenchmarkAminoDecodeDelta(b *testing.B) { benchmarkDecodeDelta(b, newEncoder("amino")) }
func BenchmarkJsonDecodeDelta(b *testing.B)  { benchmarkDecodeDelta(b, newEncoder("json")) }
func benchmarkDecodeDelta(b *testing.B, enc encoder) {
	wd := newTestWatchData()
	data, _ := enc.encodeFunc(wd)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.decodeFunc(data)
	}
}

type encoder interface {
	name() string
	encodeFunc(*watcher.WatchData) ([]byte, error)
	decodeFunc([]byte) (*watcher.WatchData, error)
}

func newEncoder(encType string) encoder {
	switch encType {
	case "amino":
		return &aminoEncoder{}
	case "json":
		return &jsonEncoder{}
	default:
	}
	panic("unsupport encoder")
}

// amino encoder
type aminoEncoder struct{}

func (ae *aminoEncoder) name() string { return "amino" }
func (ae *aminoEncoder) encodeFunc(data *watcher.WatchData) ([]byte, error) {
	return data.MarshalToAmino(nil)
}
func (ae *aminoEncoder) decodeFunc(data []byte) (*watcher.WatchData, error) {
	wd := &watcher.WatchData{}
	err := wd.UnmarshalFromAmino(nil, data)
	return wd, err
}

// json encoder
type jsonEncoder struct{}

func (je *jsonEncoder) name() string { return "json" }
func (je *jsonEncoder) encodeFunc(data *watcher.WatchData) ([]byte, error) {
	return json.Marshal(data)
}
func (je *jsonEncoder) decodeFunc(data []byte) (*watcher.WatchData, error) {
	wd := &watcher.WatchData{}
	err := json.Unmarshal(data, wd)
	return wd, err
}
