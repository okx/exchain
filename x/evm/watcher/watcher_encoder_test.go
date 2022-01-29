package watcher

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

var (
	testAccAddr1 = sdk.AccAddress("0x01")
	testAccAddr2 = sdk.AccAddress("0x02")
)
var (
	jsonEnc = jsoniter.ConfigCompatibleWithStandardLibrary
	cdc     = amino.NewCodec()
)

var testWatchData = []*WatchData{
	{},
	{
		DirtyAccount: []*sdk.AccAddress{&testAccAddr1, &testAccAddr2},
	},
	{
		Batches: []*Batch{{Key: []byte("0x01")}, {Value: []byte("0x01")}, {TypeValue: 1}},
	},
	{
		DelayEraseKey: [][]byte{[]byte("0x01"), []byte("0x02")},
	},
	{
		BloomData: []*types.KV{{Key: []byte("0x01")}, {Value: []byte("0x01")}},
	},
	{
		DirtyList: [][]byte{[]byte("0x01"), []byte("0x02")},
	},
	{
		DirtyAccount:  []*sdk.AccAddress{&testAccAddr1, {}, &testAccAddr2},
		Batches:       []*Batch{{Key: []byte("0x01")}, {}, {TypeValue: 1}},
		DelayEraseKey: [][]byte{[]byte("0x01"), {}, []byte("0x02")},
		BloomData:     []*types.KV{{Key: []byte("0x01")}, {}, {Value: []byte("0x01")}},
		DirtyList:     [][]byte{[]byte("0x01"), {}, []byte("0x02")},
	},
	{
		DirtyAccount:  []*sdk.AccAddress{&testAccAddr1, &testAccAddr2},
		Batches:       []*Batch{{Key: []byte("0x01")}, {Value: []byte("0x02")}, {TypeValue: 1}},
		DelayEraseKey: [][]byte{[]byte("0x01"), []byte("0x02")},
		BloomData:     []*types.KV{{Key: []byte("0x01")}, {Value: []byte("0x01")}},
		DirtyList:     [][]byte{[]byte("0x01"), []byte("0x02")},
	},
}

func newTestWatchData() *WatchData {
	return testWatchData[len(testWatchData)-1]
}

func TestWatchDataEncoder(t *testing.T) { testWatchDataAmino(t) }
func testWatchDataAmino(t *testing.T) {
	for i, wd := range testWatchData {
		expect, err := cdc.MarshalBinaryBare(wd)
		require.NoError(t, err, fmt.Sprintf("num %v", i))

		actual, err := wd.MarshalToAmino(cdc)
		require.NoError(t, err, fmt.Sprintf("num %v", i))
		require.EqualValues(t, expect, actual, fmt.Sprintf("num %v", i))

		var expectValue WatchData
		err = cdc.UnmarshalBinaryBare(expect, &expectValue)
		require.NoError(t, err, fmt.Sprintf("num %v", i))

		var actualValue WatchData
		err = actualValue.UnmarshalFromAmino(cdc, expect)
		require.NoError(t, err, fmt.Sprintf("num %v", i))
		require.EqualValues(t, expectValue, actualValue, fmt.Sprintf("num %v", i))
	}
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
	encodeFunc(*WatchData) ([]byte, error)
	decodeFunc([]byte) (*WatchData, error)
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
func (ae *aminoEncoder) encodeFunc(data *WatchData) ([]byte, error) {
	return data.MarshalToAmino(nil)
}
func (ae *aminoEncoder) decodeFunc(data []byte) (*WatchData, error) {
	wd := &WatchData{}
	err := wd.UnmarshalFromAmino(nil, data)
	return wd, err
}

// json encoder
type jsonEncoder struct{}

func (je *jsonEncoder) name() string { return "json" }
func (je *jsonEncoder) encodeFunc(data *WatchData) ([]byte, error) {
	return jsonEnc.Marshal(data)
}
func (je *jsonEncoder) decodeFunc(data []byte) (*WatchData, error) {
	wd := &WatchData{}
	err := jsonEnc.Unmarshal(data, wd)
	return wd, err
}
