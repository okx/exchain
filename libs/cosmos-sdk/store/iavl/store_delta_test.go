package iavl

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	iavltree "github.com/okex/exchain/libs/iavl"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tm-db"
)

func newTestTreeDelta() (types.CommitID, iavltree.TreeDelta) {
	tmtypes.DownloadDelta = true
	db := dbm.NewMemDB()
	tree, err := iavltree.NewMutableTree(db, cacheSize)
	if err != nil {
		panic(err)
	}
	iavlStore := UnsafeNewStore(tree)

	k1, v1 := []byte("key1"), []byte("val1")
	k2, v2 := []byte("key2"), []byte("val2")

	// set data
	iavlStore.Set(k1, v1)
	iavlStore.Set(k2, v2)

	iavltree.SetProduceDelta(true)
	cid, treeDelta, _ := iavlStore.Commit(nil, nil)
	return cid, treeDelta
}

//test newTestTreeDelta
func TestTreeDelta(t *testing.T) { testTreeDelta(t) }
func testTreeDelta(t *testing.T) {
	cid, tdata := newTestTreeDelta()

	emptyDelta := iavltree.TreeDelta{
		//NodesDelta:         map[string]*iavltree.NodeJson{},
		OrphansDelta:       []*iavltree.NodeJson{},
		CommitOrphansDelta: map[string]int64{},
	}
	assert.NotEmpty(t, cid.Hash)
	assert.EqualValues(t, 1, cid.Version)
	assert.NotEqual(t, emptyDelta, tdata)
}

type encoder interface {
	name() string
	encodeFunc(iavltree.TreeDelta) ([]byte, error)
	decodeFunc([]byte) (iavltree.TreeDelta, error)
}

type aminoEncoder struct{}

func newAmino() *aminoEncoder         { return &aminoEncoder{} }
func (ae *aminoEncoder) name() string { return "amino" }
func (ae *aminoEncoder) encodeFunc(data iavltree.TreeDelta) ([]byte, error) {
	return data.MarshalToAmino()
}
func (ae *aminoEncoder) decodeFunc(data []byte) (iavltree.TreeDelta, error) {
	td := iavltree.TreeDelta{
		NodesDelta:         map[string]*iavltree.NodeJson{},
		OrphansDelta:       make([]*iavltree.NodeJson, 0),
		CommitOrphansDelta: map[string]int64{},
	}
	err := td.UnmarshalFromAmino(data)
	return td, err
}

// different encode function to handle delta
func TestEncodeTreeDelta(t *testing.T) {
	testEncodeTreeDelta(t, newAmino())
}
func testEncodeTreeDelta(t *testing.T, enc encoder) {
	_, delta := newTestTreeDelta()
	_, err := enc.encodeFunc(delta)
	require.NoError(t, err, enc.name())
}

// decode data after it is marshaled to bytes, and the result is same before
func TestDecodeTreeDelta(t *testing.T) {
	testDecodeTreeDelta(t, newAmino())
}
func testDecodeTreeDelta(t *testing.T, enc encoder) {
	_, delta1 := newTestTreeDelta()
	data, err := enc.encodeFunc(delta1)

	delta2, err := enc.decodeFunc(data)
	require.NoError(t, err, enc.name())
	assert.EqualValues(t, delta1, delta2, enc.name())
}
