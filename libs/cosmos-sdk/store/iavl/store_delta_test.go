package iavl

import (
	"fmt"
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

type encodeFunc func(data interface{}) ([]byte, error)

func aminoEncDelta(data interface{}) ([]byte, error) {
	td, ok := data.(iavltree.TreeDelta)
	if !ok {
		return nil, fmt.Errorf("no supported this type")
	}
	return td.MarshalToAmino()
}

//test different encode function to handle delta
func TestEncodeDelta(t *testing.T) {
	testEncodeDelta(t, "amino", aminoEncDelta)
}
func testEncodeDelta(t *testing.T, name string, encFunc encodeFunc) {
	_, delta := newTestTreeDelta()
	_, err := encFunc(delta)
	require.NoError(t, err, name)
}
