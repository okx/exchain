package iavl

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/require"
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

var (
	cacheSize = 100
	treeData  = map[string]string{
		"hello": "goodbye",
		"aloha": "shalom",
	}
	nMoreData = 0
)

func randBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, _ = crand.Read(b)
	return b
}

// make a tree with data from above and save it
func newAlohaTree(t *testing.T, db dbm.DB) (*iavl.MutableTree, types.CommitID) {
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	for k, v := range treeData {
		tree.Set([]byte(k), []byte(v))
	}

	for i := 0; i < nMoreData; i++ {
		key := randBytes(12)
		value := randBytes(50)
		tree.Set(key, value)
	}

	hash, ver, _, err := tree.SaveVersion(false)
	require.Nil(t, err)

	return tree, types.CommitID{Version: ver, Hash: hash}
}

func TestLoadStore(t *testing.T) {
	db := dbm.NewMemDB()
	flatKVDB := dbm.NewMemDB()
	tree, _ := newAlohaTree(t, db)
	store := UnsafeNewStore(tree)

	// Create non-pruned height H
	require.True(t, tree.Set([]byte("hello"), []byte("hallo")))
	hash, verH, _, err := tree.SaveVersion(false)
	cIDH := types.CommitID{Version: verH, Hash: hash}
	require.Nil(t, err)

	// Create pruned height Hp
	require.True(t, tree.Set([]byte("hello"), []byte("hola")))
	hash, verHp, _, err := tree.SaveVersion(false)
	cIDHp := types.CommitID{Version: verHp, Hash: hash}
	require.Nil(t, err)

	// TODO: Prune this height

	// Create current height Hc
	require.True(t, tree.Set([]byte("hello"), []byte("ciao")))
	hash, verHc, _, err := tree.SaveVersion(false)
	cIDHc := types.CommitID{Version: verHc, Hash: hash}
	require.Nil(t, err)

	// Querying an existing store at some previous non-pruned height H
	hStore, err := store.GetImmutable(verH)
	require.NoError(t, err)
	require.Equal(t, string(hStore.Get([]byte("hello"))), "hallo")

	// Querying an existing store at some previous pruned height Hp
	hpStore, err := store.GetImmutable(verHp)
	require.NoError(t, err)
	require.Equal(t, string(hpStore.Get([]byte("hello"))), "hola")

	// Querying an existing store at current height Hc
	hcStore, err := store.GetImmutable(verHc)
	require.NoError(t, err)
	require.Equal(t, string(hcStore.Get([]byte("hello"))), "ciao")

	// Querying a new store at some previous non-pruned height H
	newHStore, err := LoadStore(db, flatKVDB, cIDH, false, 0)
	require.NoError(t, err)
	require.Equal(t, string(newHStore.Get([]byte("hello"))), "hallo")

	// Querying a new store at some previous pruned height Hp
	newHpStore, err := LoadStore(db, flatKVDB, cIDHp, false, 0)
	require.NoError(t, err)
	require.Equal(t, string(newHpStore.Get([]byte("hello"))), "hola")

	// Querying a new store at current height H
	newHcStore, err := LoadStore(db, flatKVDB, cIDHc, false, 0)
	require.NoError(t, err)
	require.Equal(t, string(newHcStore.Get([]byte("hello"))), "ciao")
}

func TestGetImmutable(t *testing.T) {
	db := dbm.NewMemDB()
	tree, cID := newAlohaTree(t, db)
	store := UnsafeNewStore(tree)

	require.True(t, tree.Set([]byte("hello"), []byte("adios")))
	hash, ver, _, err := tree.SaveVersion(false)
	cID = types.CommitID{Version: ver, Hash: hash}
	require.Nil(t, err)

	_, err = store.GetImmutable(cID.Version + 1)
	require.NoError(t, err)

	newStore, err := store.GetImmutable(cID.Version - 1)
	require.NoError(t, err)
	require.Equal(t, newStore.Get([]byte("hello")), []byte("goodbye"))

	newStore, err = store.GetImmutable(cID.Version)
	require.NoError(t, err)
	require.Equal(t, newStore.Get([]byte("hello")), []byte("adios"))

	res := newStore.Query(abci.RequestQuery{Data: []byte("hello"), Height: cID.Version, Path: "/key", Prove: true})
	require.Equal(t, res.Value, []byte("adios"))
	require.NotNil(t, res.Proof)

	require.Panics(t, func() { newStore.Set(nil, nil) })
	require.Panics(t, func() { newStore.Delete(nil) })
	require.Panics(t, func() { newStore.CommitterCommit(nil) })
}

func TestTestGetImmutableIterator(t *testing.T) {
	db := dbm.NewMemDB()
	tree, cID := newAlohaTree(t, db)
	store := UnsafeNewStore(tree)

	newStore, err := store.GetImmutable(cID.Version)
	require.NoError(t, err)

	iter := newStore.Iterator([]byte("aloha"), []byte("hellz"))
	expected := []string{"aloha", "hello"}
	var i int

	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}

	require.Equal(t, len(expected), i)
}

func TestIAVLStoreGetSetHasDelete(t *testing.T) {
	db := dbm.NewMemDB()
	tree, _ := newAlohaTree(t, db)
	iavlStore := UnsafeNewStore(tree)

	key := "hello"

	exists := iavlStore.Has([]byte(key))
	require.True(t, exists)

	value := iavlStore.Get([]byte(key))
	require.EqualValues(t, value, treeData[key])

	value2 := "notgoodbye"
	iavlStore.Set([]byte(key), []byte(value2))

	value = iavlStore.Get([]byte(key))
	require.EqualValues(t, value, value2)

	iavlStore.Delete([]byte(key))

	exists = iavlStore.Has([]byte(key))
	require.False(t, exists)
}

func TestIAVLStoreNoNilSet(t *testing.T) {
	db := dbm.NewMemDB()
	tree, _ := newAlohaTree(t, db)
	iavlStore := UnsafeNewStore(tree)
	require.Panics(t, func() { iavlStore.Set([]byte("key"), nil) }, "setting a nil value should panic")
}

func TestIAVLIterator(t *testing.T) {
	db := dbm.NewMemDB()
	tree, _ := newAlohaTree(t, db)
	iavlStore := UnsafeNewStore(tree)
	iter := iavlStore.Iterator([]byte("aloha"), []byte("hellz"))
	expected := []string{"aloha", "hello"}
	var i int

	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)

	iter = iavlStore.Iterator([]byte("golang"), []byte("rocks"))
	expected = []string{"hello"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)

	iter = iavlStore.Iterator(nil, []byte("golang"))
	expected = []string{"aloha"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)

	iter = iavlStore.Iterator(nil, []byte("shalom"))
	expected = []string{"aloha", "hello"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)

	iter = iavlStore.Iterator(nil, nil)
	expected = []string{"aloha", "hello"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)

	iter = iavlStore.Iterator([]byte("golang"), nil)
	expected = []string{"hello"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, treeData[expectedKey])
		i++
	}
	require.Equal(t, len(expected), i)
}

func TestIAVLReverseIterator(t *testing.T) {
	db := dbm.NewMemDB()

	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)

	iavlStore.Set([]byte{0x00}, []byte("0"))
	iavlStore.Set([]byte{0x00, 0x00}, []byte("0 0"))
	iavlStore.Set([]byte{0x00, 0x01}, []byte("0 1"))
	iavlStore.Set([]byte{0x00, 0x02}, []byte("0 2"))
	iavlStore.Set([]byte{0x01}, []byte("1"))

	var testReverseIterator = func(t *testing.T, start []byte, end []byte, expected []string) {
		iter := iavlStore.ReverseIterator(start, end)
		var i int
		for i = 0; iter.Valid(); iter.Next() {
			expectedValue := expected[i]
			value := iter.Value()
			require.EqualValues(t, string(value), expectedValue)
			i++
		}
		require.Equal(t, len(expected), i)
	}

	testReverseIterator(t, nil, nil, []string{"1", "0 2", "0 1", "0 0", "0"})
	testReverseIterator(t, []byte{0x00}, nil, []string{"1", "0 2", "0 1", "0 0", "0"})
	testReverseIterator(t, []byte{0x00}, []byte{0x00, 0x01}, []string{"0 0", "0"})
	testReverseIterator(t, []byte{0x00}, []byte{0x01}, []string{"0 2", "0 1", "0 0", "0"})
	testReverseIterator(t, []byte{0x00, 0x01}, []byte{0x01}, []string{"0 2", "0 1"})
	testReverseIterator(t, nil, []byte{0x01}, []string{"0 2", "0 1", "0 0", "0"})
}

func TestIAVLPrefixIterator(t *testing.T) {
	db := dbm.NewMemDB()
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)

	iavlStore.Set([]byte("test1"), []byte("test1"))
	iavlStore.Set([]byte("test2"), []byte("test2"))
	iavlStore.Set([]byte("test3"), []byte("test3"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(0)}, []byte("test4"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(1)}, []byte("test4"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(255)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(0)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(1)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(255)}, []byte("test4"))

	var i int

	iter := types.KVStorePrefixIterator(iavlStore, []byte("test"))
	expected := []string{"test1", "test2", "test3"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, expectedKey)
		i++
	}
	iter.Close()
	require.Equal(t, len(expected), i)

	iter = types.KVStorePrefixIterator(iavlStore, []byte{byte(55), byte(255), byte(255)})
	expected2 := [][]byte{
		{byte(55), byte(255), byte(255), byte(0)},
		{byte(55), byte(255), byte(255), byte(1)},
		{byte(55), byte(255), byte(255), byte(255)},
	}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected2[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, []byte("test4"))
		i++
	}
	iter.Close()
	require.Equal(t, len(expected), i)

	iter = types.KVStorePrefixIterator(iavlStore, []byte{byte(255), byte(255)})
	expected2 = [][]byte{
		{byte(255), byte(255), byte(0)},
		{byte(255), byte(255), byte(1)},
		{byte(255), byte(255), byte(255)},
	}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected2[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, []byte("test4"))
		i++
	}
	iter.Close()
	require.Equal(t, len(expected), i)
}

func TestIAVLReversePrefixIterator(t *testing.T) {
	db := dbm.NewMemDB()
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)

	iavlStore.Set([]byte("test1"), []byte("test1"))
	iavlStore.Set([]byte("test2"), []byte("test2"))
	iavlStore.Set([]byte("test3"), []byte("test3"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(0)}, []byte("test4"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(1)}, []byte("test4"))
	iavlStore.Set([]byte{byte(55), byte(255), byte(255), byte(255)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(0)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(1)}, []byte("test4"))
	iavlStore.Set([]byte{byte(255), byte(255), byte(255)}, []byte("test4"))

	var i int

	iter := types.KVStoreReversePrefixIterator(iavlStore, []byte("test"))
	expected := []string{"test3", "test2", "test1"}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, expectedKey)
		i++
	}
	require.Equal(t, len(expected), i)

	iter = types.KVStoreReversePrefixIterator(iavlStore, []byte{byte(55), byte(255), byte(255)})
	expected2 := [][]byte{
		{byte(55), byte(255), byte(255), byte(255)},
		{byte(55), byte(255), byte(255), byte(1)},
		{byte(55), byte(255), byte(255), byte(0)},
	}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected2[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, []byte("test4"))
		i++
	}
	require.Equal(t, len(expected), i)

	iter = types.KVStoreReversePrefixIterator(iavlStore, []byte{byte(255), byte(255)})
	expected2 = [][]byte{
		{byte(255), byte(255), byte(255)},
		{byte(255), byte(255), byte(1)},
		{byte(255), byte(255), byte(0)},
	}
	for i = 0; iter.Valid(); iter.Next() {
		expectedKey := expected2[i]
		key, value := iter.Key(), iter.Value()
		require.EqualValues(t, key, expectedKey)
		require.EqualValues(t, value, []byte("test4"))
		i++
	}
	require.Equal(t, len(expected), i)
}

func nextVersion(iStore *Store) {
	key := []byte(fmt.Sprintf("Key for tree: %d", iStore.LastCommitID().Version))
	value := []byte(fmt.Sprintf("Value for tree: %d", iStore.LastCommitID().Version))
	iStore.Set(key, value)
	iStore.CommitterCommit(nil)
}

func TestIAVLNoPrune(t *testing.T) {
	db := dbm.NewMemDB()
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)
	nextVersion(iavlStore)

	for i := 1; i < 100; i++ {
		for j := 1; j <= i; j++ {
			require.True(t, iavlStore.VersionExists(int64(j)),
				"Missing version %d with latest version %d. Should be storing all versions",
				j, i)
		}

		nextVersion(iavlStore)
	}
}

func TestIAVLStoreQuery(t *testing.T) {
	db := dbm.NewMemDB()
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)

	k1, v1 := []byte("key1"), []byte("val1")
	k2, v2 := []byte("key2"), []byte("val2")
	v3 := []byte("val3")

	ksub := []byte("key")
	KVs0 := []types.KVPair{}
	KVs1 := []types.KVPair{
		{Key: k1, Value: v1},
		{Key: k2, Value: v2},
	}
	KVs2 := []types.KVPair{
		{Key: k1, Value: v3},
		{Key: k2, Value: v2},
	}
	valExpSubEmpty := cdc.MustMarshalBinaryLengthPrefixed(KVs0)
	valExpSub1 := cdc.MustMarshalBinaryLengthPrefixed(KVs1)
	valExpSub2 := cdc.MustMarshalBinaryLengthPrefixed(KVs2)

	cid, _ := iavlStore.CommitterCommit(nil)
	ver := cid.Version
	query := abci.RequestQuery{Path: "/key", Data: k1, Height: ver}
	querySub := abci.RequestQuery{Path: "/subspace", Data: ksub, Height: ver}

	// query subspace before anything set
	qres := iavlStore.Query(querySub)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, valExpSubEmpty, qres.Value)

	// set data
	iavlStore.Set(k1, v1)
	iavlStore.Set(k2, v2)

	// set data without commit, doesn't show up
	qres = iavlStore.Query(query)
	require.Equal(t, uint32(0), qres.Code)
	require.Nil(t, qres.Value)

	// commit it, but still don't see on old version
	cid, _ = iavlStore.CommitterCommit(nil)
	qres = iavlStore.Query(query)
	require.Equal(t, uint32(0), qres.Code)
	require.Nil(t, qres.Value)

	// but yes on the new version
	query.Height = cid.Version
	qres = iavlStore.Query(query)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, v1, qres.Value)

	// and for the subspace
	qres = iavlStore.Query(querySub)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, valExpSub1, qres.Value)

	// modify
	iavlStore.Set(k1, v3)
	cid, _ = iavlStore.CommitterCommit(nil)

	// query will return old values, as height is fixed
	qres = iavlStore.Query(query)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, v1, qres.Value)

	// update to latest in the query and we are happy
	query.Height = cid.Version
	qres = iavlStore.Query(query)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, v3, qres.Value)
	query2 := abci.RequestQuery{Path: "/key", Data: k2, Height: cid.Version}

	qres = iavlStore.Query(query2)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, v2, qres.Value)
	// and for the subspace
	qres = iavlStore.Query(querySub)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, valExpSub2, qres.Value)

	// default (height 0) will show latest -1
	query0 := abci.RequestQuery{Path: "/key", Data: k1}
	qres = iavlStore.Query(query0)
	require.Equal(t, uint32(0), qres.Code)
	require.Equal(t, v1, qres.Value)
}

func testCommitDelta(t *testing.T) {
	emptyDelta := &iavl.TreeDelta{NodesDelta: map[string]*iavl.NodeJson{}, OrphansDelta: []*iavl.NodeJson{}, CommitOrphansDelta: map[string]int64{}}
	tmtypes.DownloadDelta = true
	iavl.SetProduceDelta(false)

	db := dbm.NewMemDB()
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(t, err)

	iavlStore := UnsafeNewStore(tree)

	k1, v1 := []byte("key1"), []byte("val1")
	k2, v2 := []byte("key2"), []byte("val2")

	// set data
	iavlStore.Set(k1, v1)
	iavlStore.Set(k2, v2)

	// normal case (not use delta and not produce delta)
	cid, treeDelta := iavlStore.CommitterCommit(nil)
	assert.NotEmpty(t, cid.Hash)
	assert.EqualValues(t, 1, cid.Version)
	assert.Equal(t, emptyDelta, treeDelta)

	// not use delta and produce delta
	iavl.SetProduceDelta(true)
	cid1, treeDelta1 := iavlStore.CommitterCommit(nil)
	assert.NotEmpty(t, cid1.Hash)
	assert.EqualValues(t, 2, cid1.Version)
	assert.NotEqual(t, emptyDelta, treeDelta1)

	// use delta and produce delta
	cid2, treeDelta2 := iavlStore.CommitterCommit(treeDelta1)
	assert.NotEmpty(t, cid2.Hash)
	assert.EqualValues(t, 3, cid2.Version)
	assert.NotEqual(t, emptyDelta, treeDelta2)
	assert.Equal(t, treeDelta1, treeDelta2)

	// use delta and not produce delta
	iavl.SetProduceDelta(false)
	cid3, treeDelta3 := iavlStore.CommitterCommit(treeDelta1)
	assert.NotEmpty(t, cid3.Hash)
	assert.EqualValues(t, 4, cid3.Version)
	assert.Equal(t, emptyDelta, treeDelta3)
}
func TestCommitDelta(t *testing.T) {
	if os.Getenv("SUB_PROCESS") == "1" {
		testCommitDelta(t)
		return
	}

	var outb, errb bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run=TestCommitDelta")
	cmd.Env = append(os.Environ(), "SUB_PROCESS=1")
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		isFailed := false
		if strings.Contains(outb.String(), "FAIL:") ||
			strings.Contains(errb.String(), "FAIL:") {
			fmt.Print(cmd.Stderr)
			fmt.Print(cmd.Stdout)
			isFailed = true
		}
		assert.Equal(t, isFailed, false)

		return
	}
}

func TestIAVLDelta(t *testing.T) {
	emptyDelta := iavl.TreeDelta{NodesDelta: map[string]*iavl.NodeJson{}, OrphansDelta: []*iavl.NodeJson{}, CommitOrphansDelta: map[string]int64{}}

	db := dbm.NewMemDB()
	tree, _ := newAlohaTree(t, db)

	// Create non-pruned height H
	require.True(t, tree.Set([]byte("hello"), []byte("hallo")))

	// normal case (not use delta and not produce delta)
	iavl.SetProduceDelta(false)
	h, v, delta, err := tree.SaveVersion(false)
	require.NoError(t, err)
	assert.NotEmpty(t, h)
	assert.EqualValues(t, 2, v)
	assert.Equal(t, delta, emptyDelta)

	// not use delta and produce delta
	iavl.SetProduceDelta(true)
	h1, v1, delta1, err := tree.SaveVersion(false)
	require.NoError(t, err)
	assert.NotEmpty(t, h1)
	assert.EqualValues(t, 3, v1)
	// delta is empty or not depends on SetProduceDelta()
	assert.NotEqual(t, delta1, emptyDelta)

	// use delta and produce delta
	tree.SetDelta(&delta1)
	h2, v2, delta2, err := tree.SaveVersion(true)
	require.NoError(t, err)
	assert.NotEmpty(t, h2)
	assert.EqualValues(t, 4, v2)
	assert.NotEqual(t, delta2, emptyDelta)
	assert.Equal(t, delta1, delta2)

	// use delta and not produce delta
	iavl.SetProduceDelta(false)
	tree.SetDelta(&delta1)
	h3, v3, delta3, err := tree.SaveVersion(true)
	require.NoError(t, err)
	assert.NotEmpty(t, h3)
	assert.EqualValues(t, 5, v3)
	assert.Equal(t, delta3, emptyDelta)
}

func BenchmarkIAVLIteratorNext(b *testing.B) {
	db := dbm.NewMemDB()
	treeSize := 1000
	tree, err := iavl.NewMutableTree(db, cacheSize)
	require.NoError(b, err)

	for i := 0; i < treeSize; i++ {
		key := randBytes(4)
		value := randBytes(50)
		tree.Set(key, value)
	}

	iavlStore := UnsafeNewStore(tree)
	iterators := make([]types.Iterator, b.N/treeSize)

	for i := 0; i < len(iterators); i++ {
		iterators[i] = iavlStore.Iterator([]byte{0}, []byte{255, 255, 255, 255, 255})
	}

	b.ResetTimer()
	for i := 0; i < len(iterators); i++ {
		iter := iterators[i]
		for j := 0; j < treeSize; j++ {
			iter.Next()
		}
	}
}
