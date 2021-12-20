package iavl

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"strconv"
	"testing"

	db "github.com/tendermint/tm-db"
)

func TestDelete(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	tree.set([]byte("k1"), []byte("Fred"))
	hash, version, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)

	require.NoError(t, tree.DeleteVersion(version))

	k1Value, _, _ := tree.GetVersionedWithProof([]byte("k1"), version)
	require.Nil(t, k1Value)

	key := tree.ndb.rootKey(version)
	err = memDB.Set(key, hash)
	require.NoError(t, err)
	tree.versions.Set(version, true)

	k1Value, _, err = tree.GetVersionedWithProof([]byte("k1"), version)
	require.Nil(t, err)
	require.Equal(t, 0, bytes.Compare([]byte("Fred"), k1Value))
}

func TestTraverse(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	for i := 0; i < 6; i++ {
		tree.set([]byte(fmt.Sprintf("k%d", i)), []byte(fmt.Sprintf("v%d", i)))
	}

	require.Equal(t, 11, tree.nodeSize(), "Size of tree unexpected")
}

func TestMutableTree_DeleteVersions(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	type entry struct {
		key   []byte
		value []byte
	}

	versionEntries := make(map[int64][]entry)

	// create 10 tree versions, each with 1000 random key/value entries
	for i := 0; i < 10; i++ {
		entries := make([]entry, 1000)

		for j := 0; j < 1000; j++ {
			k := randBytes(10)
			v := randBytes(10)

			entries[j] = entry{k, v}
			_ = tree.Set(k, v)
		}

		_, v, _, err := tree.SaveVersion(false)
		require.NoError(t, err)

		versionEntries[v] = entries
	}

	// delete even versions
	versionsToDelete := []int64{2, 4, 6, 8}
	require.NoError(t, tree.DeleteVersions(versionsToDelete...))

	// ensure even versions have been deleted
	for _, v := range versionsToDelete {
		require.False(t, tree.versions.Get(v))

		_, err := tree.LazyLoadVersion(v)
		require.Error(t, err)
	}

	// ensure odd number versions exist and we can query for all set entries
	for _, v := range []int64{1, 3, 5, 7, 9, 10} {
		require.True(t, tree.versions.Get(v))

		_, err := tree.LazyLoadVersion(v)
		require.NoError(t, err)

		for _, e := range versionEntries[v] {
			_, val := tree.Get(e.key)
			require.Equal(t, e.value, val)
		}
	}
}

func TestMutableTree_DeleteVersionsRange(t *testing.T) {
	require := require.New(t)

	mdb := db.NewMemDB()
	tree, err := NewMutableTree(mdb, 0)
	require.NoError(err)

	const maxLength = 100
	const fromLength = 10

	versions := make([]int64, 0, maxLength)
	for count := 1; count <= maxLength; count++ {
		versions = append(versions, int64(count))
		countStr := strconv.Itoa(count)
		// Set kv pair and save version
		tree.Set([]byte("aaa"), []byte("bbb"))
		tree.Set([]byte("key"+countStr), []byte("value"+countStr))
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(err, "SaveVersion should not fail")
	}

	tree, err = NewMutableTree(mdb, 0)
	require.NoError(err)
	targetVersion, err := tree.LoadVersion(int64(maxLength))
	require.NoError(err)
	require.Equal(targetVersion, int64(maxLength), "targetVersion shouldn't larger than the actual tree latest version")

	err = tree.DeleteVersionsRange(fromLength, int64(maxLength/2))
	require.NoError(err, "DeleteVersionsTo should not fail")

	for _, version := range versions[:fromLength-1] {
		require.True(tree.versions.Get(version), "versions %d no more than 10 should exist", version)

		v, err := tree.LazyLoadVersion(version)
		require.NoError(err, version)
		require.Equal(v, version)

		_, value := tree.Get([]byte("aaa"))
		require.Equal(string(value), "bbb")

		for _, count := range versions[:version] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
	}

	for _, version := range versions[fromLength : int64(maxLength/2)-1] {
		require.False(tree.versions.Get(version), "versions %d more 10 and no more than 50 should have been deleted", version)

		_, err := tree.LazyLoadVersion(version)
		require.Error(err)
	}

	for _, version := range versions[int64(maxLength/2)-1:] {
		require.True(tree.versions.Get(version), "versions %d more than 50 should exist", version)

		v, err := tree.LazyLoadVersion(version)
		require.NoError(err)
		require.Equal(v, version)

		_, value := tree.Get([]byte("aaa"))
		require.Equal(string(value), "bbb")

		for _, count := range versions[:fromLength] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
		for _, count := range versions[int64(maxLength/2)-1 : version] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
	}
}

func TestMutableTree_InitialVersion(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 9})
	require.NoError(t, err)

	tree.Set([]byte("a"), []byte{0x01})
	_, version, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	assert.EqualValues(t, 10, version)

	tree.Set([]byte("b"), []byte{0x02})
	_, version, _, err = tree.SaveVersion(false)
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	// Reloading the tree with the same initial version is fine
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 9})
	require.NoError(t, err)
	version, err = tree.Load()
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	// Reloading the tree with an initial version beyond the lowest should error
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 11})
	require.NoError(t, err)
	_, err = tree.Load()
	require.Error(t, err)

	// Reloading the tree with a lower initial version is fine, and new versions can be produced
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 3})
	require.NoError(t, err)
	version, err = tree.Load()
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	tree.Set([]byte("c"), []byte{0x03})
	_, version, _, err = tree.SaveVersion(false)
	require.NoError(t, err)
	assert.EqualValues(t, 12, version)
}

func TestMutableTree_SetInitialVersion(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)
	tree.SetInitialVersion(9)

	tree.Set([]byte("a"), []byte{0x01})
	_, version, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	assert.EqualValues(t, 10, version)
}

func BenchmarkMutableTree_Set(b *testing.B) {
	db := db.NewDB("test", db.MemDBBackend, "")
	t, err := NewMutableTree(db, 100000)
	require.NoError(b, err)
	for i := 0; i < 1000000; i++ {
		t.Set(randBytes(10), []byte{})
	}
	b.ReportAllocs()
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Set(randBytes(10), []byte{})
	}
}

func BenchmarkUpdateBranch(b *testing.B) {
	nodeNums := 100000
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		CommitIntervalHeight = 10
		memDB := db.NewMemDB()
		tree, _ := NewMutableTree(memDB, 10000)
		_, _, _, _ = tree.SaveVersion(false)
		ks := "k%d"
		vs := "v%d"

		for i := 0; i < nodeNums; i++ {
			k := fmt.Sprintf(ks, i)
			v := fmt.Sprintf(vs, i)
			tree.Set([]byte(k), []byte(v))
		}
		tree.ndb.updateBranch(tree.root, map[string]*Node{})
		treeMap.resetMap()
	}
}

func BenchmarkUpdateBranchParallelV1(b *testing.B) {
	nodeNums := 100000
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		CommitIntervalHeight = 10
		memDB := db.NewMemDB()
		tree, _ := NewMutableTree(memDB, 10000)
		_, _, _, _ = tree.SaveVersion(false)
		ks := "k%d"
		vs := "v%d"

		for i := 0; i < nodeNums; i++ {
			k := fmt.Sprintf(ks, i)
			v := fmt.Sprintf(vs, i)
			tree.Set([]byte(k), []byte(v))
		}
		tree.ndb.updateBranchParallelV1(tree.root, map[string]*Node{}, 1)
		treeMap.resetMap()
	}
}

func BenchmarkUpdateBranchParallelV2(b *testing.B) {
	nodeNums := 100000
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		CommitIntervalHeight = 10
		memDB := db.NewMemDB()
		tree, _ := NewMutableTree(memDB, 10000)
		_, _, _, _ = tree.SaveVersion(false)
		ks := "k%d"
		vs := "v%d"

		for i := 0; i < nodeNums; i++ {
			k := fmt.Sprintf(ks, i)
			v := fmt.Sprintf(vs, i)
			tree.Set([]byte(k), []byte(v))
		}
		tree.ndb.updateBranchParallelV2(tree.root, map[string]*Node{})
		treeMap.resetMap()
	}
}
