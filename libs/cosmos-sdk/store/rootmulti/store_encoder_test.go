package rootmulti

import (
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	iavltree "github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/require"
)

func newTestTreeDeltaMap() iavltree.TreeDeltaMap {
	// new multiplex store
	var db dbm.DB = dbm.NewMemDB()
	ms := newMultiStoreWithMounts(db, types.PruneNothing)
	ms.LoadLatestVersion()

	// set value into store map
	k, v := []byte("wind"), []byte("blows")
	k1, v1 := []byte("key1"), []byte("val1")
	k2, v2 := []byte("key2"), []byte("val2")
	store1 := ms.getStoreByName("store1").(types.KVStore)
	store1.Set(k, v)
	store1.Set(k1, v1)
	store2 := ms.getStoreByName("store2").(types.KVStore)
	store2.Set(k2, v2)

	// each store to be committed and return its delta
	_, returnedDeltas := ms.CommitterCommitMap(nil)

	return returnedDeltas
}

// test decode function
func TestAminoDecoder(t *testing.T) { testDecodeTreeDelta(t, newEncoder("amino")) }
func TestJsonDecoder(t *testing.T)  { testDecodeTreeDelta(t, newEncoder("json")) }
func testDecodeTreeDelta(t *testing.T, enc encoder) {
	deltaList1 := newTestTreeDeltaMap()
	data, err := enc.encodeFunc(deltaList1)
	require.NoError(t, err, enc.name())

	deltaList2, err := enc.decodeFunc(data)
	require.NoError(t, err, enc.name())
	require.Equal(t, deltaList1, deltaList2)
}

// benchmark encode performance
func BenchmarkAminoEncodeDelta(b *testing.B) { benchmarkEncodeDelta(b, newEncoder("amino")) }
func BenchmarkJsonEncodeDelta(b *testing.B)  { benchmarkEncodeDelta(b, newEncoder("json")) }
func benchmarkEncodeDelta(b *testing.B, enc encoder) {
	data := newTestTreeDeltaMap()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.encodeFunc(data)
	}

}

// benchmark decode performance
func BenchmarkAminoDecodeDelta(b *testing.B) { benchmarkDecodeDelta(b, newEncoder("amino")) }
func BenchmarkJsonDecodeDelta(b *testing.B)  { benchmarkDecodeDelta(b, newEncoder("json")) }
func benchmarkDecodeDelta(b *testing.B, enc encoder) {
	deltaList1 := newTestTreeDeltaMap()
	data, _ := enc.encodeFunc(deltaList1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.decodeFunc(data)
	}
}

type encoder interface {
	name() string
	encodeFunc(iavltree.TreeDeltaMap) ([]byte, error)
	decodeFunc([]byte) (iavltree.TreeDeltaMap, error)
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
func (ae *aminoEncoder) encodeFunc(data iavltree.TreeDeltaMap) ([]byte, error) {
	return data.MarshalToAmino(nil)
}
func (ae *aminoEncoder) decodeFunc(data []byte) (iavltree.TreeDeltaMap, error) {
	deltaList := iavltree.TreeDeltaMap{}
	err := deltaList.UnmarshalFromAmino(nil, data)
	return deltaList, err
}

// json encoder
type jsonEncoder struct{}

func (je *jsonEncoder) name() string { return "json" }
func (je *jsonEncoder) encodeFunc(data iavltree.TreeDeltaMap) ([]byte, error) {
	return itjs.Marshal(data)
}
func (je *jsonEncoder) decodeFunc(data []byte) (iavltree.TreeDeltaMap, error) {
	deltaList := iavltree.TreeDeltaMap{}
	err := itjs.Unmarshal(data, &deltaList)
	return deltaList, err
}
